package config

import (
	juno "github.com/desmos-labs/juno/types"
)

var (
	_ juno.Config = &Config{}
)

// Config contains the configuration data for the parser
type Config struct {
	juno.Config
	databaseConfig *DatabaseConfig
}

// NewConfig allows to build a new Config instance
func NewConfig(junoCfg juno.Config, databaseCfg *DatabaseConfig) juno.Config {
	return &Config{
		Config:         junoCfg,
		databaseConfig: databaseCfg,
	}
}

func (c *Config) GetDatabaseConfig() juno.DatabaseConfig {
	return c.databaseConfig
}

// --------------------------------------------------------------------------------------------------------------------

var _ juno.DatabaseConfig = &DatabaseConfig{}

// DatabaseConfig extends juno.databaseConfig allowing to specify whether or not to store historical data
type DatabaseConfig struct {
	juno.DatabaseConfig
	StoreHistoricalData bool `toml:"store_historical_data"`
}

// NewDatabaseConfig allows to build a new DatabaseConfig instance
func NewDatabaseConfig(junoDbCfg juno.DatabaseConfig, storeHistoricalData bool) *DatabaseConfig {
	return &DatabaseConfig{
		DatabaseConfig:      junoDbCfg,
		StoreHistoricalData: storeHistoricalData,
	}
}

// ShouldStoreHistoricalData tells whether or not to persist historical data
func (d *DatabaseConfig) ShouldStoreHistoricalData() bool {
	return d.StoreHistoricalData
}
