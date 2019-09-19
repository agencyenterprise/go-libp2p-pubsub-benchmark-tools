package config

import (
	"github.com/agencyenterprise/gossip-host/pkg/logger"
)

// Load reads the passed config file location and parses it into a config struct.
func Load(confLoc, listens, peers string) (*Config, error) {
	var conf Config

	if err := parseConfigFile(&conf, confLoc); err != nil {
		logger.Errorf("err parsing config file:\n%v", err)
		return nil, err
	}

	parseListens(&conf, listens)
	parsePeers(&conf, peers)

	return &conf, nil
}
