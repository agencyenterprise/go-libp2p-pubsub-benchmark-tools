package host

import (
	"context"

	pb "github.com/agencyenterprise/gossip-host/internal/pb/publisher"
	"github.com/agencyenterprise/gossip-host/pkg/logger"
	"github.com/davecgh/go-spew/spew"
)

// PublishMessage implements PublisherServer
func (s *Server) PublishMessage(ctx context.Context, in *pb.Message) (*pb.PublishReply, error) {
	logger.Info("Received message")
	logger.Infof("%v", s)
	spew.Dump(in)

	var err error
	if err = s.PblshMessage(in); err != nil {
		logger.Errorf("err publishing message:\n%v", err)
	}

	// TODO: send message!
	return &pb.PublishReply{
		MsgId:   in.Id,
		Success: err == nil,
	}, nil
}
