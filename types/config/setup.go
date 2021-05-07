package config

import (
	initcmd "github.com/desmos-labs/juno/cmd/init"
	juno "github.com/desmos-labs/juno/types"
	"github.com/spf13/cobra"
)

const (
	flagDataType = "Application-data-type"
)

// SetupConfigFlags implements initcmd.ConfigFlagSetup
func SetupConfigFlags(cmd *cobra.Command) {
	cmd.Flags().String(flagDataType, DataTypeUpdated, `Set the data type that will be persisted after parsing the chain. 
The two possible values are the following:

1. %s if you want to persist only the most up-to-date version of the data. 
   This means for example that you will not be able to query historic account balances or delegations
   This should be used when running the parser for an explorer (eg. BigDipper).

2. %s if you want to persist only historic version of your data.
   This will store the historical delegations, account balances and token prices.
   Note that by running the parsed with this option will cause your database to become very large very quickly.
   This should be used when running the parser for a history-based utility Application (eg. ForboleX).`)
}

// CreateConfig implements initcmd.ConfigCreator
func CreateConfig(cmd *cobra.Command) juno.Config {
	junoCfg := initcmd.DefaultConfigCreator(cmd)

	dataType, _ := cmd.Flags().GetString(flagDataType)

	return NewConfig(
		junoCfg,
		NewApplicationConfig(dataType),
	)
}
