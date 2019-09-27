package host

import pb "github.com/agencyenterprise/gossip-host/pkg/pb/publisher"

// Host listens on grpc
type Host struct {
	Server *Server
}

// Server is used to implement PublisherServer.
type Server struct {
	PblshMessage func(msg *pb.Message) error
}
