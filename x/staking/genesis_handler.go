package staking

import (
	"encoding/json"
	"fmt"

	bstakingutils "github.com/forbole/bdjuno/x/staking/utils"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/rs/zerolog/log"

	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/x/staking/types"
)

func HandleGenesis(
	appState map[string]json.RawMessage, cdc codec.Marshaler, db *database.BigDipperDb,
) error {
	log.Debug().Str("module", "staking").Msg("parsing genesis")

	// Read the genesis state
	var genState stakingtypes.GenesisState
	err := cdc.UnmarshalJSON(appState[stakingtypes.ModuleName], &genState)
	if err != nil {
		return err
	}

	// Save the params
	err = saveParams(genState.Params, db)
	if err != nil {
		return fmt.Errorf("error while storing staking genesis params: %s", err)
	}

	// Save the validators
	err = saveValidators(genState.Validators, cdc, db)
	if err != nil {
		return fmt.Errorf("error while storing staking genesis validators: %s", err)
	}

	// Save the description
	err = saveValidatorDescription(genState.Validators, db)
	if err != nil {
		return fmt.Errorf("error while storing staking genesis validator descriptions: %s", err)
	}

	err = saveValidatorsCommissions(genState.Validators, db)
	if err != nil {
		return fmt.Errorf("error while storing staking genesis validators commissions: %s", err)
	}

	// Save the delegations
	err = saveDelegations(genState, db)
	if err != nil {
		return fmt.Errorf("error while storing staking genesis delegations: %s", err)
	}

	// Save the unbonding delegations
	err = saveUnbondingDelegations(genState, db)
	if err != nil {
		return fmt.Errorf("error while storing staking genesis unbonding delegations: %s", err)
	}

	// Save the re-delegations
	err = saveRedelegations(genState, db)
	if err != nil {
		return fmt.Errorf("error while storing staking genesis redelegations: %s", err)
	}

	return nil
}

// saveParams saves the given params into the database
func saveParams(params stakingtypes.Params, db *database.BigDipperDb) error {
	return db.SaveStakingParams(types.NewStakingParams(
		params.BondDenom,
	))
}

// saveValidators stores the validators data present inside the given genesis state
func saveValidators(validators stakingtypes.Validators, cdc codec.Marshaler, db *database.BigDipperDb) error {
	bdValidators := make([]types.Validator, len(validators))
	for i, validator := range validators {
		consAddr, err := bstakingutils.GetValidatorConsAddr(cdc, validator)
		if err != nil {
			return err
		}

		consPubKey, err := bstakingutils.GetValidatorConsPubKey(cdc, validator)
		if err != nil {
			return err
		}

		bdValidators[i] = types.NewValidator(
			consAddr.String(),
			validator.OperatorAddress,
			consPubKey.String(),
			sdk.AccAddress(consAddr).String(),
			&validator.Commission.MaxChangeRate,
			&validator.Commission.MaxRate,
		)
	}

	return db.SaveValidators(bdValidators)
}

// saveValidatorsCommissions save the initial commission for each validator
func saveValidatorsCommissions(validators stakingtypes.Validators, db *database.BigDipperDb) error {
	for _, account := range validators {
		err := db.SaveValidatorCommission(types.NewValidatorCommission(
			account.OperatorAddress,
			&account.Commission.Rate,
			&account.MinSelfDelegation,
			1,
		))
		if err != nil {
			return err
		}
	}

	return nil
}

// saveDelegations stores the delegations data present inside the given genesis state
func saveDelegations(genState stakingtypes.GenesisState, db *database.BigDipperDb) error {
	var delegations []types.Delegation
	for _, validator := range genState.Validators {
		tokens := validator.Tokens
		delegatorShares := validator.DelegatorShares

		for _, delegation := range getDelegations(genState.Delegations, validator.OperatorAddress) {
			delegationAmount := tokens.ToDec().Mul(delegation.Shares).Quo(delegatorShares).TruncateInt()
			delegations = append(delegations, types.NewDelegation(
				delegation.DelegatorAddress,
				validator.OperatorAddress,
				sdk.NewCoin(genState.Params.BondDenom, delegationAmount),
				delegation.Shares.String(),
				1,
			))
		}
	}

	if err := db.SaveDelegations(delegations); err != nil {
		return err
	}
	return nil
}

// saveUnbondingDelegations stores the unbonding delegations data present inside the given genesis state
func saveUnbondingDelegations(genState stakingtypes.GenesisState, db *database.BigDipperDb) error {
	var unbondingDelegations []types.UnbondingDelegation
	for _, validator := range genState.Validators {
		valUD := getUnbondingDelegations(genState.UnbondingDelegations, validator.OperatorAddress)
		for _, ud := range valUD {
			for _, entry := range ud.Entries {
				unbondingDelegations = append(unbondingDelegations, types.NewUnbondingDelegation(
					ud.DelegatorAddress,
					validator.OperatorAddress,
					sdk.NewCoin(genState.Params.BondDenom, entry.InitialBalance),
					entry.CompletionTime,
					entry.CreationHeight,
				))
			}
		}
	}

	return db.SaveUnbondingDelegations(unbondingDelegations)
}

// saveRedelegations stores the redelegations data present inside the given genesis state
func saveRedelegations(genState stakingtypes.GenesisState, db *database.BigDipperDb) error {
	var redelegations []types.Redelegation
	for _, redelegation := range genState.Redelegations {
		for _, entry := range redelegation.Entries {
			redelegations = append(redelegations, types.NewRedelegation(
				redelegation.DelegatorAddress,
				redelegation.ValidatorSrcAddress,
				redelegation.ValidatorDstAddress,
				sdk.NewCoin(genState.Params.BondDenom, entry.InitialBalance),
				entry.CompletionTime,
				entry.CreationHeight,
			))
		}
	}

	return db.SaveRedelegations(redelegations)
}

// getDelegations returns the list of all the delegations that are
// related to the validator having the given validator address
func getDelegations(genData stakingtypes.Delegations, valAddr string) stakingtypes.Delegations {
	var delegations stakingtypes.Delegations
	for _, delegation := range genData {
		if delegation.ValidatorAddress == valAddr {
			delegations = append(delegations, delegation)
		}
	}
	return delegations
}

// getUnbondingDelegations returns the list of all the unbonding delegations
// that are related to the validator having the given validator address
func getUnbondingDelegations(genData stakingtypes.UnbondingDelegations, valAddr string) stakingtypes.UnbondingDelegations {
	var unbondingDelegations stakingtypes.UnbondingDelegations
	for _, unbondingDelegation := range genData {
		if unbondingDelegation.ValidatorAddress == valAddr {
			unbondingDelegations = append(unbondingDelegations, unbondingDelegation)
		}
	}
	return unbondingDelegations
}

// saveValidatorDescription saves the description for the given validators
func saveValidatorDescription(validators stakingtypes.Validators, db *database.BigDipperDb) error {
	for _, account := range validators {
		err := db.SaveValidatorDescription(types.NewValidatorDescription(
			account.OperatorAddress,
			account.Description,
			1,
		))
		if err != nil {
			return err
		}
	}

	return nil
}
