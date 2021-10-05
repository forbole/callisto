package distribution

import "gopkg.in/yaml.v3"

// Config contains the configuration about distribution frequency
type Config struct {
	RewardsFrequency int64 `yaml:"rewards_frequency,omitempty"`
}

// NewConfig returns a new Config instance
func NewConfig(frequency int64) *Config {
	return &Config{
		RewardsFrequency: frequency,
	}
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return NewConfig(100)
}

func ParseConfig(bz []byte) (*Config, error) {
	type T struct {
		Config *Config `yaml:"distribution"`
	}
	var cfg T
	err := yaml.Unmarshal(bz, &cfg)
	return cfg.Config, err
}
