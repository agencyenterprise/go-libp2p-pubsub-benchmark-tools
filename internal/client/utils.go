package client

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	pb "github.com/agencyenterprise/gossip-host/internal/pb/publisher"
	"github.com/agencyenterprise/gossip-host/pkg/logger"
)

func parseMessageFile(loc string) (*pb.Message, error) {
	b, err := ioutil.ReadFile(loc)
	if err != nil {
		logger.Errorf("err reading file %s:\n%v", err)
		return nil, err
	}

	var msg pb.Message
	if err = json.Unmarshal(b, &msg); err != nil {
		logger.Errorf("err unmarshaling message:\n%v", err)
		return nil, err
	}

	return &msg, nil
}

func parsePeers(peers string) []string {
	peersArr := strings.Split(peers, ",")
	for idx := range peersArr {
		peersArr[idx] = strings.TrimSpace(peersArr[idx])
	}
	return peersArr
}
