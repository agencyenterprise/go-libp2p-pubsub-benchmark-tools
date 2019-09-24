package host

import (
	"context"
	"net"

	pb "github.com/agencyenterprise/gossip-host/internal/pb/publisher"
	"github.com/agencyenterprise/gossip-host/pkg/logger"

	"google.golang.org/grpc"
)

func New(messagePublisher func(msg *pb.Message) error) *Host {
	return &Host{
		&Server{
			messagePublisher,
		},
	}
}

func (h *Host) Listen(ctx context.Context, addr string) error {
	var lstnCfg net.ListenConfig
	lis, err := lstnCfg.Listen(ctx, "tcp", addr)
	if err != nil {
		logger.Errorf("failed to listen: %v", err)
		return err
	}
	defer lis.Close()

	s := grpc.NewServer()
	pb.RegisterPublisherServer(s, h.Server)
	if err := s.Serve(lis); err != nil {
		logger.Errorf("failed to serve: %v", err)
		return err
	}

	return nil
}
