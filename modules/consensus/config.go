package consensus

import "gopkg.in/yaml.v3"

// Config contains the configuration about chain start height
type Config struct {
	StartHeight int64 `yaml:"start_height,omitempty"`
}

// NewConfig returns a new Config instance
func NewConfig(height int64) *Config {
	return &Config{
		StartHeight: height,
	}
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return NewConfig(1)
}

func ParseConfig(bz []byte) (*Config, error) {
	type T struct {
		Config *Config `yaml:"parsing"`
	}
	var cfg T
	err := yaml.Unmarshal(bz, &cfg)
	return cfg.Config, err
}
