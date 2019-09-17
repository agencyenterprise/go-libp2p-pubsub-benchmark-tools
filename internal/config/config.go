package config

import (
	"github.com/agencyenterprise/gossip-host/pkg/logger"

	"github.com/BurntSushi/toml"
)

// Load reads the passed config file location and parses it into a config struct.
func Load(confLoc string) (*Config, error) {
	var config Config

	if _, err := toml.DecodeFile(confLoc, &config); err != nil {
		logger.Errorf("err unmarshaling config\n%v", err)
		return nil, err
	}

	return &config, nil
}
