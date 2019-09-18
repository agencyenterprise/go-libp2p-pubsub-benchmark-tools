package config

import (
	"github.com/agencyenterprise/gossip-host/pkg/cerr"
)

// Config is a struct to hold the config options
type Config struct {
	Host Host
}

// Host contains configs for the host
type Host struct {
	// Listen are addresses on which to listen
	Listens []string `json:"listens"`
	// Peers are peers to be bootstrapped (e.g. /ip4/127.0.0.1/tcp/63785/ipfs/QmWjz6xb8v9K4KnYEwP5Yk75k5mMBCehzWFLCvvQpYxF3d)
	Peers []string `json:"peers"`
	// Transports are the transport protocols which the host is to use (e.g. "tcp", "ws", etc)
	Transports []string `json:"transports"`
	// Muxers are the transport muxers (e.g. yamux, mplex, etc.)
	Muxers [][]string `json:"muxers"`
	// Security specifies the security to use
	Security string `json:"security"`
	// Disable relay disables the relay
	DisableRelay bool `json:"disableRelay"`
}

// ErrNilConfig is returned when a config is expected but none is given
const ErrNilConfig = cerr.Error("unknown nil config")
