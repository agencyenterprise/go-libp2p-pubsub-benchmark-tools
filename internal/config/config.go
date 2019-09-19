package config

import (
	"github.com/agencyenterprise/gossip-host/pkg/logger"

	"github.com/spf13/viper"
)

// Load reads the passed config file location and parses it into a config struct.
func Load(confLoc, listens, peers string) (*Config, error) {
	var conf Config

	viper.SetConfigName(trimExtension(confLoc))

	if err := viper.ReadInConfig(); err != nil {
		logger.Errorf("err reading configuration file:%s\n%v", confLoc, err)
		return nil, err
	}

	if err := viper.Unmarshal(&conf); err != nil {
		logger.Errorf("err unmarshaling config\n%v", err)
		return nil, err
	}
	logger.Infof("config: %v", conf)

	/*
		file, err := ioutil.ReadFile(confLoc)
		if err != nil {
			logger.Errorf("err reading configuration file:%s\n%v", confLoc, err)
			return nil, err
		}

		if err = json.Unmarshal([]byte(file), &conf); err != nil {
			logger.Errorf("err unmarshaling config\n%v", err)
			return nil, err
		}
	*/

	return &conf, nil
}
