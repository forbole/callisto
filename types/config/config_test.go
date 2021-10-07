package config_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/forbole/bdjuno/types/config"
)

func TestConfigParser_WithPricefeedConfig(t *testing.T) {
	var configStr = `
[cosmos]
modules = []
prefix = "cosmos"

[pricefeed]
tokens = [
    { name = "Osmo", units = [{ denom = "uosmo", exponent = 0 }, { denom = "osmo", exponent = 6 }] },
    { name = "Ion", units = [{ denom = "ion", exponent = 0 }, { denom = "ion", exponent = 6 }] }
]

[rpc]
address = "http://localhost:26657"
client_name = "juno"

[grpc]
address = "localhost:9090"
insecure = true

[parsing]
fast_sync = true
listen_new_blocks = true
parse_genesis = true
parse_old_blocks = true
start_height = 1
workers = 1

[database]
host = "localhost"
max_idle_connections = 0
max_open_connections = 0
name = "database-name"
password = "password"
port = 5432
schema = "public"
ssl_mode = ""
user = "user"

[pruning]
interval = 10
keep_every = 500
keep_recent = 100

[logging]
format = "text"
level = "debug"`

	cfg, err := config.Parser([]byte(configStr))
	require.NoError(t, err)

	bdJunoCfg, ok := cfg.(*config.Config)
	require.True(t, ok)
	require.Len(t, bdJunoCfg.GetPricefeedConfig().GetTokens(), 2)
}

func TestConfigParser_WithoutPricefeedConfig(t *testing.T) {
	var configStr = `
[cosmos]
modules = []
prefix = "cosmos"

[rpc]
address = "http://localhost:26657"
client_name = "juno"

[grpc]
address = "localhost:9090"
insecure = true

[parsing]
fast_sync = true
listen_new_blocks = true
parse_genesis = true
parse_old_blocks = true
start_height = 1
workers = 1

[database]
host = "localhost"
max_idle_connections = 0
max_open_connections = 0
name = "database-name"
password = "password"
port = 5432
schema = "public"
ssl_mode = ""
user = "user"

[pruning]
interval = 10
keep_every = 500
keep_recent = 100

[logging]
format = "text"
level = "debug"`

	cfg, err := config.Parser([]byte(configStr))
	require.NoError(t, err)

	bdJunoCfg, ok := cfg.(*config.Config)
	require.True(t, ok)
	require.Len(t, bdJunoCfg.GetPricefeedConfig().GetTokens(), 0)
}
