package staking

import (
	"fmt"

	"github.com/forbole/juno/v2/cmd/parse"
	"github.com/forbole/juno/v2/types/config"
	"github.com/spf13/cobra"

	"github.com/forbole/bdjuno/v2/database"
	"github.com/forbole/bdjuno/v2/modules"
	"github.com/forbole/bdjuno/v2/modules/staking"
)

// slashesCmd returns a Cobra command that allows to fix the delegations for all the slashed validators.
func slashesCmd(parseConfig *parse.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "slashes",
		Short: "Fix the delegations for all the slashed validators, taking their delegations from the latest known height",
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
			height, err := parseCtx.Node.LatestHeight()
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
				err = db.ReplaceValidatorDelegations(height, validator.OperatorAddress, staking.ConvertDelegationsResponses(delegations))
				if err != nil {
					return fmt.Errorf("error while refreshing validator delegations: %s", err)
				}
			}

			return nil
		},
	}
}
