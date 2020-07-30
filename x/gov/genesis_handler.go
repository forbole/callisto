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
	"github.com/rs/zerolog/log"
	tmtypes "github.com/tendermint/tendermint/types"
)

func GenesisHandler(codec *codec.Codec, genesisDoc *tmtypes.GenesisDoc, appState map[string]json.RawMessage, w worker.Worker) error {
	log.Debug().Str("module", "staking").Msg("parsing genesis")

	bigDipperDb, ok := w.Db.(database.BigDipperDb)
	if !ok {
		return fmt.Errorf("given database instance is not a BigDipperDb")
	}
	// Read the genesis state
	var genState staking.GenesisState
	if err := codec.UnmarshalJSON(appState[staking.ModuleName], &genState); err != nil {
		return err
	}

	// Save the validators
	if err := saveValidators(genState.Validators, bigDipperDb); err != nil {
		return err
	}

	if err := saveValidatorsCommissions(genState.Validators, genesisDoc, bigDipperDb); err != nil {
		return err
	}

	// Save the delegations
	if err := saveDelegations(genState, genesisDoc, bigDipperDb); err != nil {
		return err
	}

	// Save the unbonding delegations
	if err := saveUnbondingDelegations(genState, genesisDoc, bigDipperDb); err != nil {
		return err
	}

	// Save the re-delegations
	if err := saveRedelegations(genState, bigDipperDb); err != nil {
		return err
	}

	return nil
}

// saveValidators stores the validators data present inside the given genesis state
func saveValidators(validators staking.Validators, db database.BigDipperDb) error {
	bdValidators := make([]types.Validator, len(validators))
	for i, validator := range validators {
		bdValidators[i] = types.NewValidator(
			validator.ConsAddress(),
			validator.OperatorAddress,
			validator.GetConsPubKey(),
			validator.Description,
			sdk.AccAddress(validator.ConsAddress()),
		)
	}

	return db.SaveValidatorsData(bdValidators)
}

//saveValidatorsCommissions save initial commission for validators
func saveValidatorsCommissions(
	validators staking.Validators, genesisDoc *tmtypes.GenesisDoc, db database.BigDipperDb,
) error {
	accounts := make([]types.ValidatorCommission, len(validators))
	for index, account := range validators {
		accounts[index] = types.NewValidatorCommission(
			account.OperatorAddress,
			account.Commission.Rate.String(),
			account.MinSelfDelegation.Int64(),
			0, genesisDoc.GenesisTime,
		)
	}

	err := db.SaveValidatorCommissions(accounts)
	if err != nil {
		return err
	}
	return nil
}

// saveDelegations stores the delegations data present inside the given genesis state
func saveDelegations(genState staking.GenesisState, genesisDoc *tmtypes.GenesisDoc, db database.BigDipperDb) error {
	var delegations []types.Delegation
	var delegationShares []types.DelegationShare
	for _, validator := range genState.Validators {
		tokens := validator.Tokens
		delegatorShares := validator.DelegatorShares

		for _, delegation := range getDelegations(genState.Delegations, validator.OperatorAddress) {
			delegationAmount := tokens.ToDec().Mul(delegation.Shares).Quo(delegatorShares).TruncateInt()
			delegations = append(delegations, types.NewDelegation(
				delegation.DelegatorAddress,
				validator.OperatorAddress,
				sdk.NewCoin(genState.Params.BondDenom, delegationAmount),
				0,
				genesisDoc.GenesisTime,
			))
			delegationShares = append(delegationShares, types.NewDelegationShare(
				validator.OperatorAddress,
				delegation.DelegatorAddress,
				delegation.Shares.String(),
				0,
				genesisDoc.GenesisTime,
			))
		}
	}

	if err := db.SaveDelegationsShares(delegationShares); err != nil {
		return err
	}
	if err := db.SaveDelegations(delegations); err != nil {
		return err
	}
	return nil
}

// saveUnbondingDelegations stores the unbonding delegations data present inside the given genesis state
func saveUnbondingDelegations(genState staking.GenesisState, genesisDoc *tmtypes.GenesisDoc, db database.BigDipperDb) error {
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

	return db.SaveUnbondingDelegations(unbondingDelegations)
}

// saveRedelegations stores the redelegations data present inside the given genesis state
func saveRedelegations(genState staking.GenesisState, db database.BigDipperDb) error {
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
	if err := db.SaveRedelegations(redelegations); err != nil {
		return err
	}
	return nil
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

/*
// Store the accounts
	accounts := make([]exported.ValidatorI, len(stakingGenesisState.Validators))
	for index, account := range stakingGenesisState.Validators {
		accounts[index] = account.(exported.Account)
		selfAddress := sdk.AccAddress(account[index].Bytes())
		bstaking.NewSelfDelegation(validatorAddress,delegation.Shares.Int64(),
					float64(delegation.Shares.Int64())/float64(validator.DelegatorShares.Int64()*100,
					0,genDoc.GenesisTime)stakingGenesisState.Delegations
		//find the self delegation address has delegated to someone?
	}
*/
