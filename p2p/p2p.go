// Package p2p implements the vite P2P network

package p2p

import (
	"errors"
	"fmt"
	"github.com/vitelabs/go-vite/crypto/ed25519"
	"github.com/vitelabs/go-vite/log15"
	"github.com/vitelabs/go-vite/monitor"
	"github.com/vitelabs/go-vite/p2p/block"
	"github.com/vitelabs/go-vite/p2p/discovery"
	"github.com/vitelabs/go-vite/p2p/nat"
	"net"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

var p2pServerLog = log15.New("module", "p2p/server")
var errSvrStarted = errors.New("server has started")
var errSvrStopped = errors.New("server has stopped")

type Discovery interface {
	Lookup(discovery.NodeID) []*discovery.Node
	Resolve(discovery.NodeID) *discovery.Node
	RandomNodes([]*discovery.Node) int
	Start()
	Stop()
	SubNodes(ch chan<- *discovery.Node)
}

type Config struct {
	Name            string
	NetID           NetworkID          // which network server runs on
	MaxPeers        uint               // max peers can be connected
	MaxPendingPeers uint               // max peers waiting for connect
	MaxInboundRatio uint               // max inbound peers: MaxPeers / MaxInboundRatio
	Port            uint               // TCP and UDP listen port
	Database        string             // the directory for storing node table
	PrivateKey      ed25519.PrivateKey // use for encrypt message, the corresponding public key use for NodeID
	Protocols       []*Protocol        // protocols server supported
	BootNodes       []string
}

type Server struct {
	*Config
	running   int32          // atomic
	wg        sync.WaitGroup // Wait for all jobs done
	term      chan struct{}
	pending   chan struct{} // how many connection can wait for handshake
	addPeer   chan *conn
	delPeer   chan *Peer
	discv     Discovery
	handshake *Handshake
	BootNodes []*discovery.Node
	peers     *PeerSet
	blockList *block.CuckooSet
	self      *discovery.Node
	//agent          *agent
	log      log15.Logger
	ln       net.Listener
	nodeChan chan *discovery.Node
}

func New(cfg Config) (svr *Server, err error) {
	safeCfg := EnsureConfig(cfg)

	ID, err := discovery.Priv2NodeID(cfg.PrivateKey)
	if err != nil {
		return
	}

	addr := "0.0.0.0:" + strconv.FormatUint(uint64(cfg.Port), 10)
	// udp discover
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return
	}

	log := log15.New("module", "p2p/server")

	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return
	}
	log.Info(fmt.Sprintf("udp listen at %s", udpAddr))

	// tcp listener
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		// close the open udp connection
		udpConn.Close()
		return
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		// close the open udp connection
		udpConn.Close()
		return
	}
	log.Info(fmt.Sprintf("tcp listen at %s", tcpAddr))

	svr = &Server{
		Config:    safeCfg,
		log:       log,
		peers:     NewPeerSet(),
		term:      make(chan struct{}),
		pending:   make(chan struct{}, cfg.MaxPendingPeers),
		addPeer:   make(chan *conn, 1),
		delPeer:   make(chan *Peer, 1),
		BootNodes: addFirmNodes(cfg.BootNodes),
		blockList: block.NewCuckooSet(100),
		ln:        listener,
		self: &discovery.Node{
			ID:  ID,
			IP:  udpAddr.IP,
			UDP: uint16(udpAddr.Port),
			TCP: uint16(tcpAddr.Port),
		},
		nodeChan: make(chan *discovery.Node, 10),
	}

	svr.discv = discovery.New(&discovery.Config{
		Priv:      svr.PrivateKey,
		DBPath:    svr.Database,
		BootNodes: svr.BootNodes,
		Conn:      udpConn,
		Self:      svr.self,
	})

	svr.discv.SubNodes(svr.nodeChan)

	//svr.agent = newAgent(svr)

	return
}

func (svr *Server) Peers() []*PeerInfo {
	return svr.peers.Info()
}

func (svr *Server) PeersCount() (amount int) {
	return svr.peers.Size()
}

func (svr *Server) NodeInfo() *NodeInfo {
	protocols := make([]string, len(svr.Protocols))
	for i, protocol := range svr.Protocols {
		protocols[i] = protocol.String()
	}

	return &NodeInfo{
		ID:    svr.self.ID.String(),
		Name:  svr.Name,
		Url:   svr.self.String(),
		NetID: svr.NetID,
		Address: &address{
			IP:  svr.self.IP,
			TCP: svr.self.TCP,
			UDP: svr.self.UDP,
		},
		Protocols: protocols,
	}
}

func (svr *Server) URL() string {
	return svr.self.String()
}

func (svr *Server) Available() bool {
	return svr.PeersCount() > 0
}

func (svr *Server) maxOutboundPeers() uint {
	return svr.MaxPeers - svr.maxInboundPeers()
}

func (svr *Server) maxInboundPeers() uint {
	return svr.MaxPeers / svr.MaxInboundRatio
}

func (svr *Server) Start() error {
	if !atomic.CompareAndSwapInt32(&svr.running, 0, 1) {
		return errSvrStarted
	}

	svr.log.Info("p2p server start")

	err := svr.setHandshake()
	if err != nil {
		return err
	}

	// mapping udp and tcp
	go nat.Map(svr.term, "udp", int(svr.self.UDP), int(svr.self.UDP), "vite p2p udp", 0, svr.updateNode)
	go nat.Map(svr.term, "tcp", int(svr.self.TCP), int(svr.self.TCP), "vite p2p tcp", 0, svr.updateNode)

	svr.discv.Start()

	svr.wg.Add(1)
	go svr.dialLoop()

	// tcp listener
	svr.wg.Add(1)
	go svr.listenLoop()

	// peer manager
	svr.wg.Add(1)
	go svr.loop()

	svr.log.Info("p2p server started")
	return nil
}

func (svr *Server) updateNode(addr *nat.Addr) {
	if addr.Proto == "tcp" {
		svr.self.TCP = uint16(addr.Port)
	} else {
		svr.self.UDP = uint16(addr.Port)
	}
}

func (svr *Server) setHandshake() error {
	cmdsets := make([]*CmdSet, len(svr.Protocols))
	for i, p := range svr.Protocols {
		cmdsets[i] = p.CmdSet()
	}

	svr.handshake = &Handshake{
		Version: Version,
		Name:    svr.Name,
		NetID:   svr.NetID,
		ID:      svr.self.ID,
		CmdSets: cmdsets,
	}

	return nil
}

func (svr *Server) dialLoop() {
	defer svr.wg.Done()

	dialer := &net.Dialer{
		Timeout: 3 * time.Second,
	}

	for {
		select {
		case <-svr.term:
			return
		case svr.pending <- struct{}{}:
			var node *discovery.Node
			for node = range svr.nodeChan {
				if err := svr.checkConn(node.ID, outbound); err == nil {
					break
				}
			}

			svr.log.Info(fmt.Sprintf("got node: %s", node))
			if conn, err := dialer.Dial("tcp", node.TCPAddr().String()); err == nil {
				go svr.setupConn(conn, outbound)
			} else {
				svr.log.Error(fmt.Sprintf("dial node %s failed: %v", node, err))
			}
		}
	}
}

func (svr *Server) listenLoop() {
	defer svr.wg.Done()
	defer svr.ln.Close()

	var conn net.Conn
	var err error

	for {
		select {
		case svr.pending <- struct{}{}:
			for {
				conn, err = svr.ln.Accept()

				if err == nil {
					break
				}
			}

			go svr.setupConn(conn, inbound)
		case <-svr.term:
			return
		}
	}
}

func (svr *Server) setupConn(c net.Conn, flag connFlag) {
	ts := &conn{
		AsyncMsgConn: NewAsyncMsgConn(c, nil),
		flags:        flag,
	}

	svr.log.Info(fmt.Sprintf("begin handshake with %s", c.RemoteAddr()))

	// handshake data, add remoteIP and remotePort
	handshake := *svr.handshake
	tcpAddr := c.RemoteAddr().(*net.TCPAddr)
	handshake.RemoteIP = tcpAddr.IP
	handshake.RemotePort = uint16(tcpAddr.Port)
	data, err := handshake.Serialize()
	if err != nil {
		ts.Close(nil)
		return
	}
	sig := ed25519.Sign(svr.PrivateKey, data)
	data = append(sig, data...)

	their, err := ts.Handshake(data)

	if err != nil {
		ts.Close(err)
		svr.log.Error(fmt.Sprintf("handshake with %s error: %v", c.RemoteAddr(), err))
	} else if their.NetID != svr.NetID {
		err = fmt.Errorf("different NetID: our %s, their %s", svr.NetID, their.NetID)
		ts.Close(err)
		svr.log.Error(fmt.Sprintf("handshake with %s error: %v", c.RemoteAddr(), err))
	} else {
		ts.name = their.Name
		ts.cmdSets = their.CmdSets

		// use to discribe the connection
		ts.remoteID = their.ID
		ts.remoteIP = handshake.RemoteIP
		ts.remotePort = handshake.RemotePort

		ts.localID = svr.self.ID
		ts.localIP = their.RemoteIP
		ts.localPort = their.RemotePort

		svr.log.Info(fmt.Sprintf("handshake with %s@%s done", ts.remoteID, c.RemoteAddr()))
		svr.addPeer <- ts
	}

	<-svr.pending
}

func (svr *Server) checkConn(id discovery.NodeID, flag connFlag) error {
	if uint(svr.peers.Size()) >= svr.MaxPeers {
		return DiscTooManyPeers
	}

	if flag.is(inbound) && uint(svr.peers.inbound) >= svr.maxInboundPeers() {
		return DiscTooManyPassivePeers
	}

	if svr.peers.Has(id) {
		return DiscAlreadyConnected
	}

	if id == svr.self.ID {
		return DiscSelf
	}

	return nil
}

func (svr *Server) loop() {
	defer svr.wg.Done()

	//shouldSchedule := make(chan struct{})

	//go svr.agent.scheduleTasks(svr.term, shouldSchedule)

loop:
	for {
		select {
		case <-svr.term:
			break loop
		//case <-shouldSchedule:
		//	svr.agent.createTasks()
		case c := <-svr.addPeer:
			err := svr.checkConn(c.remoteID, c.flags)

			if err == nil {
				if p, err := NewPeer(c, svr.Protocols); err == nil {
					svr.peers.Add(p)

					peersCount := svr.peers.Size()
					svr.log.Info(fmt.Sprintf("create new peer %s, total: %d", p, peersCount))
					monitor.LogDuration("p2p/peer", "add", int64(peersCount))

					go svr.runPeer(p)
				} else {
					svr.log.Error(fmt.Sprintf("create new peer error: %v", err))
				}
			} else {
				c.Close(err)
				svr.log.Error("cannot create new peer", "error", err)
			}

		case p := <-svr.delPeer:
			svr.peers.Del(p)

			peersCount := svr.peers.Size()
			svr.log.Info("delete peer", "ID", p.ID().String(), "total", peersCount)
			monitor.LogDuration("p2p/peer", "del", int64(peersCount))
		}
	}

	svr.peers.Traverse(func(id discovery.NodeID, p *Peer) {
		p.Disconnect(DiscQuitting)
	})
}

func (svr *Server) runPeer(p *Peer) {
	err := p.run()
	if err != nil {
		svr.log.Error("run peer error", "error", err)
	}
	svr.delPeer <- p
}

func (svr *Server) Stop() {
	if !atomic.CompareAndSwapInt32(&svr.running, 1, 0) {
		return
	}

	close(svr.term)
	svr.discv.Stop()
	svr.wg.Wait()
}

// @section NodeInfo
type NodeInfo struct {
	ID        string    `json:"remoteID"`
	Name      string    `json:"name"`
	Url       string    `json:"url"`
	NetID     NetworkID `json:"netId"`
	Address   *address  `json:"address"`
	Protocols []string  `json:"protocols"`
}

type address struct {
	IP  net.IP `json:"ip"`
	TCP uint16 `json:"tcp"`
	UDP uint16 `json:"udp"`
}