package config

import (
	juno "github.com/desmos-labs/juno/types"
)

const (
	// DataTypeUpdated represents the data type that should be used when you want to persist only the most up-to-date
	// version of data.
	DataTypeUpdated = "updated"

	// DataTypeHistoric represents the data type that should be used when you want to persist the historic verion of
	// the chain data (eg. historic delegations, balances, and so on).
	DataTypeHistoric = "historic"
)

var (
	_ juno.Config = &Config{}
)

// Config contains the configuration data for the parser
type Config struct {
	juno.Config
	Application *ApplicationConfig `toml:"application"`
}

// NewConfig allows to build a new Config instance
func NewConfig(junoCfg juno.Config, applicationCfg *ApplicationConfig) *Config {
	return &Config{
		Config:      junoCfg,
		Application: applicationCfg,
	}
}

// GetDataType returns the type of data that should be persisted
func (c *Config) GetDataType() string {
	return c.Application.DataType
}

// ApplicationConfig contains the configuration telling what kind of Application to parse the data for.
type ApplicationConfig struct {
	DataType string `toml:"data_type"`
}

// NewApplicationConfig allows to build a new ApplicationConfig instance
func NewApplicationConfig(dataType string) *ApplicationConfig {
	return &ApplicationConfig{
		DataType: dataType,
	}
}
