package client

import (
	"github.com/agencyenterprise/gossip-host/pkg/cerr"
)

// ErrPublishFailure is returned when the peer couldn't publish
const ErrPublishFailure = cerr.Error("peer could not publish message")
