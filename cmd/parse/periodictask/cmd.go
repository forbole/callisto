package periodictask

import (
	parsecmdtypes "github.com/forbole/juno/v3/cmd/parse/types"
	"github.com/spf13/cobra"
)

// NewPeriodicTaskCmd returns the Cobra command allowing to fix various things related to the periodic operations
func NewPeriodicTaskCmd(parseConfig *parsecmdtypes.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "periodic-task",
		Short: "Refresh data that depend on periodic operations",
	}

	cmd.AddCommand(
		stakingCmd(parseConfig),
	)

	return cmd
}
