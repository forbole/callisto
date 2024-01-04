package wormhole

import (
	parsecmdtypes "github.com/forbole/juno/v5/cmd/parse/types"
	"github.com/spf13/cobra"
)

// NewWormholeCmd returns the Cobra command that allows to fix all the things related to the x/wormhole module
func NewWormholeCmd(parseConfig *parsecmdtypes.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "wormhole",
		Short: "Fix things related to the x/wormhole module",
	}

	cmd.AddCommand(
		updateGuardianSetCmd(parseConfig),
		updateGuardianValidatorsCmd(parseConfig),
	)

	return cmd
}
