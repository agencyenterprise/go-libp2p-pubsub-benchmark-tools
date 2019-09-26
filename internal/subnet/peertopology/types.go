package peertopology

import (
	"github.com/agencyenterprise/gossip-host/internal/host"
)

// PeerTopology is a peering algorithm that connects hosts
type PeerTopology interface {
	Build(hosts []*host.Host) error
}
