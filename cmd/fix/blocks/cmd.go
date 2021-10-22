package blocks

import (
	"github.com/forbole/juno/v2/cmd/parse"
	"github.com/spf13/cobra"
)

// NewBlocksCmd returns the Cobra command that allows to fix all the things related to blocks
func NewBlocksCmd(parseConfig *parse.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "blocks",
		Short: "Fix things related to blocks",
	}

	cmd.AddCommand(
		blocksCmd(parseConfig),
	)

	return cmd
}
