package subnet

import (
	"context"
	"net"

	"github.com/agencyenterprise/gossip-host/pkg/logger"
	"github.com/agencyenterprise/gossip-host/pkg/subnet/config"
)

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
	_ = hosts

	// build network topology

	return nil
}
