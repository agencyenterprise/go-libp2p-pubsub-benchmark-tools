package subnet

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/agencyenterprise/gossip-host/pkg/host"
	"github.com/agencyenterprise/gossip-host/pkg/logger"
	"github.com/agencyenterprise/gossip-host/pkg/subnet/config"
	"github.com/agencyenterprise/gossip-host/pkg/subnet/peertopology"
)

// Start begins the subnet
func Start(ctx context.Context, conf config.Config) error {
	// parse pubsub cidr
	pubsubIP, pubsubNet, err := net.ParseCIDR(conf.Subnet.PubsubCIDR)
	if err != nil {
		logger.Errorf("err parsing pubsub CIDRs %s:\n%v", conf.Subnet.PubsubCIDR, err)
		return err
	}

	// parse RPC cidr
	rpcIP, rpcNet, err := net.ParseCIDR(conf.Subnet.RPCCIDR)
	if err != nil {
		logger.Errorf("err parsing rpc CIDRs %s:\n%v", conf.Subnet.RPCCIDR, err)
		return err
	}

	// build hosts
	hosts, err := buildHosts(ctx, conf, pubsubIP, rpcIP, pubsubNet, rpcNet, conf.Subnet.PubsubPortRange, conf.Subnet.RPCPortRange)
	if err != nil {
		logger.Errorf("err bulding hosts:\n%v", err)
		return err
	}

	// build the host pubsubs and rpc
	ch := make(chan error)
	if err = buildPubsubAndRPC(ch, hosts); err != nil {
		logger.Errorf("err bulding pubsub and rpc:\n%v", err)
		return err
	}

	// build network topology
	if err = peertopology.ConnectPeersForTopology(conf.Subnet.PeerTopology, hosts); err != nil {
		logger.Errorf("err building topology for %s:\n%v", conf.Subnet.PeerTopology, err)
		return err
	}

	// build router/discover
	if err = buildDiscovery(hosts); err != nil {
		logger.Errorf("err building discovery:\n%v", err)
		return err
	}

	// capture the ctrl+c signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT)

	for _, h := range hosts {
		// TODO: need to add ready signal when host has started
		go func(ch chan error, stop chan os.Signal, hst *host.Host) {
			if err = hst.Start(ch, stop); err != nil {
				logger.Errorf("host id %s err:\n%v", hst.ID(), err)
			}
		}(ch, stop, h)
	}

	select {
	case <-stop:
		// note: I don't like '^C' showing up on the same line as the next logged line...
		fmt.Println("")
		logger.Info("Received stop signal from os. Shutting down...")

	case err := <-ch:
		logger.Errorf("received err on rpc channel:\n%v", err)
		return err

	case <-ctx.Done():
		if err := ctx.Err(); err != nil {
			logger.Errorf("err on the context:\n%v", err)
			return err
		}
	}

	return nil
}
