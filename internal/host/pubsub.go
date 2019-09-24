package host

import (
	"context"
	"time"

	"github.com/agencyenterprise/gossip-host/pkg/logger"

	pb "github.com/agencyenterprise/gossip-host/internal/pb/publisher"
	peer "github.com/libp2p/go-libp2p-peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

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
		logger.Infof("Pubsub message received: %v,%v,%v,%d,%d", hostID, nxt.GetFrom(), msg.Id, time.Now().UnixNano(), msg.Sequence)
	}
}

func (publisher *publisher) publish(msg *pb.Message) error {
	var b []byte

	bs, err := msg.XXX_Marshal(b, true)
	if err != nil {
		logger.Errorf("err marshaling message:\n%v", err)
		return err
	}

	return publisher.ps.Publish(pubsubTopic, bs)
}
