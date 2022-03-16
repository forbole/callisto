package staking

import (
	parsecmdtypes "github.com/forbole/juno/v3/cmd/parse/types"
	"github.com/spf13/cobra"
)

// NewStakingCmd returns the Cobra command that allows to fix all the things related to the x/staking module
func NewStakingCmd(parseConfig *parsecmdtypes.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "staking",
		Short: "Fix things related to the x/staking module",
	}

	cmd.AddCommand(
		validatorsCmd(parseConfig),
	)

	return cmd
}
