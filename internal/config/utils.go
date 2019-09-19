package config

import (
	"path/filepath"
	"strings"

	"github.com/agencyenterprise/gossip-host/pkg/logger"

	"github.com/spf13/viper"
)

func trimExtension(filename string) string {
	return strings.TrimSuffix(filename, filepath.Ext(filename))
}

func parseListens(conf *Config, listens string) {
	if listens != "" {
		listensArr := strings.Split(listens, ",")
		for idx := range listensArr {
			listensArr[idx] = strings.TrimSpace(listensArr[idx])
		}
		conf.Host.Listens = listensArr
	}
}

func parsePeers(conf *Config, peers string) {
	if peers != "" {
		peersArr := strings.Split(peers, ",")
		for idx := range peersArr {
			peersArr[idx] = strings.TrimSpace(peersArr[idx])
		}
		conf.Host.Peers = peersArr
	}
}

func parseConfigFile(conf *Config, confLoc string) error {
	viper.SetConfigName(trimExtension(confLoc))
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		logger.Errorf("err reading configuration file: %s\n%v", confLoc, err)
		return err
	}

	if err := viper.Unmarshal(&conf); err != nil {
		logger.Errorf("err unmarshaling config\n%v", err)
		return err
	}

	return nil
}
