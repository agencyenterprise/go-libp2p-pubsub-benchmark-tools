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

func loadDefaultBox() *packr.Box {
	return packr.New("defaults", defaultsLoc)
}

func loadDefaultConfig(box *packr.Box) ([]byte, error) {
	// Get the string representation of a file, or an error if it doesn't exist:
	return box.Find(defaultConfigName)
}

func parseDefaults(conf *Config) error {
	box := loadDefaultBox()

	defaultConfig, err := loadDefaultConfig(box)
	if err != nil {
		logger.Errorf("err loading default config:\n%v", err)
		return err
	}

	if err := json.Unmarshal(defaultConfig, conf); err != nil {
		logger.Errorf("err unmarshaling config\n%v", err)
		return err
	}

	return nil
}

// note: this could panic!
func mergeDefaults(conf, defaults *Config) {
	// subnet
	if conf.Subnet.NumHosts == 0 {
		conf.Subnet.NumHosts = defaults.Subnet.NumHosts
	}
	if conf.Subnet.PubsubCIDR == "" {
		conf.Subnet.PubsubCIDR = defaults.Subnet.PubsubCIDR
	}
	if len(conf.Subnet.PubsubPortRange) == 0 {
		conf.Subnet.PubsubPortRange = defaults.Subnet.PubsubPortRange
	}
	if conf.Subnet.RPCCIDR == "" {
		conf.Subnet.RPCCIDR = defaults.Subnet.RPCCIDR
	}
	if len(conf.Subnet.RPCPortRange) == 0 {
		conf.Subnet.RPCPortRange = defaults.Subnet.RPCPortRange
	}

	// host
	//if conf.Host.Priv == nil {
	//conf.Host.Priv = defaults.Host.Priv
	//}
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
