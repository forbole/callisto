package config

import (
	juno "github.com/desmos-labs/juno/types"
	"github.com/pelletier/go-toml"

	"github.com/forbole/bdjuno/types"
)

var _ juno.Config = &Config{}

// Config contains the data about the BDJuno configuration
type Config struct {
	juno.Config
	PricefeedConfig    *PricefeedConfig
	DistributionConfig *DistributionConfig
}

type tomlConfig struct {
	PricefeedConfig    *PricefeedConfig    `toml:"pricefeed"`
	DistributionConfig *DistributionConfig `toml:"distribution"`
}

// NewConfig returns a new Config instance
func NewConfig(junoConfig juno.Config, pricefeedConfig *PricefeedConfig, distributionConfig *DistributionConfig) *Config {
	return &Config{
		Config:             junoConfig,
		PricefeedConfig:    pricefeedConfig,
		DistributionConfig: distributionConfig,
	}
}

// Parser represents the method that should be used to parse a configuration
func Parser(fileContents []byte) (juno.Config, error) {
	junoCfg, err := juno.DefaultConfigParser(fileContents)
	if err != nil {
		return nil, err
	}

	var tomlCfg tomlConfig
	err = toml.Unmarshal(fileContents, &tomlCfg)
	if err != nil {
		return nil, err
	}

	return NewConfig(junoCfg, tomlCfg.PricefeedConfig, tomlCfg.DistributionConfig), nil
}

// GetRPCConfig implements juno.Config
func (c *Config) GetRPCConfig() juno.RPCConfig {
	return c.Config.GetRPCConfig()
}

// GetGrpcConfig implements juno.Config
func (c *Config) GetGrpcConfig() juno.GrpcConfig {
	return c.Config.GetGrpcConfig()
}

// GetCosmosConfig implements juno.Config
func (c *Config) GetCosmosConfig() juno.CosmosConfig {
	return c.Config.GetCosmosConfig()
}

// GetDatabaseConfig implements juno.Config
func (c *Config) GetDatabaseConfig() juno.DatabaseConfig {
	return c.Config.GetDatabaseConfig()
}

// GetLoggingConfig implements juno.Config
func (c *Config) GetLoggingConfig() juno.LoggingConfig {
	return c.Config.GetLoggingConfig()
}

// GetParsingConfig implements juno.Config
func (c *Config) GetParsingConfig() juno.ParsingConfig {
	return c.Config.GetParsingConfig()
}

// GetPruningConfig implements juno.Config
func (c *Config) GetPruningConfig() juno.PruningConfig {
	return c.Config.GetPruningConfig()
}

// GetTelemetryConfig implements juno.Config
func (c *Config) GetTelemetryConfig() juno.TelemetryConfig {
	return c.Config.GetTelemetryConfig()
}

// GetPricefeedConfig return the current PricefeedConfig
func (c *Config) GetPricefeedConfig() *PricefeedConfig {
	if c.PricefeedConfig == nil {
		return &PricefeedConfig{Tokens: nil}
	}
	return c.PricefeedConfig
}

// GetDistributionConfig return current distribution frequency
func (c *Config) GetDistributionConfig() *DistributionConfig {
	if c.DistributionConfig == nil {
		return &DistributionConfig{DistributionFrequency: 0}
	}
	return c.DistributionConfig
}

// --------------------------------------------------------------------------------------------------------------------

// PricefeedConfig contains the configuration about the pricefeed module
type PricefeedConfig struct {
	Tokens []types.Token `toml:"tokens"`
}

// GetTokens returns the list of tokens for which to get the prices
func (p *PricefeedConfig) GetTokens() []types.Token {
	return p.Tokens
}

// DistributionConfig contains the configuration about distribution frequency
type DistributionConfig struct {
	DistributionFrequency int64 `toml:"distribution_frequency"`
}

// GetDistributionFrequency returns distribution frequency int64 value
func (b *DistributionConfig) GetDistributionFrequency() int64 {
	return b.DistributionFrequency
}
