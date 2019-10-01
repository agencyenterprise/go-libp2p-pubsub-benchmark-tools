package host

import (
	"context"

	pb "github.com/agencyenterprise/gossip-host/pkg/pb/publisher"
	"github.com/agencyenterprise/gossip-host/pkg/logger"
	"github.com/davecgh/go-spew/spew"
)

// PublishMessage implements PublisherServer
func (s *Server) PublishMessage(ctx context.Context, in *pb.Message) (*pb.PublishReply, error) {
	logger.Info("received rpc message; will now publish to subscribers")
	spew.Dump(in)

	var err error
	if err = s.PblshMessage(in); err != nil {
		logger.Errorf("err publishing message:\n%v", err)
	}

	return &pb.PublishReply{
		MsgId:   in.Id,
		Success: err == nil,
	}, nil
}
