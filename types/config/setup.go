package config

import (
	initcmd "github.com/desmos-labs/juno/cmd/init"
	juno "github.com/desmos-labs/juno/types"
	"github.com/spf13/cobra"
)

const (
	flagDatabaseStoreHistoricData = "database-store-historic-data"
)

// SetupConfigFlags implements initcmd.ConfigFlagSetup
func SetupConfigFlags(cmd *cobra.Command) {
	cmd.Flags().Bool(flagDatabaseStoreHistoricData, false,
		"Whether or not to persist historic data inside the data")
}

// CreateConfig implements initcmd.ConfigCreator
func CreateConfig(cmd *cobra.Command) juno.Config {
	junoCfg := initcmd.DefaultConfigCreator(cmd)

	storeHistoricData, _ := cmd.Flags().GetBool(flagDatabaseStoreHistoricData)

	return NewConfig(
		junoCfg,
		NewDatabaseConfig(
			junoCfg.GetDatabaseConfig(),
			storeHistoricData,
		),
	)
}
