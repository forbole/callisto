package types

import (
	juno "github.com/desmos-labs/juno/types"
	"github.com/pelletier/go-toml"
)

const (
	ApplicationTypeExplorer = "explorer"
	ApplicationTypeUtility  = "utility"
)

var (
	_ juno.Config = &Config{}
)

// ParseConfig allows to read the given file contents as a Config instance
func ParseConfig(fileContents []byte) (juno.Config, error) {
	var cfg Config
	err := toml.Unmarshal(fileContents, &cfg)
	return &cfg, err
}

// applicationConfig contains the configuration telling what kind of application to parse the data for.
type applicationConfig struct {
	Type string
}

// Config contains the configuration data for the parser
type Config struct {
	juno.Config
	application *applicationConfig `toml:"persistence"`
}

// GetApplicationType returns the type of the application specified inside the configuration
func (c *Config) GetApplicationType() string {
	return c.application.Type
}
