package client

import (
	"context"
	"time"

	"github.com/agencyenterprise/gossip-host/pkg/logger"
	pb "github.com/agencyenterprise/gossip-host/pkg/pb/publisher"

	"google.golang.org/grpc"
)

func Send(address string, msg *pb.Message) error {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		logger.Errorf("did not connect:\n%v", err)
		return err
	}
	defer conn.Close()

	c := pb.NewPublisherClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	r, err := c.PublishMessage(ctx, msg)
	if err != nil {
		logger.Errorf("could not publish message:\n %v", err)
		return err
	}

	logger.Infof("Message publish ok: %v", r.GetSuccess())
	if !r.GetSuccess() {
		return ErrPublishFailure
	}

	return nil
}
