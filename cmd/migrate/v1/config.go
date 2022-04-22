package v1

import (
	"github.com/pelletier/go-toml"
)

type TomlConfig struct {
	PricefeedConfig    *PricefeedConfig    `toml:"pricefeed"`
	DistributionConfig *DistributionConfig `toml:"distribution"`
}

func ParseConfig(bz []byte) (TomlConfig, error) {
	var tomlCfg TomlConfig
	err := toml.Unmarshal(bz, &tomlCfg)
	if err != nil {
		return TomlConfig{}, err
	}

	return TomlConfig{
		PricefeedConfig:    tomlCfg.PricefeedConfig,
		DistributionConfig: tomlCfg.DistributionConfig,
	}, nil
}

type PricefeedConfig struct {
	Tokens []Token `toml:"tokens"`
}

type Token struct {
	Name  string      `toml:"name"`
	Units []TokenUnit `toml:"units"`
}

type TokenUnit struct {
	Denom    string   `toml:"denom"`
	Exponent int      `toml:"exponent"`
	Aliases  []string `toml:"aliases"`
	PriceID  string   `toml:"price_id"`
}

type DistributionConfig struct {
	DistributionFrequency int64 `toml:"distribution_frequency"`
}
