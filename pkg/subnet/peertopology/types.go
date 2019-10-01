package peertopology

import (
	"github.com/agencyenterprise/gossip-host/pkg/cerr"
	"github.com/agencyenterprise/gossip-host/pkg/host"
)

const (
	// ErrUnknownTopology is thrown when a passed topology is not known
	ErrUnknownTopology = cerr.Error("unknown peering topology")
)

// PeerTopology is a peering algorithm that connects hosts
type PeerTopology interface {
	Build(hosts []*host.Host) error
}
