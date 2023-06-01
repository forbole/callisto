package stakers

import (
	"fmt"

	parsecmdtypes "github.com/forbole/juno/v4/cmd/parse/types"
	"github.com/forbole/juno/v4/types/config"
	"github.com/spf13/cobra"

	"github.com/forbole/bdjuno/v4/database"
	"github.com/forbole/bdjuno/v4/modules/stakers"
	modulestypes "github.com/forbole/bdjuno/v4/modules/types"
)

// stakersCmd returns the Cobra command allowing to refresh x/stakers values
func stakersCmd(parseConfig *parsecmdtypes.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "protocol-validators",
		Short: "Refresh all stakers information in database",
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

			// Build stakers module
			stakersModule := stakers.NewModule(sources.StakersSource, parseCtx.EncodingConfig.Codec, db)

			height, err := parseCtx.Node.LatestHeight()
			if err != nil {
				return fmt.Errorf("error while getting chain latest block height: %s", err)
			}

			err = stakersModule.UpdateProtocolValidatorsInfo(height)
			if err != nil {
				return fmt.Errorf("error while updating stakers: %s", err)
			}

			err = stakersModule.UpdateProtocolValidatorsCommission(height)
			if err != nil {
				return fmt.Errorf("error while updating stakers commission: %s", err)
			}

			err = stakersModule.UpdateProtocolValidatorsDescription(height)
			if err != nil {
				return fmt.Errorf("error while updating stakers description: %s", err)
			}

			return nil
		},
	}
}
