package distribution

import "gopkg.in/yaml.v3"

// Config contains the configuration about distribution frequency
type Config struct {
	DistributionFrequency int64 `yaml:"distribution_frequency"`
}

// GetDistributionFrequency returns distribution frequency int64 value
func (b *Config) GetDistributionFrequency() int64 {
	return b.DistributionFrequency
}

func ParseConfig(bz []byte) (*Config, error) {
	type T struct {
		Config *Config `yaml:"distribution"`
	}
	var cfg T
	err := yaml.Unmarshal(bz, &cfg)
	return cfg.Config, err
}
