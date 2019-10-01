package subnet

import (
	"github.com/agencyenterprise/gossip-host/pkg/cerr"
)

const (
	// ErrIPOutOfCIDRRange is returned when an IP is incremented but out of CIDR range
	ErrIPOutOfCIDRRange = cerr.Error("ip out of CIDR range")
	// ErrNilIPNet is thrown when an IPNet pointer is nil
	ErrNilIPNet = cerr.Error("ip net is nil")
	// ErrNilPort is thrown when a port pointer is nil
	ErrNilPort = cerr.Error("port is nil")
)
