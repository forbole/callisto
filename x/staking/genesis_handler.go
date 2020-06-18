package staking

import (
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/desmos-labs/juno/parse/worker"
	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/x/staking/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

func GenesisHandler(
	codec *codec.Codec, genesisDoc *tmtypes.GenesisDoc, appState map[string]json.RawMessage, w worker.Worker,
) error {
	bigDipperDb, ok := w.Db.(database.BigDipperDb)
	if !ok {
		return fmt.Errorf("provided database is not a BigDipper database")
	}

	// Read the genesis state
	var genState staking.GenesisState
	if err := codec.UnmarshalJSON(appState[staking.ModuleName], &genState); err != nil {
		return err
	}

	// Get the validators
	validators := make([]types.Validator, len(genState.Validators))
	for _, validator := range genState.Validators {
		validators = append(validators, types.NewValidator(
			validator.ConsAddress(),
			validator.OperatorAddress,
			validator.GetConsPubKey(),
		))
	}

	// Get the delegations
	delegations := make([]types.ValidatorDelegations, len(genState.Validators))
	for _, validator := range genState.Validators {
		delegations = append(delegations, types.ValidatorDelegations{
			ConsAddress:          validator.ConsAddress(),
			Delegations:          getDelegations(genState.Delegations, validator.OperatorAddress),
			UnbondingDelegations: getUnbondingDelegations(genState.UnbondingDelegations, validator.OperatorAddress),
			Height:               0,
			Timestamp:            genesisDoc.GenesisTime,
		})
	}

	// Save the validators
	if err := bigDipperDb.SaveValidators(validators); err != nil {
		return err
	}

	// Save the delegations
	if err := bigDipperDb.SaveValidatorsDelegations(delegations, 0, genesisDoc.GenesisTime); err != nil {
		return err
	}

	return nil
}

// getDelegations returns the list of all the delegations that are
// related to the validator having the given validator address
func getDelegations(genDelegations staking.Delegations, valAddr sdk.ValAddress) staking.Delegations {
	var delegations staking.Delegations
	for _, delegation := range genDelegations {
		if delegation.ValidatorAddress.Equals(valAddr) {
			delegations = append(delegations, delegation)
		}
	}
	return delegations
}

// getUnbondingDelegations returns the list of all the unbonding delegations
// that are related to the validator having the given validator address
func getUnbondingDelegations(
	genUnbondingDelegations staking.UnbondingDelegations, valAddr sdk.ValAddress,
) staking.UnbondingDelegations {
	var unbondingDelegations staking.UnbondingDelegations
	for _, unbondingDelegation := range genUnbondingDelegations {
		if unbondingDelegation.ValidatorAddress.Equals(valAddr) {
			unbondingDelegations = append(unbondingDelegations, unbondingDelegation)
		}
	}
	return unbondingDelegations
}
