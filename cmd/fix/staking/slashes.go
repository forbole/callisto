package staking

import (
	"fmt"

	"github.com/desmos-labs/juno/v2/cmd/parse"
	"github.com/desmos-labs/juno/v2/types/config"
	"github.com/spf13/cobra"

	"github.com/forbole/bdjuno/v2/database"
	"github.com/forbole/bdjuno/v2/modules"
	"github.com/forbole/bdjuno/v2/modules/staking"
)

// StakingCmd returns the Cobra command that allows to fix all the things related to the x/staking module
func StakingCmd(parseConfig *parse.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "staking",
		Short: "Fix things related to the x/staking module",
	}

	cmd.AddCommand(
		validatorsCmd(parseConfig),
		slashesCmd(parseConfig),
	)

	return cmd
}

// validatorsCmd returns a Cobra command that allows to fix the validator infos for all validators.
func validatorsCmd(parseConfig *parse.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "validators",
		Short: "Fix the information about validators taking them from the latest height stored inside the database.",
		RunE: func(cmd *cobra.Command, args []string) error {
			parseCtx, err := parse.GetParsingContext(parseConfig)
			if err != nil {
				return err
			}

			sources, err := modules.BuildSources(config.Cfg.Node, parseCtx.EncodingConfig)
			if err != nil {
				return err
			}

			// Get the database
			db := database.Cast(parseCtx.Database)

			// Build the staking module
			stakingModule := staking.NewModule(sources.StakingSource, nil, nil, nil, parseCtx.EncodingConfig.Marshaler, db)

			// Get latest height
			height, err := db.GetLastBlockHeight()
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

// slashesCmd returns a Cobra command that allows to fix the delegations for all the slashed validators.
func slashesCmd(parseConfig *parse.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "slashes",
		Short: "Fix the delegations for all the slashed validators, taking their delegations from the latest height stored inside the database.",
		RunE: func(cmd *cobra.Command, args []string) error {
			parseCtx, err := parse.GetParsingContext(parseConfig)
			if err != nil {
				return err
			}

			sources, err := modules.BuildSources(config.Cfg.Node, parseCtx.EncodingConfig)
			if err != nil {
				return err
			}

			// Get the database
			db := database.Cast(parseCtx.Database)

			// Get latest height
			height, err := db.GetLastBlockHeight()
			if err != nil {
				return fmt.Errorf("error while getting latest block height: %s", err)
			}

			// Get all validators
			validators, err := sources.StakingSource.GetValidatorsWithStatus(height, "")
			if err != nil {
				return fmt.Errorf("error while getting validators: %s", err)
			}

			for _, validator := range validators {
				// Get the validator delegations
				delegations, err := sources.StakingSource.GetValidatorDelegations(height, validator.OperatorAddress)
				if err != nil {
					return fmt.Errorf("error while getting validator delegations: %s", err)
				}

				// Delete the old delegations
				err = db.DeleteValidatorDelegations(validator.OperatorAddress)
				if err != nil {
					return fmt.Errorf("error while deleting validator delegations: %s", err)
				}

				// Save the delegations
				err = db.SaveDelegations(staking.ConvertDelegationsResponses(height, delegations))
				if err != nil {
					return fmt.Errorf("error while saving delegations: %s", err)
				}
			}

			return nil
		},
	}
}
