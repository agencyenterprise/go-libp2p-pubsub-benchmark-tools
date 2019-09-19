package config

import (
	"encoding/json"
	"path/filepath"
	"strings"

	"github.com/agencyenterprise/gossip-host/pkg/logger"

	"github.com/gobuffalo/packr/v2"
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
	v := viper.New()

	v.SetConfigName(trimExtension(confLoc))
	v.AddConfigPath(".")

	if err := v.ReadInConfig(); err != nil {
		logger.Errorf("err reading configuration file: %s\n%v", confLoc, err)
		return err
	}

	if err := v.Unmarshal(conf); err != nil {
		logger.Errorf("err unmarshaling config\n%v", err)
		return err
	}

	return nil
}

func loadDefaults() ([]byte, error) {
	// set up a new box by giving it a name and an optional (relative) path to a folder on disk:
	box := packr.New("defaults", defaultsLoc)

	// Get the string representation of a file, or an error if it doesn't exist:
	return box.Find(defaultsName)
}

func parseDefaults(conf *Config) error {
	defaults, err := loadDefaults()
	if err != nil {
		logger.Errorf("err loading default config:\n%v", err)
		return err
	}

	if err := json.Unmarshal(defaults, conf); err != nil {
		logger.Errorf("err unmarshaling config\n%v", err)
		return err
	}

	return nil
}

// note: this could panic!
func mergeDefaults(conf, defaults *Config) {
	if len(conf.Host.Listens) == 0 {
		logger.Info("zero length listens")
		conf.Host.Listens = defaults.Host.Listens
	}
	if len(conf.Host.Peers) == 0 {
		conf.Host.Peers = defaults.Host.Peers
	}
	if len(conf.Host.Transports) == 0 {
		conf.Host.Transports = defaults.Host.Transports
	}
	if len(conf.Host.Muxers) == 0 {
		conf.Host.Muxers = defaults.Host.Muxers
	}
	if conf.Host.Security == "" {
		conf.Host.Security = defaults.Host.Security
	}
}
