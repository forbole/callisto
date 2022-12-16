package staking

import (
	"fmt"

	modulestypes "github.com/forbole/bdjuno/v3/modules/types"

	parsecmdtypes "github.com/forbole/juno/v4/cmd/parse/types"
	"github.com/forbole/juno/v4/types/config"
	"github.com/spf13/cobra"

	"github.com/forbole/bdjuno/v3/database"
	"github.com/forbole/bdjuno/v3/modules/staking"
)

// validatorsCmd returns a Cobra command that allows to fix the validator infos for all validators.
func validatorsCmd(parseConfig *parsecmdtypes.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "validators",
		Short: "Fix the information about validators taking them from the latest known height",
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

			// Build the staking module
			stakingModule := staking.NewModule(sources.StakingSource, nil, parseCtx.EncodingConfig.Marshaler, db)

			// Get latest height
			height, err := parseCtx.Node.LatestHeight()
			if err != nil {
				return fmt.Errorf("error while getting latest block height: %s", err)
			}

			// Get all validators
			validators, err := sources.StakingSource.GetValidatorsWithStatus(height, "")
			if err != nil {
				return fmt.Errorf("error while getting validators: %s", err)
			}

			// Refresh each validator
			for _, validator := range validators {
				err = stakingModule.RefreshValidatorInfos(height, validator.OperatorAddress)
				if err != nil {
					return fmt.Errorf("error while refreshing validator: %s", err)
				}
			}

			return nil
		},
	}
}
