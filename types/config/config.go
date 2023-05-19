package config

import (
	initcmd "github.com/forbole/juno/v4/cmd/init"
	junoconfig "github.com/forbole/juno/v4/types/config"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// Config represents the BDJuno configuration
type Config struct {
	JunoConfig junoconfig.Config `yaml:"-,inline"`
}

// NewConfig returns a new Config instance
func NewConfig(junoCfg junoconfig.Config) Config {
	return Config{
		JunoConfig: junoCfg,
	}
}

// GetBytes implements WritableConfig
func (c Config) GetBytes() ([]byte, error) {
	return yaml.Marshal(&c)
}

// Creator represents a configuration creator
func Creator(_ *cobra.Command) initcmd.WritableConfig {
	return NewConfig(junoconfig.DefaultConfig())
}
