package peertopology

import (
	"strings"

	"github.com/agencyenterprise/gossip-host/pkg/host"
	"github.com/agencyenterprise/gossip-host/pkg/logger"
	"github.com/agencyenterprise/gossip-host/pkg/subnet/peertopology/whiteblocks"
)

// ConnectPeersForTopology builds the specified topology
func ConnectPeersForTopology(topology string, hosts []*host.Host) error {
	switch strings.ToLower(topology) {
	case "whiteblocks":
		return whiteblocks.Build(hosts)

	default:
		logger.Errorf("unknown peering topology %s", topology)
		return ErrUnknownTopology
	}
}
