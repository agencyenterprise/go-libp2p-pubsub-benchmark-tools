package config

import (
	"github.com/agencyenterprise/gossip-host/pkg/cerr"

	lcrypto "github.com/libp2p/go-libp2p-core/crypto"
)

const (
	defaultsLoc       string = "./defaults"
	defaultConfigName string = "config.default.json"
	defaultPEMName    string = "private_key.pem"
)

// Config is a struct to hold the config options
type Config struct {
	Host Host `json:"host,omitempty"`
}

// Host contains configs for the host
type Host struct {
	// PrivPEM is the host's private key location in PKCS#8, ASN.1 DER PEM format
	PrivPEM string `json:"privPEM,omitempty"`
	// Priv is the parsed, host's private key
	Priv lcrypto.PrivKey
	// Listen are addresses on which to listen
	Listens []string `json:"listens,omitempty"`
	// RPCAddress is the address to listen on for RPC
	RPCAddress string `json:"rcpAddress,omitempty"`
	// Peers are peers to be bootstrapped (e.g. /ip4/127.0.0.1/tcp/63785/ipfs/QmWjz6xb8v9K4KnYEwP5Yk75k5mMBCehzWFLCvvQpYxF3d)
	Peers []string `json:"peers,omitempty"`
	// Transports are the transport protocols which the host is to use (e.g. "tcp", "ws", etc)
	Transports []string `json:"transports,omitempty"`
	// Muxers are the transport muxers (e.g. yamux, mplex, etc.)
	Muxers [][]string `json:"muxers,omitempty"`
	// Security specifies the security to use
	Security string `json:"security,omitempty"`
	// OmitRelay disables the relay
	OmitRelay bool `json:"omitRelay,omitempty"`
	// OmitConnectionManager enables the connection manager
	OmitConnectionManager bool `json:"omitConnectionManager,omitempty"`
	// OmitNatPortMap enables the nat port map
	OmitNATPortMap bool `json:"omitNATPortMap,omitempty"`
}

// ErrNilConfig is returned when a config is expected but none is given
const ErrNilConfig = cerr.Error("unknown nil config")

// ErrIncorrectKeyType is returned when the private key is not of the correct type
const ErrIncorrectKeyType = cerr.Error("incorrect private key type")
