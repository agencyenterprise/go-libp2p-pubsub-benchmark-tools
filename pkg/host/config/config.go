package config

import (
	"github.com/agencyenterprise/gossip-host/pkg/logger"
)

// Load reads the passed config file location and parses it into a config struct.
func Load(confLoc, listens, rpcListen, peers string) (Config, error) {
	var (
		conf, defaults Config
		err            error
	)

	if err = parseDefaults(&defaults); err != nil {
		logger.Errorf("err parsing defaults:\n%v", err)
	}

	if err = parseConfigFile(&conf, confLoc); err != nil {
		logger.Errorf("err parsing config file:\n%v", err)
		return conf, err
	}

	if conf.Host.PrivPEM != "" {
		if err := loadAndSavePriv(&conf); err != nil {
			logger.Errorf("could not load private key:\n%v", err)
			return conf, err
		}
	}

	mergeDefaults(&conf, &defaults)
	parseListens(&conf, listens)
	parsePeers(&conf, peers)
	if rpcListen != "" {
		conf.Host.RPCAddress = rpcListen
	}

	return conf, nil
}
