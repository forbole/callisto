package staking

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/desmos-labs/juno/parse/worker"
	"github.com/forbole/bdjuno/database"
	dbtypes "github.com/forbole/bdjuno/database/types"
	"github.com/rs/zerolog/log"
	"github.com/tendermint/tendermint/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

func GenesisHandler(codec *codec.Codec, genesisDoc *types.GenesisDoc, appState map[string]json.RawMessage, w worker.Worker) error {
	log.Debug().Str("module", "auth").Msg("parsing genesis")

	bigDipperDb, ok := w.Db.(database.BigDipperDb)
	if !ok {
		return fmt.Errorf("given database instance is not a BigDipperDb")
	}
	// Read the genesis state
	var genState staking.GenesisState
	if err := codec.UnmarshalJSON(appState[staking.ModuleName], &genState); err != nil {
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

	var stakingGenesisState staking.GenesisState
	if err := codec.UnmarshalJSON(appState[staking.ModuleName], &stakingGenesisState); err != nil {
		return err
	}

	if err := codec.UnmarshalJSON(appState[staking.ModuleName], &stakingGenesisState); err != nil {
		return err
	}
	err := InitialCommission(stakingGenesisState, bigDipperDb)
	if err != nil {
		return err
	}
	return nil
}

func InitialCommission(stakingGenesisState staking.GenesisState, db database.BigDipperDb) error {
	// Store the accounts
	accounts := make([]dbtypes.ValidatorCommission, len(stakingGenesisState.Validators))
	for index, account := range stakingGenesisState.Validators {
		accounts[index] = dbtypes.NewValidatorCommission(account.ConsAddress.String(), time.Now,
			account.Commission, account.MinSelfDelegation, 0)
	}

	err := db.SaveValidatorCommissions(accounts)
	if err != nil {
		return err
	}
	return nil
}

func InitialInformation(stakingGenesisState staking.GenesisState, db database.BigDipperDb) error {
	accounts := make([]dbtypes.ValidatorInfoRow, len(stakingGenesisState.Validators))
	for index, account := range stakingGenesisState.Validators {
		accounts[index] = dbtypes.NewValidatorInfoRow(
			account.ConsAddress().String(),
			account.OperatorAddress.String(),
			account.Description.Moniker,
			account.Description.Identity,
			account.Description.Website,
			account.Description.SecurityContact,
			account.Description.Details,
		)
	}

	err := db.SaveInitialValidatorInfo(accounts)
	if err != nil {
		return err
	}
	return nil
}

// saveDelegations stores the delegations data present inside the given genesis state
func saveDelegations(genState staking.GenesisState, genesisDoc *tmtypes.GenesisDoc, db database.BigDipperDb) error {
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
				0,
				genesisDoc.GenesisTime,
			))
		}
	}

	return db.SaveDelegations(delegations)
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
