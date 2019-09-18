package config

import (
	"encoding/json"
	"io/ioutil"

	"github.com/agencyenterprise/gossip-host/pkg/logger"
)

// Load reads the passed config file location and parses it into a config struct.
func Load(confLoc string) (*Config, error) {
	var conf Config

	file, err := ioutil.ReadFile(confLoc)
	if err != nil {
		logger.Errorf("err reading configuration file:%s\n%v", confLoc, err)
		return nil, err
	}

	if err = json.Unmarshal([]byte(file), &conf); err != nil {
		logger.Errorf("err unmarshaling config\n%v", err)
		return nil, err
	}

	return &conf, nil
}
