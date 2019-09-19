package config

import (
	"github.com/agencyenterprise/gossip-host/pkg/logger"
)

// Load reads the passed config file location and parses it into a config struct.
func Load(confLoc, listens, peers string) (*Config, error) {
	var (
		conf, defaults Config
	)

	if err := parseDefaults(&defaults); err != nil {
		logger.Errorf("err parsing defaults:\n%v", err)
	}

	if err := parseConfigFile(&conf, confLoc); err != nil {
		logger.Errorf("err parsing config file:\n%v", err)
		return nil, err
	}

	mergeDefaults(&conf, &defaults)
	parseListens(&conf, listens)
	parsePeers(&conf, peers)

	logger.Infof("configuration: %v", conf)
	return &conf, nil
}
