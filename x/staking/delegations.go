package staking

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/desmos-labs/juno/parse/client"
	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/x/staking/types"
	"github.com/rs/zerolog/log"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
)

// updateDelegations updates all the validators delegations
func updateDelegations(cp client.ClientProxy, db database.BigDipperDb) error {
	log.Debug().
		Str("module", "staking").
		Str("operation", "delegations").
		Msg("getting delegations")

	var block tmctypes.ResultBlock
	err := cp.QueryLCD("/blocks/latest", &block)
	if err != nil {
		return err
	}

	validators, err := db.GetValidators()
	if err != nil {
		return err
	}

	validatorDelegations := make([]types.ValidatorDelegations, len(validators))
	for index, validator := range validators {

		// Get the delegations
		dEndpoint := fmt.Sprintf("/staking/validators/%s/delegations?height=%d",
			validator.GetOperator().String(), block.Block.Height)
		var delegations staking.DelegationResponses
		if _, err := cp.QueryLCDWithHeight(dEndpoint, &delegations); err != nil {
			return err
		}

		// Get the unbonding delegations
		udEndpoint := fmt.Sprintf("/staking/validators/%s/unbonding_delegations?height=%d",
			validator.GetOperator().String(), block.Block.Height)
		var unbondingDelegations staking.UnbondingDelegations
		if _, err := cp.QueryLCDWithHeight(udEndpoint, &unbondingDelegations); err != nil {
			return err
		}

		validatorDelegations[index] = types.ValidatorDelegations{
			ConsAddress:          validator.GetConsAddr(),
			Delegations:          delegations,
			UnbondingDelegations: unbondingDelegations,
			Height:               block.Block.Height,
			Timestamp:            block.Block.Time,
		}
	}

	log.Debug().
		Str("module", "staking").
		Str("operation", "delegations").
		Msg("saving delegations")
	return db.SaveValidatorsDelegations(validatorDelegations, block.Block.Height, block.Block.Time)
}
