package client

import (
	"context"
	"errors"
	"time"

	"github.com/agencyenterprise/gossip-host/pkg/logger"
	pb "github.com/agencyenterprise/gossip-host/pkg/pb/publisher"

	grpc "google.golang.org/grpc"
)

func Gossip(msgLoc, peers string, timeout int) error {
	msg, err := parseMessageFile(msgLoc)
	if err != nil || msg == nil {
		logger.Errorf("err parsing message file:\n%v", err)
		return err
	}
	logger.Infof("message is %s", msg.String())
	var failed = false

	peersArr := parsePeers(peers)
	for _, peer := range peersArr {
		// Set up a connection to the server.
		conn, err := grpc.Dial(peer, grpc.WithInsecure())
		if err != nil {
			logger.Errorf("did not connect to %s:\n%v", peer, err)
			failed = true
			continue
		}

		c := pb.NewPublisherClient(conn)

		// Contact the server and print out its response.
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)

		r, err := c.PublishMessage(ctx, msg)
		if err != nil {
			logger.Errorf("could not gossip message to %s:\n %v", peer, err)
			failed = true
			conn.Close()
			cancel()
			continue
		}

		logger.Infof("ok for %s: %v", peer, r.GetSuccess())
		if !r.GetSuccess() {
			failed = true
		}

		conn.Close()
		cancel()
	}

	if failed {
		return errors.New("gossip failed")
	}

	return nil
}

func CloseAll(peers string, timeout int) error {
	peersArr := parsePeers(peers)
	var failed = false

	for _, peer := range peersArr {
		// Set up a connection to the server.
		conn, err := grpc.Dial(peer, grpc.WithInsecure())
		if err != nil {
			logger.Errorf("did not connect to %s:\n%v", peer, err)
			failed = true
			continue
		}

		c := pb.NewPublisherClient(conn)

		// Contact the server and print out its response.
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)

		r, err := c.CloseAllPeerConnections(ctx, nil)
		if err != nil {
			logger.Errorf("err closing all peer connections for %s:\n%v", peer, err)
			failed = true
			conn.Close()
			cancel()
			continue
		}

		logger.Infof("ok for %s: %v", peer, r.GetSuccess())
		if !r.GetSuccess() {
			failed = true
		}

		conn.Close()
		cancel()
	}

	if failed {
		return errors.New("close all connections failed")
	}

	return nil
}

func ClosePeers(peers string, closePeers string, timeout int) error {
	peersArr := parsePeers(peers)
	closePeersList := parsePeers(closePeers)
	peersList := &pb.PeersList{
		Peers: closePeersList,
	}
	var failed = false

	for _, peer := range peersArr {
		// Set up a connection to the server.
		conn, err := grpc.Dial(peer, grpc.WithInsecure())
		if err != nil {
			logger.Errorf("did not connect to %s:\n%v", peer, err)
			failed = true
			continue
		}

		c := pb.NewPublisherClient(conn)

		// Contact the server and print out its response.
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)

		r, err := c.ClosePeerConnections(ctx, peersList)
		if err != nil {
			logger.Errorf("err closing all peer connections for %s:\n%v", peer, err)
			failed = true
			conn.Close()
			cancel()
			continue
		}

		logger.Infof("ok for %s: %v", peer, r.GetSuccess())
		if !r.GetSuccess() {
			failed = true
		}

		conn.Close()
		cancel()
	}

	if failed {
		return errors.New("close peer connections failed")
	}

	return nil
}

func OpenPeers(peers string, openPeers string, timeout int) error {
	peersArr := parsePeers(peers)
	openPeersList := parsePeers(openPeers)
	peersList := &pb.PeersList{
		Peers: openPeersList,
	}
	var failed = false

	for _, peer := range peersArr {
		// Set up a connection to the server.
		conn, err := grpc.Dial(peer, grpc.WithInsecure())
		if err != nil {
			logger.Errorf("did not connect to %s:\n%v", peer, err)
			failed = true
			continue
		}

		c := pb.NewPublisherClient(conn)

		// Contact the server and print out its response.
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)

		r, err := c.OpenPeersConnections(ctx, peersList)
		if err != nil {
			logger.Errorf("err opening peer connections for %s to %v:\n%v", peer, openPeers, err)
			failed = true
			conn.Close()
			cancel()
			continue
		}

		for _, peerConn := range r.GetPeerConnections() {
			if !peerConn.GetSuccess() {
				logger.Errorf("peer %v failed to connect to peer %v", peer, peerConn.GetPeer())
				failed = true
			}
		}

		conn.Close()
		cancel()
	}

	if failed {
		return errors.New("open peers failed")
	}

	return nil
}

func ListPeers(peers string, timeout int) error {
	peersArr := parsePeers(peers)
	var failed = false

	for _, peer := range peersArr {
		// Set up a connection to the server.
		conn, err := grpc.Dial(peer, grpc.WithInsecure())
		if err != nil {
			logger.Errorf("did not connect to %s:\n%v", peer, err)
			return err
		}

		c := pb.NewPublisherClient(conn)

		// Contact the server and print out its response.
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)

		r, err := c.ListConnectedPeers(ctx, nil)
		if err != nil {
			logger.Errorf("err closing all peer connections for peer %s:\n%v", peer, err)
			conn.Close()
			cancel()
			failed = true
			continue
		}

		for _, connPeers := range r.GetPeers() {
			logger.Infof("peers for %s:\n%v", peer, connPeers)
		}

		conn.Close()
		cancel()
	}

	if failed {
		return errors.New("open peers failed")
	}

	return nil
}

func IDs(peers string, timeout int) error {
	peersArr := parsePeers(peers)
	var failed = false

	for _, peer := range peersArr {
		// Set up a connection to the server.
		conn, err := grpc.Dial(peer, grpc.WithInsecure())
		if err != nil {
			logger.Errorf("did not connect to %s:\n%v", peer, err)
			failed = true
			continue
		}

		c := pb.NewPublisherClient(conn)

		// Contact the server and print out its response.
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)

		r, err := c.ID(ctx, nil)
		if err != nil {
			logger.Errorf("err getting id for %s:\n%v", peer, err)
			conn.Close()
			cancel()
			failed = true
			continue
		}

		logger.Infof("id for %s is %s", peer, r.GetID())
		conn.Close()
		cancel()
	}

	if failed {
		return errors.New("get id failed")
	}

	return nil
}

func Listens(peers string, timeout int) error {
	peersArr := parsePeers(peers)
	var failed = false

	for _, peer := range peersArr {
		// Set up a connection to the server.
		conn, err := grpc.Dial(peer, grpc.WithInsecure())
		if err != nil {
			logger.Errorf("did not connect to %s:\n%v", peer, err)
			failed = true
			continue
		}

		c := pb.NewPublisherClient(conn)

		// Contact the server and print out its response.
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)

		r, err := c.ListenAddresses(ctx, nil)
		if err != nil {
			logger.Errorf("err getting listens for %s:\n%v", peer, err)
			conn.Close()
			cancel()
			failed = true
			continue
		}

		logger.Infof("listens addresses for %s is %v", peer, r.GetAddresses())
		conn.Close()
		cancel()
	}

	if failed {
		return errors.New("get id failed")
	}

	return nil
}
