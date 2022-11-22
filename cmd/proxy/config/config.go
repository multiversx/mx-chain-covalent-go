package config

import (
	"io/ioutil"

	"github.com/pelletier/go-toml"
)

// Config holds the config for covalent proxy
type Config struct {
	Port                   uint32                 `toml:"port"`
	HyperBlockPath         string                 `toml:"hyperBlockPath"`
	HyperBlocksPath        string                 `toml:"hyperBlocksPath"`
	HyperBlocksBatchSize   uint32                 `toml:"hyperBlocksBatchSize"`
	ElrondProxyUrl         string                 `toml:"elrondProxyUrl"`
	RequestTimeOutSec      uint64                 `toml:"requestTimeOutSec"`
	HyperBlockQueryOptions HyperBlockQueryOptions `toml:"hyperBlockQueryOptions"`
}

// HyperBlockQueryOptions holds the hyper block query params options
type HyperBlockQueryOptions struct {
	WithLogs            bool   `toml:"withLogs"`
	WithAlteredAccounts bool   `toml:"withAlteredAccounts"`
	NotarizedAtSource   bool   `toml:"notarizedAtSource"`
	Tokens              string `toml:"tokens"`
}

// HyperBlocksQueryOptions holds the hyper blocks query params options
type HyperBlocksQueryOptions struct {
	QueryOptions HyperBlockQueryOptions
	BatchSize    uint64
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
