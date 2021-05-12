package config_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/forbole/bdjuno/types/config"
)

func TestParseConfig(t *testing.T) {
	data := `
[database]
  store_historical_data = true
  host = "localhost"
  name = "juno"
  password = "password"
  port = 5432
  schema = "public"
  ssl_mode = ""
  user = "user"
`

	cfg, err := config.ParseConfig([]byte(data))
	require.NoError(t, err)

	dbConfig, ok := cfg.GetDatabaseConfig().(*config.DatabaseConfig)
	require.True(t, ok)

	require.Equal(t, true, dbConfig.ShouldStoreHistoricalData())
}
