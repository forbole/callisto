package config

import (
	"git.ooo.ua/vipcoin/lib/vault"
	dbconfig "github.com/forbole/juno/v5/database/config"
	loggingconfig "github.com/forbole/juno/v5/logging/config"
	nodeconfig "github.com/forbole/juno/v5/node/config"
	"github.com/forbole/juno/v5/node/remote"
	parserconfig "github.com/forbole/juno/v5/parser/config"
	junoconf "github.com/forbole/juno/v5/types/config"
	"github.com/spf13/cobra"
)

// CheckVaultConfig - checking the ability to load the configuration from vault
func CheckVaultConfig(serviceName string, isLocal bool, cmdAndConf *cobra.Command) *cobra.Command {
	var (
		config = getDebugConfig()
		err    error
	)

	if !isLocal {
		config, err = loadFromVault(serviceName)
		if err != nil {
			return cmdAndConf
		}
	}

	cmdAndConf.PreRunE = func(_ *cobra.Command, _ []string) error {
		// Set the global configuration
		junoconf.Cfg = config
		return nil
	}

	return cmdAndConf
}

func loadFromVault(serviceName string) (junoconf.Config, error) {
	vaultClient, err := vault.NewClient(vault.NamespaceCubbyhole)
	if err != nil {
		return junoconf.Config{}, err
	}

	vaultData, err := vaultClient.Pull(serviceName)
	if err != nil {
		return junoconf.Config{}, err
	}

	return junoconf.DefaultConfigParser(vaultData)
}

// getDebugConfig - get the debug config
func getDebugConfig() junoconf.Config {
	return junoconf.Config{
		Chain: junoconf.ChainConfig{
			Bech32Prefix: "ovg",
			Modules: []string{
				"modules",
				"messages",
				"auth",
				"bank",
				"staking",
				"gov",
				"consensus",
				"mint",
				"slashing",
				"overgold",
			},
		},
		Node: nodeconfig.Config{
			Type: "remote",
			Details: &remote.Details{
				RPC: &remote.RPCConfig{
					ClientName:     "juno",
					Address:        "http://35.205.93.149:26657",
					MaxConnections: 40,
				},
				GRPC: &remote.GRPCConfig{
					Address:  "http://35.205.93.149:9090",
					Insecure: true,
				},
			},
		},
		Parser: parserconfig.Config{
			GenesisFilePath: "/bdjuno/genesis.json",
			Workers:         5,
			StartHeight:     0,
			AvgBlockTime:    nil,
			ParseNewBlocks:  false,
			ParseOldBlocks:  true,
			ParseGenesis:    true,
			FastSync:        false,
		},
		Database: dbconfig.Config{
			URL:                "postgresql://postgres:postgres@localhost:5432/juno?sslmode=disable",
			MaxOpenConnections: 10,
			MaxIdleConnections: 10,
			PartitionSize:      0,
			PartitionBatchSize: 0,
			SSLModeEnable:      "disable",
			SSLRootCert:        "",
			SSLCert:            "",
			SSLKey:             "",
		},
		Logging: loggingconfig.Config{
			LogLevel:  "debug",
			LogFormat: "text",
		},
	}
}
