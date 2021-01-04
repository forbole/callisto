package staking

import (
	"encoding/json"
	"time"

	"github.com/desmos-labs/juno/client"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/rs/zerolog/log"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/x/staking/types"
)

func HandleGenesis(
	genDoc *tmtypes.GenesisDoc, appState map[string]json.RawMessage,
	cdc *codec.Codec, cp *client.Proxy, db *database.BigDipperDb,
) error {
	log.Debug().Str("module", "staking").Msg("parsing genesis")

	// Read the genesis state
	var genState staking.GenesisState
	err := cdc.UnmarshalJSON(appState[staking.ModuleName], &genState)
	if err != nil {
		return err
	}

	// Save the params
	err = saveParams(genState.Params, db)
	if err != nil {
		return err
	}

	// Save the validators
	err = saveValidators(genState.Validators, db)
	if err != nil {
		return err
	}

	err = saveValidatorsCommissions(genState.Validators, genDoc, db)
	if err != nil {
		return err
	}

	// Save the delegations
	err = saveDelegations(genState, genDoc, db)
	if err != nil {
		return err
	}

	// Save the unbonding delegations
	err = saveUnbondingDelegations(genState, genDoc, db)
	if err != nil {
		return err
	}

	// Save the re-delegations
	err = saveRedelegations(genDoc.GenesisTime, genState, db)
	if err != nil {
		return err
	}

	// Save the description
	err = saveDescription(genDoc.GenesisTime, genState.Validators, db)
	if err != nil {
		return err
	}

	return nil
}

// saveParams saves the given params into the database
func saveParams(params staking.Params, db *database.BigDipperDb) error {
	return db.SaveStakingParams(types.NewStakingParams(
		params.BondDenom,
	))
}

// saveValidators stores the validators data present inside the given genesis state
func saveValidators(validators staking.Validators, db *database.BigDipperDb) error {
	bdValidators := make([]types.Validator, len(validators))
	for i, validator := range validators {
		bdValidators[i] = types.NewValidator(
			validator.ConsAddress(),
			validator.OperatorAddress,
			validator.GetConsPubKey(),
			sdk.AccAddress(validator.ConsAddress()),
			&validator.Commission.MaxChangeRate,
			&validator.Commission.MaxRate,
		)
	}

	return db.SaveValidators(bdValidators)
}

//saveValidatorsCommissions save initial commission for validators
func saveValidatorsCommissions(
	validators staking.Validators, genesisDoc *tmtypes.GenesisDoc, db *database.BigDipperDb,
) error {
	for _, account := range validators {
		err := db.SaveValidatorCommission(types.NewValidatorCommission(
			account.OperatorAddress,
			&account.Commission.Rate,
			&account.MinSelfDelegation,
			1, genesisDoc.GenesisTime,
		))
		if err != nil {
			return err
		}
	}

	return nil
}

// saveDelegations stores the delegations data present inside the given genesis state
func saveDelegations(genState staking.GenesisState, genesisDoc *tmtypes.GenesisDoc, db *database.BigDipperDb) error {
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
				genesisDoc.GenesisTime,
			))
		}
	}

	if err := db.SaveCurrentDelegations(delegations); err != nil {
		return err
	}
	return nil
}

// saveUnbondingDelegations stores the unbonding delegations data present inside the given genesis state
func saveUnbondingDelegations(genState staking.GenesisState, genesisDoc *tmtypes.GenesisDoc, db *database.BigDipperDb) error {
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
					genesisDoc.GenesisTime,
				))
			}
		}
	}

	return db.SaveCurrentUnbondingDelegations(unbondingDelegations)
}

// saveRedelegations stores the redelegations data present inside the given genesis state
func saveRedelegations(genTime time.Time, genState staking.GenesisState, db *database.BigDipperDb) error {
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
				genTime,
			))
		}
	}

	return db.SaveCurrentRedelegations(redelegations)
}

// getDelegations returns the list of all the delegations that are
// related to the validator having the given validator address
func getDelegations(genData staking.Delegations, valAddr sdk.ValAddress) staking.Delegations {
	var delegations staking.Delegations
	for _, delegation := range genData {
		if delegation.ValidatorAddress.Equals(valAddr) {
			delegations = append(delegations, delegation)
		}
	}
	return delegations
}

// getUnbondingDelegations returns the list of all the unbonding delegations
// that are related to the validator having the given validator address
func getUnbondingDelegations(genData staking.UnbondingDelegations, valAddr sdk.ValAddress) staking.UnbondingDelegations {
	var unbondingDelegations staking.UnbondingDelegations
	for _, unbondingDelegation := range genData {
		if unbondingDelegation.ValidatorAddress.Equals(valAddr) {
			unbondingDelegations = append(unbondingDelegations, unbondingDelegation)
		}
	}
	return unbondingDelegations
}

//saveValidatorsCommissions save initial commission for validators
func saveDescription(genTime time.Time, validators staking.Validators, db *database.BigDipperDb) error {
	for _, account := range validators {
		err := db.SaveValidatorDescription(types.NewValidatorDescription(
			account.OperatorAddress,
			account.Description,
			1,
			genTime,
		))
		if err != nil {
			return err
		}
	}

	return nil
}
