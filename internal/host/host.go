package host

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/agencyenterprise/gossip-host/internal/config"
	"github.com/agencyenterprise/gossip-host/pkg/logger"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/routing"
	kaddht "github.com/libp2p/go-libp2p-kad-dht"
	mplex "github.com/libp2p/go-libp2p-mplex"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	secio "github.com/libp2p/go-libp2p-secio"
	yamux "github.com/libp2p/go-libp2p-yamux"
	"github.com/libp2p/go-libp2p/p2p/discovery"
	tcp "github.com/libp2p/go-tcp-transport"
	ws "github.com/libp2p/go-ws-transport"
	"github.com/multiformats/go-multiaddr"
)

type mdnsNotifee struct {
	h   host.Host
	ctx context.Context
}

// HandlePeerFound...
func (m *mdnsNotifee) HandlePeerFound(pi peer.AddrInfo) {
	m.h.Connect(m.ctx, pi)
}

// Start starts a new gossip host
func Start(config *config.Config) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	transports := libp2p.ChainOptions(
		libp2p.Transport(tcp.NewTCPTransport),
		libp2p.Transport(ws.New),
	)

	muxers := libp2p.ChainOptions(
		libp2p.Muxer("/yamux/1.0.0", yamux.DefaultTransport),
		libp2p.Muxer("/mplex/6.7.0", mplex.DefaultTransport),
	)

	security := libp2p.Security(secio.ID, secio.New)

	listenAddrs := libp2p.ListenAddrStrings(
		"/ip4/0.0.0.0/tcp/0",
		"/ip4/0.0.0.0/tcp/0/ws",
	)

	var dht *kaddht.IpfsDHT
	newDHT := func(h host.Host) (routing.PeerRouting, error) {
		var err error
		dht, err = kaddht.New(ctx, h)
		if err != nil {
			logger.Errorf("err creating new kaddht\n%v", err)
		}

		return dht, err
	}
	routing := libp2p.Routing(newDHT)

	host, err := libp2p.New(
		ctx,
		transports,
		listenAddrs,
		muxers,
		security,
		routing,
	)
	if err != nil {
		logger.Errorf("err creating new libp2p host\n%v", err)
		return err
	}

	ps, err := pubsub.NewGossipSub(ctx, host)
	if err != nil {
		logger.Errorf("err creating new gossip sub\n%v", err)
		return err
	}
	sub, err := ps.Subscribe(pubsubTopic)
	if err != nil {
		logger.Errorf("err subscribing\n%v", err)
		return err
	}
	go pubsubHandler(ctx, sub)

	for _, addr := range host.Addrs() {
		logger.Infof("Listening on %v", addr)
	}

	targetAddr, err := multiaddr.NewMultiaddr("/ip4/127.0.0.1/tcp/63785/ipfs/QmWjz6xb8v9K4KnYEwP5Yk75k5mMBCehzWFLCvvQpYxF3d")
	if err != nil {
		logger.Errorf("err parsing targetAddr from multiaddr\n%v", err)
		return err
	}

	targetInfo, err := peer.AddrInfoFromP2pAddr(targetAddr)
	if err != nil {
		logger.Errorf("err parsing targetInfo from peer addr\n%v", err)
		return err
	}

	err = host.Connect(ctx, *targetInfo)
	if err != nil {
		logger.Errorf("err connecting\n%v", err)
		return err
	}

	fmt.Println("Connected to", targetInfo.ID)

	mdns, err := discovery.NewMdnsService(ctx, host, time.Second*10, "")
	if err != nil {
		logger.Errorf("err discovering\n%v", err)
		return err
	}
	mdns.RegisterNotifee(&mdnsNotifee{h: host, ctx: ctx})

	err = dht.Bootstrap(ctx)
	if err != nil {
		logger.Errorf("err bootstrapping\n%v", err)
		return err
	}

	donec := make(chan struct{}, 1)
	go chatInputLoop(ctx, host, ps, donec)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT)

	select {
	case <-stop:
		logger.Info("shutting down...")
		host.Close()
		logger.Info("exiting...")

	case <-donec:
		logger.Info("shutting down...")
		host.Close()
	}

	return nil
}
