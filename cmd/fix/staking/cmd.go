package staking

import (
	"github.com/forbole/juno/v2/cmd/parse"
	"github.com/spf13/cobra"
)

// NewStakingCmd returns the Cobra command that allows to fix all the things related to the x/staking module
func NewStakingCmd(parseConfig *parse.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "staking",
		Short: "Fix things related to the x/staking module",
	}

	cmd.AddCommand(
		validatorsCmd(parseConfig),
		slashesCmd(parseConfig),
	)

	return cmd
}
