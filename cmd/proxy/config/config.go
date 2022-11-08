package config

import (
	"io/ioutil"

	"github.com/pelletier/go-toml"
)

// Config holds the config for covalent proxy
type Config struct {
	Port                   uint32                 `toml:"port"`
	HyperBlockPath         string                 `toml:"hyperBlockPath"`
	ElrondProxyUrl         string                 `toml:"elrondProxyUrl"`
	RequestTimeOutSec      uint64                 `toml:"requestTimeOutSec"`
	HyperBlockQueryOptions HyperBlockQueryOptions `toml:"hyperBlockQueryOptions"`
}

// HyperBlockQueryOptions holds the hyper block query params options
type HyperBlockQueryOptions struct {
	WithLogs            bool   `toml:"withLogs"`
	WithAlteredAccounts bool   `toml:"withLogs"`
	NotarizedAtSource   bool   `toml:"notarizedAtSource"`
	Tokens              string `toml:"tokens"`
	WithMetaData        bool   `toml:"withMetaData"`
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
