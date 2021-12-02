package feegrant

import (
	"github.com/forbole/juno/v2/cmd/parse"
	"github.com/spf13/cobra"
)

// NewStakingCmd returns the Cobra command that allows to fix all the things related to the x/staking module
func NewFeegrantCmd(parseConfig *parse.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "feegrant",
		Short: "Fix things related to the x/feegrant module",
	}

	cmd.AddCommand(
		allowanceCmd(parseConfig),
	)

	return cmd
}
