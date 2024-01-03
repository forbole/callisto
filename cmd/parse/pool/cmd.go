package pool

import (
	parsecmdtypes "github.com/forbole/juno/v5/cmd/parse/types"
	"github.com/spf13/cobra"
)

// NewPoolCmd returns the Cobra command allowing to fix various things related to the x/pool module
func NewPoolCmd(parseConfig *parsecmdtypes.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pool",
		Short: "Fix things related to the x/pool module",
	}

	cmd.AddCommand(
		poolCmd(parseConfig),
	)

	return cmd
}
