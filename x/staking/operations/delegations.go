package operations

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/desmos-labs/juno/parse/client"
	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/x/staking/types"
	"github.com/rs/zerolog/log"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
)

// UpdateDelegations updates all the validators delegations
func UpdateDelegations(cp client.ClientProxy, db database.BigDipperDb) error {
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
		dEndpoint := fmt.Sprintf("/staking/validators/%s/delegationResponses?height=%d",
			validator.GetOperator().String(), block.Block.Height)
		var delegationResponses staking.DelegationResponses
		if _, err := cp.QueryLCDWithHeight(dEndpoint, &delegationResponses); err != nil {
			return err
		}

		delegations := make([]staking.Delegation, len(delegationResponses))
		for index, delegation := range delegationResponses {
			delegations[index] = delegation.Delegation
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
