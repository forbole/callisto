package top_accounts

import (
	parsecmdtypes "github.com/forbole/juno/v3/cmd/parse/types"
	"github.com/spf13/cobra"
)

func NewTopAccountsCmd(parseConfig *parsecmdtypes.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use: "top-accounts",
	}

	cmd.AddCommand(
		allCmd(parseConfig),
	)

	return cmd
}
