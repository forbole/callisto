package wormhole

import (
	"fmt"

	parsecmdtypes "github.com/forbole/juno/v5/cmd/parse/types"
	"github.com/forbole/juno/v5/types/config"
	"github.com/spf13/cobra"

	"github.com/forbole/bdjuno/v4/database"
	modulestypes "github.com/forbole/bdjuno/v4/modules/types"
	"github.com/forbole/bdjuno/v4/modules/wormhole"
)

// updateGuardianSetCmd returns the Cobra command allowing to refresh x/wormhole
// guardian set list in database
func updateGuardianSetCmd(parseConfig *parsecmdtypes.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "guardian-set",
		Short: "Refresh guardian set",
		RunE: func(cmd *cobra.Command, args []string) error {
			parseCtx, err := parsecmdtypes.GetParserContext(config.Cfg, parseConfig)
			if err != nil {
				return err
			}

			sources, err := modulestypes.BuildSources(config.Cfg.Node, parseCtx.EncodingConfig)
			if err != nil {
				return err
			}

			// Get the database
			db := database.Cast(parseCtx.Database)

			// Build wormhole module
			wormholeModule := wormhole.NewModule(sources.WormholeSource, parseCtx.EncodingConfig.Marshaler, db)

			err = wormholeModule.UpdateGuardianSet()
			if err != nil {
				return fmt.Errorf("error while updating guardian set: %s", err)
			}

			return nil
		},
	}
}
