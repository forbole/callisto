package ccv

import (
	parsecmdtypes "github.com/forbole/juno/v4/cmd/parse/types"
	"github.com/spf13/cobra"
)

// NewCcvCmd returns the Cobra command that allows to fix all the things related to the ccv module
func NewCcvCmd(parseConfig *parsecmdtypes.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ccv",
		Short: "Fix things related to the x/ccv module",
	}

	cmd.AddCommand(
		ccvAllChainsCmd(parseConfig),
	)

	return cmd
}
