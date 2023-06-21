package staking

import (
	"fmt"

	parsecmdtypes "github.com/forbole/juno/v5/cmd/parse/types"
	"github.com/forbole/juno/v5/types/config"
	"github.com/spf13/cobra"

	"github.com/forbole/bdjuno/v4/database"
	"github.com/forbole/bdjuno/v4/modules/staking"
	modulestypes "github.com/forbole/bdjuno/v4/modules/types"
)

// poolCmd returns the Cobra command allowing to refresh x/staking pool
func poolCmd(parseConfig *parsecmdtypes.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "pool",
		Short: "Refresh staking pool",
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

			// Build staking module
			stakingModule := staking.NewModule(sources.StakingSource, parseCtx.EncodingConfig.Codec, db)

			err = stakingModule.UpdateStakingPool()
			if err != nil {
				return fmt.Errorf("error while updating staking pool: %s", err)
			}

			return nil
		},
	}
}
