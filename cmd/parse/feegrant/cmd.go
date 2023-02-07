package feegrant

import (
	parsecmdtypes "github.com/forbole/juno/v4/cmd/parse/types"
	"github.com/spf13/cobra"
)

// NewFeegrantCmd returns the Cobra command that allows to fix all the things related to the x/feegrant module
func NewFeegrantCmd(parseConfig *parsecmdtypes.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "feegrant",
		Short: "Fix things related to the x/feegrant module",
	}

	cmd.AddCommand(
		allowanceCmd(parseConfig),
	)

	return cmd
}
