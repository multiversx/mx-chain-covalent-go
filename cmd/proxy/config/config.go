package config

import (
	"io/ioutil"

	"github.com/pelletier/go-toml"
)

// Config holds the config for covalent proxy
type Config struct {
	Port              uint32 `toml:"Port"`
	HyperBlockPath    string `toml:"HyperBlockPath"`
	ElrondProxyUrl    string `toml:"ElrondProxyUrl"`
	RequestTimeOutSec uint64 `toml:"RequestTimeOutSec"`
}

// LoadConfig will load the Config from the provided file
func LoadConfig(tomlFile string) (*Config, error) {
	tomlBytes, err := ioutil.ReadFile(tomlFile)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = toml.Unmarshal(tomlBytes, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
