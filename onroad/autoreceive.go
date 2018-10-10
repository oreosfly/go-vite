package onroad

import (
	"github.com/vitelabs/go-vite/common/types"
	"github.com/vitelabs/go-vite/generator"
	"github.com/vitelabs/go-vite/ledger"
	"github.com/vitelabs/go-vite/log15"
	"github.com/vitelabs/go-vite/onroad/model"
	"math/big"
	"sync"
)

type SimpleAutoReceiveFilterPair struct {
	tti      types.TokenTypeId
	minValue big.Int
}

type AutoReceiveWorker struct {
	log     log15.Logger
	address types.Address

	manager          *Manager
	generator        *generator.Generator
	onroadBlocksPool *model.OnroadBlocksPool

	status           int
	isSleeping       bool
	breaker          chan struct{}
	stopListener     chan struct{}
	newOnroadTxAlarm chan struct{}

	filters map[types.TokenTypeId]big.Int

	statusMutex sync.Mutex
}

func NewAutoReceiveWorker(manager *Manager, address types.Address, filters map[types.TokenTypeId]big.Int) *AutoReceiveWorker {
	return &AutoReceiveWorker{
		manager:          manager,
		onroadBlocksPool: manager.onroadBlocksPool,
		generator:        generator.NewGenerator(manager.vite.Chain(), manager.vite.WalletManager().KeystoreManager),
		address:          address,
		status:           Create,
		isSleeping:       false,
		filters:          filters,
		log:              slog.New("worker", "a", "addr", address),
	}
}

func (w *AutoReceiveWorker) Start() {

	w.log.Info("Start()", "current status", w.status)
	w.statusMutex.Lock()
	defer w.statusMutex.Unlock()
	if w.status != Start {

		w.breaker = make(chan struct{})
		w.newOnroadTxAlarm = make(chan struct{})
		w.stopListener = make(chan struct{})

		w.onroadBlocksPool.AddCommonTxLis(w.address, func() {
			w.NewOnroadTxAlarm()
		})

		w.onroadBlocksPool.AcquireFullOnroadBlocksCache(w.address)

		go w.startWork()

		w.status = Start

	} else {
		// awake it in order to run at least once
		w.NewOnroadTxAlarm()
	}
	w.log.Info("end start")
}

func (w *AutoReceiveWorker) Stop() {
	w.log.Info("Stop()", "current status", w.status)
	w.statusMutex.Lock()
	defer w.statusMutex.Unlock()
	if w.status == Start {

		w.onroadBlocksPool.ReleaseFullOnroadBlocksCache(w.address)

		w.breaker <- struct{}{}
		close(w.breaker)

		w.onroadBlocksPool.RemoveCommonTxLis(w.address)
		close(w.newOnroadTxAlarm)

		// make sure we can stop the worker
		<-w.stopListener
		close(w.stopListener)

		w.status = Stop
	}
	w.log.Info("stopped")
}

func (w *AutoReceiveWorker) ResetAutoReceiveFilter(filters map[types.TokenTypeId]big.Int) {
	w.log.Info("ResetAutoReceiveFilter", "len", len(filters))
	w.filters = filters
	if w.Status() == Start {
		w.onroadBlocksPool.ResetCacheCursor(w.address)
	}
}

func (w *AutoReceiveWorker) startWork() {
	w.log.Info("startWork")
LOOP:
	for {
		w.isSleeping = false
		if w.Status() == Stop {
			break
		}

		tx := w.onroadBlocksPool.GetNextCommonTx(w.address)
		if tx != nil {
			minAmount, ok := w.filters[tx.TokenId]
			if !ok || tx.Amount.Cmp(&minAmount) < 0 {
				continue
			}
			w.ProcessOneBlock(tx)
			continue
		}

		w.isSleeping = true
		select {
		case <-w.newOnroadTxAlarm:
			w.log.Info("start awake")
		case <-w.breaker:
			w.log.Info("worker broken")
			break LOOP
		}
	}

	w.log.Info("startWork end called")
	w.stopListener <- struct{}{}
	w.log.Info("startWork end")
}

func (w *AutoReceiveWorker) Close() error {
	w.Stop()
	return nil
}

func (w AutoReceiveWorker) Status() int {
	w.statusMutex.Lock()
	defer w.statusMutex.Unlock()
	return w.status
}

func (w *AutoReceiveWorker) NewOnroadTxAlarm() {
	if w.isSleeping {
		w.newOnroadTxAlarm <- struct{}{}
	}
}

func (w *AutoReceiveWorker) ProcessOneBlock(sendBlock *ledger.AccountBlock) {
	if w.manager.checkExistInPool(sendBlock.ToAddress, sendBlock.FromBlockHash) {
		w.log.Info("ProcessOneBlock.checkExistInPool failed")
		return
	}

	err := w.generator.PrepareVm(nil, nil, &sendBlock.ToAddress)
	if err != nil {
		w.log.Error("ProcessOneBlock.PrepareVm failed", "error", err)
		return
	}

	genResult, genErr := w.generator.GenerateWithOnroad(*sendBlock, nil,
		func(addr types.Address, data []byte) (signedData, pubkey []byte, err error) {
			return w.generator.Sign(addr, nil, data)
		})
	if genErr != nil {
		w.log.Error("GenerateTx error ignore, ", "error", genErr)
		return
	}

	if genResult.BlockGenList == nil {
		w.log.Error("GenerateWithOnroad failed, BlockGenList is nil")
		return
	}

	poolErr := w.manager.insertCommonBlockToPool(genResult.BlockGenList)
	if poolErr != nil {
		w.log.Error("insertCommonBlockToPool failed, ", "error", poolErr)
		return
	}
}