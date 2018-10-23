package api

import (
	"github.com/vitelabs/go-vite/chain"
	"github.com/vitelabs/go-vite/common/types"
	"github.com/vitelabs/go-vite/log15"
	"github.com/vitelabs/go-vite/vite"
	"github.com/vitelabs/go-vite/vm/contracts"
	"github.com/vitelabs/go-vite/vm/util"
	"github.com/vitelabs/go-vite/vm_context"
	"math/big"
	"sort"
	"time"
)

type PledgeApi struct {
	chain chain.Chain
	log   log15.Logger
}

func NewPledgeApi(vite *vite.Vite) *PledgeApi {
	return &PledgeApi{
		chain: vite.Chain(),
		log:   log15.New("module", "rpc_api/pledge_api"),
	}
}

func (p PledgeApi) String() string {
	return "PledgeApi"
}

func (p *PledgeApi) GetPledgeData(beneficialAddr types.Address) ([]byte, error) {
	return contracts.ABIPledge.PackMethod(contracts.MethodNamePledge, beneficialAddr)
}

func (p *PledgeApi) GetCancelPledgeData(beneficialAddr types.Address, amount *big.Int) ([]byte, error) {
	return contracts.ABIPledge.PackMethod(contracts.MethodNameCancelPledge, beneficialAddr, amount)
}

type QuotaAndTxNum struct {
	Quota uint64 `json:"quota"`
	TxNum uint64 `json:"txNum"`
}

func (p *PledgeApi) GetPledgeQuota(addr types.Address) QuotaAndTxNum {
	q := p.chain.GetPledgeQuota(p.chain.GetLatestSnapshotBlock().Hash, addr)
	return QuotaAndTxNum{q, q / util.TxGas}
}

type PledgeInfo struct {
	Amount         *big.Int      `json:"amount"`
	WithdrawHeight uint64        `json:"withdrawHeight"`
	BeneficialAddr types.Address `json:"beneficialAddr"`
	WithdrawTime   int64         `json:"withdrawTime"`
}
type byWithdrawHeight []*contracts.PledgeInfo

func (a byWithdrawHeight) Len() int      { return len(a) }
func (a byWithdrawHeight) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a byWithdrawHeight) Less(i, j int) bool {
	if a[i].WithdrawHeight == a[j].WithdrawHeight {
		return a[i].BeneficialAddr.String() < a[j].BeneficialAddr.String()
	}
	return a[i].WithdrawHeight < a[j].WithdrawHeight
}

func (p *PledgeApi) GetPledgeList(addr types.Address, index int, count int) ([]*PledgeInfo, error) {
	snapshotBlock := p.chain.GetLatestSnapshotBlock()
	vmContext, err := vm_context.NewVmContext(p.chain, &snapshotBlock.Hash, nil, nil)
	if err != nil {
		return nil, err
	}
	list := contracts.GetPledgeInfoList(vmContext, addr)
	sort.Sort(byWithdrawHeight(list))
	startHeight, endHeight := index*count, (index+1)*count
	if startHeight >= len(list) {
		return []*PledgeInfo{}, nil
	}
	if endHeight > len(list) {
		endHeight = len(list)
	}
	targetList := make([]*PledgeInfo, endHeight-startHeight)
	for i, info := range list[startHeight:endHeight] {
		targetList[i] = &PledgeInfo{
			info.Amount,
			info.WithdrawHeight,
			info.BeneficialAddr,
			getWithdrawTime(snapshotBlock.Timestamp, snapshotBlock.Height, info.WithdrawHeight)}
	}
	return targetList, nil
}

const (
	secondBetweenSnapshotBlocks int64 = 1
)

func getWithdrawTime(snapshotTime *time.Time, snapshotHeight uint64, withdrawHeight uint64) int64 {
	return snapshotTime.Unix() + int64(withdrawHeight-snapshotHeight)*secondBetweenSnapshotBlocks
}