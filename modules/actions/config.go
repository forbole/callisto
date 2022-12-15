package actions

import (
	"github.com/forbole/juno/v4/node/remote"
	"gopkg.in/yaml.v3"
)

// Config contains the configuration about the actions module
type Config struct {
	Port uint            `yaml:"port"`
	Node *remote.Details `yaml:"node,omitempty"`
}

// NewConfig returns a new Config instance
func NewConfig(port uint, remoteDetails *remote.Details) *Config {
	return &Config{
		Port: port,
		Node: remoteDetails,
	}
}

// DefaultConfig returns the default configuration
func DefaultConfig() Config {
	return Config{
		Port: 3000,
		Node: nil,
	}
}

func ParseConfig(bz []byte) (*Config, error) {
	type T struct {
		Config *Config `yaml:"actions"`
	}
	var cfg T
	err := yaml.Unmarshal(bz, &cfg)
	return cfg.Config, err
}
