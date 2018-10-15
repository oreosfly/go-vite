package net

import (
	"errors"
	"fmt"
	"github.com/vitelabs/go-vite/ledger"
	"github.com/vitelabs/go-vite/p2p"
	"github.com/vitelabs/go-vite/vite/net/message"
	"sort"
	"strconv"
	"strings"
	"time"
)

type reqState byte

const (
	reqWaiting reqState = iota
	reqPending
	reqRespond
	reqDone
	reqError
	reqCancel
)

var reqStatus = [...]string{
	reqWaiting: "waiting",
	reqPending: "pending",
	reqRespond: "respond",
	reqDone:    "done",
	reqError:   "error",
	reqCancel:  "canceled",
}

func (s reqState) String() string {
	return reqStatus[s]
}

type context struct {
	*syncer
	peers *peerSet
	pool  RequestPool
	fc    *fileClient
}

type Request interface {
	Handle(ctx *context, msg *p2p.Msg, peer *Peer)
	ID() uint64
	Run(ctx *context)
	Done(err error)
	Expired() bool
	State() reqState
}

type receiveBlocks func(sblocks []*ledger.SnapshotBlock, ablocks []*ledger.AccountBlock)
type doneCallback func(id uint64, err error)

var errMissingPeer = errors.New("request missing peer")
var errUnknownResErr = errors.New("unknown response exception")
var errUnExpectedRes = errors.New("unexpected response")
var errMaxRetry = errors.New("max Retry")

const minBlocks uint64 = 3600  // minimal snapshot blocks per subLedger request
const maxBlocks uint64 = 10800 // maximal snapshot blocks per subLedger request
const maxRetry = 3

type subLedgerPiece struct {
	from, to uint64
	peer     *Peer
}

// split large subledger request to many small pieces
func splitSubLedger(from, to uint64, peers Peers) (cs []*subLedgerPiece) {
	// sort peers from low to high
	sort.Sort(peers)

	total := to - from + 1
	if total < minBlocks {
		cs = append(cs, &subLedgerPiece{
			from: from,
			to:   to,
			peer: peers[len(peers)-1], // choose the tallest peer
		})
		return
	}

	var pCount, pTo uint64 // piece length
	for _, peer := range peers {
		if peer.height > from+minBlocks {
			pTo = peer.height - minBlocks
			pCount = pTo - from + 1

			// peer not high enough
			if pCount < minBlocks {
				continue
			}

			// piece too large
			if pCount > maxBlocks {
				pCount = maxBlocks
			}

			// piece to
			pTo = from + pCount - 1

			// piece to exceed target height
			if pTo > to {
				pTo = to
				pCount = to - from + 1
			}
			// reset piece is too small, then collapse to one piece
			if to < pTo+minBlocks {
				pCount += to - pTo
			}

			cs = append(cs, &subLedgerPiece{
				from: from,
				to:   pTo,
				peer: peer,
			})

			from = pTo + 1
			if from > to {
				break
			}
		}
	}

	// reset piece, alloc to best peer
	if from < to {
		// if peer is not too much, has rest chunk, then collapse the rest chunk with last chunk
		lastChunk := cs[len(cs)-1]
		if lastChunk.peer == peers[len(peers)-1] {
			lastChunk.to = to
		} else {
			cs = append(cs, &subLedgerPiece{
				from: from,
				to:   to,
				peer: peers[len(peers)-1],
			})
		}
	}

	return
}

func splitChunk(from, to uint64) (chunks [][2]uint64) {
	// chunks may be only one block, then from == to
	if from > to || to == 0 {
		return
	}

	total := (to-from)/minBlocks + 1
	chunks = make([][2]uint64, total)

	var cTo uint64
	var i int
	for from <= to {
		if to > from+minBlocks {
			cTo = from + minBlocks - 1
		} else {
			cTo = to
		}

		chunks[i] = [2]uint64{from, cTo}

		from = cTo + 1
		i++
	}

	return chunks
}

// @request for subLedger, will get FileList and Chunk
type subLedgerRequest struct {
	id         uint64 // id & child_id
	from, to   uint64
	_retry     int
	peer       *Peer
	state      reqState
	file       *fileRequest
	chunks     []*chunkRequest
	expiration time.Time
	done       doneCallback
}

func (s *subLedgerRequest) State() reqState {
	return s.state
}

func (s *subLedgerRequest) Retry(ctx *context, peer *Peer) {
	if s._retry >= maxRetry {
		s.Done(errMaxRetry)
		return
	}

	peers := ctx.peers.Pick(peer.height)
	if len(peers) != 0 {
		s.peer = peers[0]
		ctx.pool.Retry(s.id)
		s._retry++
	} else {
		s.Done(errMissingPeer)
	}
}

func (s *subLedgerRequest) Handle(ctx *context, pkt *p2p.Msg, peer *Peer) {
	if cmd(pkt.Cmd) == FileListCode {
		msg := new(message.FileList)
		err := msg.Deserialize(pkt.Payload)
		if err != nil {
			fmt.Println("subLedgerRequest handle error: ", err)
			s.Retry(ctx, peer)
			return
		}

		if len(msg.Files) != 0 {
			// sort as StartHeight
			sort.Sort(files(msg.Files))
			// request files
			s.file = &fileRequest{
				files: msg.Files,
				nonce: msg.Nonce,
				peer:  peer,
				rec:   ctx.syncer.receiveBlocks,
				done:  s.childDone,
			}
			ctx.fc.request(s.file)
		}

		// request chunks
		for _, chunk := range msg.Chunks {
			// maybe chunk is too large
			cs := splitChunk(chunk[0], chunk[1])
			for _, c := range cs {
				msgId := ctx.pool.MsgID()

				c := &chunkRequest{
					id:         msgId,
					from:       c[0],
					to:         c[1],
					peer:       peer,
					expiration: time.Now().Add(30 * time.Second),
					done:       s.childDone,
				}
				s.chunks = append(s.chunks, c)

				ctx.pool.Add(c)
			}
		}
	} else {
		s.Retry(ctx, peer)
	}
}

func (s *subLedgerRequest) ID() uint64 {
	return s.id
}

func (s *subLedgerRequest) Run(*context) {
	err := s.peer.Send(GetSubLedgerCode, s.id, &message.GetSubLedger{
		From:    ledger.HashHeight{Height: s.from},
		Count:   s.to - s.from + 1,
		Forward: true,
	})

	if err != nil {
		s.peer.log.Error(fmt.Sprintf("send GetSubLedgerMsg<%d-%d> to %s error: %v", s.from, s.to, s.peer, err))
		s.Done(err)
	} else {
		s.state = reqPending
		s.peer.log.Info(fmt.Sprintf("send GetSubLedgerMsg<%d-%d> to %s done", s.from, s.to, s.peer))
	}
}

func (s *subLedgerRequest) childDone(id uint64, err error) {

}

func (s *subLedgerRequest) Done(err error) {
	if err != nil {
		s.state = reqError
	} else {
		s.state = reqDone
	}

	s.done(s.id, err)
}

func (s *subLedgerRequest) Expired() bool {
	return time.Now().After(s.expiration)
}

// @request file
type fileRequest struct {
	id      uint64
	state   reqState
	files   []*ledger.CompressedFileMeta
	nonce   uint64
	peer    *Peer
	rec     receiveBlocks
	current uint64 // the tallest snapshotBlock have received, as the breakpoint resume
	done    func(id uint64, err error)
}

// file1/0-3600 file2/3601-7200
func (r *fileRequest) String() string {
	files := make([]string, len(r.files))
	for i, file := range r.files {
		files[i] = file.Filename + "/" +
			strconv.FormatUint(file.StartHeight, 10) +
			"-" +
			strconv.FormatUint(file.EndHeight, 10)
	}

	return strings.Join(files, " ")
}

func (r *fileRequest) Done(err error) {
	if err != nil {
		r.state = reqError
	} else {
		r.state = reqDone
	}

	r.done(r.id, err)
}

func (r *fileRequest) Addr() string {
	return r.peer.FileAddress().String()
}

func (r *fileRequest) Retry(ctx *context) {
	for i, file := range r.files {
		if r.current < file.EndHeight {
			r.files = r.files[i:]
			break
		}
	}

	ps := ctx.peers.Pick(r.peer.height)
	if len(ps) > 0 {
		r.peer = ps[0]
		r.Run(ctx)
	} else {
		r.state = reqError
		r.Done(errMissingPeer)
	}
}
func (r *fileRequest) Run(ctx *context) {
	ctx.fc.request(r)
}

// @request for chunk
type chunkRequest struct {
	id         uint64
	from, to   uint64
	_retry     int
	peer       *Peer
	state      reqState
	expiration time.Time
	done       doneCallback
}

func (c *chunkRequest) State() reqState {
	return c.state
}

func (c *chunkRequest) Retry(ctx *context, peer *Peer) {
	if c._retry >= maxRetry {
		c.Done(errMaxRetry)
		return
	}

	// find taller peers
	peers := ctx.peers.Pick(peer.height)
	if len(peers) != 0 {
		c.peer = peers[0]
		ctx.pool.Retry(c.id)
		c._retry++
	} else {
		c.Done(errMissingPeer)
	}
}

func (c *chunkRequest) Handle(ctx *context, pkt *p2p.Msg, peer *Peer) {
	if cmd(pkt.Cmd) == SubLedgerCode {
		msg := new(message.SubLedger)
		err := msg.Deserialize(pkt.Payload)
		if err != nil {
			fmt.Println("chunkRequest handle error: ", err)
			c.Retry(ctx, peer)
			return
		}

		ctx.syncer.receiveBlocks(msg.SBlocks, msg.ABlocks)
		c.Done(nil)
	} else {
		c.Retry(ctx, peer)
	}
}

func (c *chunkRequest) ID() uint64 {
	return c.id
}

func (c *chunkRequest) Run(*context) {
	err := c.peer.Send(GetChunkCode, c.id, &message.GetChunk{
		Start: c.from,
		End:   c.to,
	})

	if err != nil {
		c.state = reqError
		c.peer.log.Error(fmt.Sprintf("send GetChunkMsg<%d/%d> to %s error: %v", c.from, c.to, c.peer, err))
	} else {
		c.state = reqPending
		c.peer.log.Info(fmt.Sprintf("send GetChunkMsg<%d/%d> to %s done", c.from, c.to, c.peer))
	}
}

func (c *chunkRequest) Done(err error) {
	if err != nil {
		c.state = reqError
	} else {
		c.state = reqDone
	}

	c.done(c.id, err)
}

func (c *chunkRequest) Expired() bool {
	return time.Now().After(c.expiration)
}

// helper
type files []*ledger.CompressedFileMeta

func (a files) Len() int {
	return len(a)
}

func (a files) Less(i, j int) bool {
	return a[i].StartHeight < a[j].StartHeight
}

func (a files) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
