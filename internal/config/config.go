package config

import (
	"crypto/rand"

	"github.com/agencyenterprise/gossip-host/pkg/logger"
	lcrypto "github.com/libp2p/go-libp2p-core/crypto"
)

// Load reads the passed config file location and parses it into a config struct.
func Load(confLoc, listens, peers string) (*Config, error) {
	var (
		conf, defaults Config
		err            error
	)

	if err = parseDefaults(&defaults); err != nil {
		logger.Errorf("err parsing defaults:\n%v", err)
	}

	if err = parseConfigFile(&conf, confLoc); err != nil {
		logger.Errorf("err parsing config file:\n%v", err)
		return nil, err
	}

	// TODO: wait for pr merge and go back to lcrypto
	/*
		if conf.Host.PrivPEM != "" {
			if err := loadAndSavePriv(&conf); err != nil {
				logger.Errorf("could not load private key:\n%v", err)
				return nil, err
			}
		}
	*/
	conf.Host.Priv, _, err = lcrypto.GenerateECDSAKeyPair(rand.Reader)
	if err != nil {
		logger.Errorf("err generating ecdsa key pair:\n%v", err)
		return nil, err
	}

	mergeDefaults(&conf, &defaults)
	parseListens(&conf, listens)
	parsePeers(&conf, peers)

	logger.Infof("configuration: %v", conf)
	return &conf, nil
}
