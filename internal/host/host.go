package host

import (
	"context"
	"crypto/rand"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/agencyenterprise/gossip-host/internal/config"
	rpcHost "github.com/agencyenterprise/gossip-host/internal/grpc/host"
	"github.com/agencyenterprise/gossip-host/pkg/logger"

	"github.com/libp2p/go-libp2p"
	connmgr "github.com/libp2p/go-libp2p-connmgr"
	lcrypto "github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/routing"
	kaddht "github.com/libp2p/go-libp2p-kad-dht"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	lconfig "github.com/libp2p/go-libp2p/config"
	"github.com/libp2p/go-libp2p/p2p/discovery"
)

type mdnsNotifee struct {
	h   host.Host
	ctx context.Context
}

// HandlePeerFound...
func (m *mdnsNotifee) HandlePeerFound(pi peer.AddrInfo) {
	logger.Infof("peer found: %v", pi)
	m.h.Connect(m.ctx, pi)
}

// Start starts a new gossip host
func Start(conf *config.Config) error {
	if conf == nil {
		logger.Error("nil config")
		return config.ErrNilConfig
	}

	var lOpts []lconfig.Option

	// create a context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// add private key
	if conf.Host.Priv == nil {
		priv, _, err := lcrypto.GenerateECDSAKeyPair(rand.Reader)
		if err != nil {
			logger.Errorf("err generating private key:\n%v", err)
			return err
		}

		lOpts = append(lOpts, libp2p.Identity(priv))
	} else {
		lOpts = append(lOpts, libp2p.Identity(conf.Host.Priv))
	}

	// create transports
	transports, err := parseTransportOptions(conf.Host.Transports)
	if err != nil {
		logger.Errorf("err parsing transports\n%v", err)
		return err
	}
	lOpts = append(lOpts, transports)

	// create muxers
	muxers, err := parseMuxerOptions(conf.Host.Muxers)
	if err != nil {
		logger.Errorf("err parsing muxers\n%v", err)
		return err
	}
	lOpts = append(lOpts, muxers)

	// create security
	security, err := parseSecurityOptions(conf.Host.Security)
	if err != nil {
		logger.Errorf("err parsing security\n%v", err)
		return err
	}
	lOpts = append(lOpts, security)

	// add listen addresses
	if len(conf.Host.Listens) > 0 {
		lOpts = append(lOpts, libp2p.ListenAddrStrings(conf.Host.Listens...))
	}

	// Conn manager
	if !conf.Host.OmitConnectionManager {
		cm := connmgr.NewConnManager(256, 512, 120)
		lOpts = append(lOpts, libp2p.ConnectionManager(cm))
	}

	// NAT port map
	if !conf.Host.OmitNATPortMap {
		lOpts = append(lOpts, libp2p.NATPortMap())
	}

	// create router
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
	lOpts = append(lOpts, routing)

	if conf.Host.OmitRelay {
		lOpts = append(lOpts, libp2p.DisableRelay())
	}

	// build the libp2p host
	host, err := libp2p.New(ctx, lOpts...)
	if err != nil {
		logger.Errorf("err creating new libp2p host\n%v", err)
		return err
	}
	defer host.Close()

	// build the gossip pub/sub
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
	go pubsubHandler(ctx, host.ID(), sub)

	// Start the RPC server
	publisher := &publisher{ps}
	rHost := rpcHost.New(publisher.publish)
	go func() {
		if err := rHost.Listen(ctx, conf.Host.RPCAddress); err != nil {
			logger.Errorf("err listening on rpc:\n%v", err)
		}
	}()

	for i, addr := range host.Addrs() {
		logger.Infof("listening #%d on: %s/ipfs/%s\n", i, addr, host.ID().Pretty())
	}

	// connect to peers
	if err := bootstrapPeers(ctx, host, conf.Host.Peers); err != nil {
		logger.Errorf("err bootstrapping peers\n%v", err)
		return err
	}

	// create discovery service
	mdns, err := discovery.NewMdnsService(ctx, host, time.Second*10, "")
	if err != nil {
		logger.Errorf("err discovering\n%v", err)
		return err
	}
	mdns.RegisterNotifee(&mdnsNotifee{h: host, ctx: ctx})

	// note: is there a reason this is after the creation of the discovery service, or can it be moved up with dht initialization?
	if err = dht.Bootstrap(ctx); err != nil {
		logger.Errorf("err bootstrapping\n%v", err)
		return err
	}

	// capture the ctrl+c signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT)

	// lock the thread
	select {
	case <-stop:
		// note: I don't like '^C' showing up on the same line as the next logged line...
		fmt.Println("")
		logger.Info("Received stop signal from os. Shutting down...")
	}

	return nil
}
