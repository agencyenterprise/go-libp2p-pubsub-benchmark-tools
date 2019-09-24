package host

import (
	"context"
	"encoding/binary"
	"time"

	"github.com/agencyenterprise/gossip-host/pkg/logger"
	"github.com/davecgh/go-spew/spew"

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

		// TODO: fix this. It isn't working
		msg := &pb.Message{}
		if err = nxt.Unmarshal(nxt.Data); err != nil {
			logger.Errorf("err unmarshaling next message:\n%v", err)
			continue
		}
		spew.Dump(msg)

		// TODO: how to increment sequence before sending out?
		// note: what is nxt.GetFrom()? It doesn't seem to be the sender id but it doesn't match the original host peer.Id?
		logger.Infof("Pubsub message received: %v,%v,%v,%v,%d,%d", hostID, nxt.GetFrom(), msg.GetId(), binary.BigEndian.Uint64(nxt.GetSeqno()), time.Now().UnixNano(), msg.Sequence)
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
