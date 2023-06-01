package stakers

import (
	parsecmdtypes "github.com/forbole/juno/v4/cmd/parse/types"
	"github.com/spf13/cobra"
)

// NewStakersCmd returns the Cobra command allowing to fix various things related to the x/stakers module
func NewStakersCmd(parseConfig *parsecmdtypes.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stakers",
		Short: "Fix things related to the x/stakers module",
	}

	cmd.AddCommand(
		stakersCmd(parseConfig),
	)

	return cmd
}
