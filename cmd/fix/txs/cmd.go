package txs

import (
	"github.com/forbole/juno/v2/cmd/parse"
	"github.com/spf13/cobra"
)

// NewTxsCmd returns the Cobra command that allows to fix all the things related to transactions
func NewTxsCmd(parseConfig *parse.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "txs",
		Short: "Fix things related to transactions",
	}

	cmd.AddCommand(
		transactionsCmd(parseConfig),
	)

	return cmd
}
