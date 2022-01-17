package blocks

import (
	"github.com/spf13/cobra"

	"github.com/forbole/juno/v2/cmd/parse"
)

// NewIBCTransfersCmd returns the Cobra command that allows to fix all the things related to IBC transfers
func NewIBCTransfersCmd(parseConfig *parse.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ibc-transfers",
		Short: "Fix things related to IBC transfers messages",
	}

	cmd.AddCommand(
		receivedCmd(parseConfig),
	)

	return cmd
}
