package host

import (
	"context"
	"time"

	"github.com/agencyenterprise/gossip-host/pkg/logger"

	pb "github.com/agencyenterprise/gossip-host/internal/host/pb"
	peer "github.com/libp2p/go-libp2p-peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

const pubsubTopic = "/libp2p/test/1.0.0"

func pubsubHandler(ctx context.Context, hostID peer.ID, sub *pubsub.Subscription) {
	for {
		nxt, err := sub.Next(ctx)
		if err != nil {
			logger.Errorf("err reading next:\n%v", err)
			continue
		}

		msg := &pb.Message{}
		if err = nxt.Unmarshal(nxt.Data); err != nil {
			logger.Errorf("err unmarshaling next message:\n%v", err)
			continue
		}

		// TODO: how to increment sequence before sending out?
		logger.Infof("%v,%v,%v,%d,%d", hostID, nxt.GetFrom(), msg.Id, time.Now().UnixNano(), msg.Sequence)
	}
}
