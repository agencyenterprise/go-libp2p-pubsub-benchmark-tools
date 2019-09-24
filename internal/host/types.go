package host

import (
	"github.com/agencyenterprise/gossip-host/pkg/cerr"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

// ErrUnknownTransportOption is returned when an unknown transport has been specified
const ErrUnknownTransportOption = cerr.Error("unknown transport option")

// ErrImproperTransportOption is returned when an improper transport has been specified (e.g. 'none' with other options)
const ErrImproperTransportOption = cerr.Error("unknown improper option")

// ErrUnknownMuxerOption is returned when an unknown muxer has been specified
const ErrUnknownMuxerOption = cerr.Error("unknown muxer option")

// ErrImproperMuxerOption is returned when an improper muxer option format has been provided. It expects ['name', 'type']
const ErrImproperMuxerOption = cerr.Error("improper muxer option")

// ErrUnknownSecurityOption is returned when an unknown security option has been specified
const ErrUnknownSecurityOption = cerr.Error("unknown security option")

const pubsubTopic = "/libp2p/test/1.0.0"

type publisher struct {
	ps *pubsub.PubSub
}
