package refresh

import (
	parsecmdtypes "github.com/forbole/juno/v3/cmd/parse/types"
	"github.com/spf13/cobra"
)

// NewRefreshCmd returns the Cobra command allowing to fix various things related to the periodic operations
func NewRefreshCmd(parseConfig *parsecmdtypes.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "refresh",
		Short: "Refresh data that depend on periodic operations",
	}

	cmd.AddCommand(
		bankCmd(parseConfig),
	)

	return cmd
}
