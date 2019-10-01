package config

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/agencyenterprise/gossip-host/pkg/logger"

	// TODO: wait for pr merge and go back to lcrypto
	//acrypto "github.com/adam-hanna/go-libp2p-core/crypto"
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

func loadDefaultBox() *packr.Box {
	return packr.New("defaults", defaultsLoc)
}

func loadDefaultConfig(box *packr.Box) ([]byte, error) {
	// Get the string representation of a file, or an error if it doesn't exist:
	return box.Find(defaultConfigName)
}

func loadDefaultPriv(box *packr.Box) ([]byte, error) {
	// Get the string representation of a file, or an error if it doesn't exist:
	return box.Find(defaultPEMName)
}

// TODO: wait for pr merge and go back to lcrypto
/*
func loadAndSavePriv(conf *Config) error {
	privB, err := loadPriv(conf.Host.PrivPEM)
	if err != nil {
		logger.Errorf("err loading private key file:\n%v", err)
		return err
	}

	priv, err := parsePrivateKey(privB)
	if err != nil {
		logger.Errorf("err parsing private key:\n%v", err)
		return err
	}

	conf.Host.Priv = priv

	return nil
}
*/

func loadPriv(loc string) ([]byte, error) {
	privateKeyFile, err := os.Open(loc)
	if err != nil {
		logger.Errorf("err loading private key pem file: %s\n%v", loc, err)
		return nil, err
	}
	defer privateKeyFile.Close()

	pemfileinfo, err := privateKeyFile.Stat()
	if err != nil {
		logger.Errorf("err statting private key file:\n%v", err)
		return nil, err
	}
	var size int64 = pemfileinfo.Size()
	pembytes := make([]byte, size)

	buffer := bufio.NewReader(privateKeyFile)
	_, err = buffer.Read(pembytes)
	return []byte(pembytes), err
}

// TODO: waiting on PR merge to lcrypto
/*
func parseDefaultPriv(box *packr.Box) (lcrypto.PrivKey, error) {
	defaultPriv, err := loadDefaultPriv(box)
	if err != nil {
		logger.Errorf("err loading default private key:\n%v", err)
		return nil, err
	}

	return parsePrivateKey(defaultPriv)
}
*/

// TODO: waiting on PR merge to lcrypto
/*
func parsePrivateKey(privB []byte) (lcrypto.PrivKey, error) {
	data, _ := pem.Decode(privB)
	if data == nil {
		logger.Error("err decoding default PEM file. Nil data block")
		return nil, errors.New("err decoding default PEM file")
	}

	cPriv, err := x509.ParsePKCS8PrivateKey(data.Bytes)
	if err != nil {
		logger.Errorf("err parsing private key bytes:\n%v", err)
		return nil, err
	}

	// TODO: remove ASAP
	priv, _, err := acrypto.KeyPairFromKey(cPriv)
	if err != nil {
		logger.Errorf("err generating lcrypto priv key:\n%v", err)
		return nil, err
	}

	return priv.(lcrypto.PrivKey), nil
}
*/

// note: this could panic!
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

	// TODO: waiting on PR merge to lcrypto
	/*
		priv, err := parseDefaultPriv(box)
		if err != nil {
			logger.Errorf("err parsing default private key:\n%v", err)
			return err
		}
		conf.Host.Priv = priv
	*/

	return nil
}

// note: this could panic!
func mergeDefaults(conf, defaults *Config) {
	if conf.Host.Priv == nil {
		conf.Host.Priv = defaults.Host.Priv
	}
	if len(conf.Host.Listens) == 0 {
		conf.Host.Listens = defaults.Host.Listens
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
	if conf.Host.RPCAddress == "" {
		conf.Host.RPCAddress = defaults.Host.RPCAddress
	}
	if conf.Host.LoggerLocation == "" {
		conf.Host.LoggerLocation = defaults.Host.LoggerLocation
	}
}
