package config

import (
	juno "github.com/desmos-labs/juno/types"
	"github.com/pelletier/go-toml"
)

type configToml struct {
	DatabaseConfig *DatabaseConfig `toml:"database"`
}

// ParseConfig allows to read the given file contents as a Config instance
func ParseConfig(fileContents []byte) (juno.Config, error) {
	// Parse the custom config
	var cfg configToml
	err := toml.Unmarshal(fileContents, &cfg)
	if err != nil {
		return nil, err
	}

	// Parse Juno config
	junoCfg, err := juno.DefaultConfigParser(fileContents)
	if err != nil {
		return nil, err
	}

	return NewConfig(
		junoCfg,
		NewDatabaseConfig(
			junoCfg.GetDatabaseConfig(),
			cfg.DatabaseConfig.StoreHistoricalData,
		),
	), err
}
