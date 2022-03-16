package gov

import (
	parsecmdtypes "github.com/forbole/juno/v3/cmd/parse/types"
	"github.com/spf13/cobra"
)

// NewGovCmd returns the Cobra command allowing to fix various things related to the x/gov module
func NewGovCmd(parseConfig *parsecmdtypes.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gov",
		Short: "Fix things related to the x/gov module",
	}

	cmd.AddCommand(
		proposalCmd(parseConfig),
	)

	return cmd
}
