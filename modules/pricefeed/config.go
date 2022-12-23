package pricefeed

import (
	"gopkg.in/yaml.v3"

	"github.com/forbole/bdjuno/v3/types"
)

// Config contains the configuration about the pricefeed module
type Config struct {
	Tokens []types.Token `yaml:"tokens"`
}

var PricefeedCfg *Config

// NewConfig returns a new Config instance
func NewConfig(tokens []types.Token) *Config {
	return &Config{
		Tokens: tokens,
	}
}

func ParseConfig(bz []byte) (*Config, error) {
	type T struct {
		Config *Config `yaml:"pricefeed"`
	}
	var cfg T
	err := yaml.Unmarshal(bz, &cfg)
	PricefeedCfg = cfg.Config

	return cfg.Config, err
}

func GetDenom() (denom string) {
	for _, token := range PricefeedCfg.Tokens {
		for _, unit := range token.Units {
			if unit.Exponent == 0 {
				denom = unit.Denom
			}
		}
	}
	return denom
}
