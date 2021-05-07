package staking

import (
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/forbole/bdjuno/types"
)

func HandleGenesis(
	doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage, cdc codec.Marshaler, db DB,
) error {
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

	// Parse genesis transactions
	err = parseGenesisTransactions(doc, appState, cdc, db)
	if err != nil {
		return fmt.Errorf("error while storing genesis transactions: %s", err)
	}

	// Save the validators
	err = saveValidators(doc, genState.Validators, cdc, db)
	if err != nil {
		return fmt.Errorf("error while storing staking genesis validators: %s", err)
	}

	// Save the delegations
	err = saveDelegations(doc, genState, db)
	if err != nil {
		return fmt.Errorf("error while storing staking genesis delegations: %s", err)
	}

	// Save the unbonding delegations
	err = saveUnbondingDelegations(doc, genState, db)
	if err != nil {
		return fmt.Errorf("error while storing staking genesis unbonding delegations: %s", err)
	}

	// Save the re-delegations
	err = saveRedelegations(doc, genState, db)
	if err != nil {
		return fmt.Errorf("error while storing staking genesis redelegations: %s", err)
	}

	return nil
}

// -------------------------------------------------------------------------------------------------------------------

func parseGenesisTransactions(
	doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage, cdc codec.Marshaler, db DB,
) error {
	var genUtilState genutiltypes.GenesisState
	err := cdc.UnmarshalJSON(appState[genutiltypes.ModuleName], &genUtilState)
	if err != nil {
		return err
	}

	for _, genTxBz := range genUtilState.GetGenTxs() {
		// Unmarshal the transaction
		var genTx tx.Tx
		if err := cdc.UnmarshalJSON(genTxBz, &genTx); err != nil {
			return err
		}

		for _, msg := range genTx.GetMsgs() {
			// Handle the message properly
			createValMsg, ok := msg.(*stakingtypes.MsgCreateValidator)
			if !ok {
				continue
			}

			err = StoreValidatorFromMsgCreateValidator(doc.InitialHeight, createValMsg, cdc, db)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// -------------------------------------------------------------------------------------------------------------------

// saveParams saves the given params into the database
func saveParams(params stakingtypes.Params, db DB) error {
	return db.SaveStakingParams(types.NewStakingParams(
		params.BondDenom,
	))
}

// saveValidators stores the validators data present inside the given genesis state
func saveValidators(
	doc *tmtypes.GenesisDoc, validators stakingtypes.Validators, cdc codec.Marshaler, db DB,
) error {
	vals := make([]types.Validator, len(validators))
	for i, val := range validators {
		validator, err := ConvertValidator(cdc, val, doc.InitialHeight)
		if err != nil {
			return err
		}

		vals[i] = validator
	}

	return db.SaveValidators(vals)
}

// saveDelegations stores the delegations data present inside the given genesis state
func saveDelegations(doc *tmtypes.GenesisDoc, genState stakingtypes.GenesisState, db DB) error {
	var delegations []types.Delegation
	for _, validator := range genState.Validators {
		tokens := validator.Tokens
		delegatorShares := validator.DelegatorShares

		for _, delegation := range findDelegations(genState.Delegations, validator.OperatorAddress) {
			delegationAmount := tokens.ToDec().Mul(delegation.Shares).Quo(delegatorShares).TruncateInt()
			delegations = append(delegations, types.NewDelegation(
				delegation.DelegatorAddress,
				validator.OperatorAddress,
				sdk.NewCoin(genState.Params.BondDenom, delegationAmount),
				doc.InitialHeight,
			))
		}
	}

	if err := db.SaveDelegations(delegations); err != nil {
		return err
	}
	return nil
}

// saveUnbondingDelegations stores the unbonding delegations data present inside the given genesis state
func saveUnbondingDelegations(doc *tmtypes.GenesisDoc, genState stakingtypes.GenesisState, db DB) error {
	var unbondingDelegations []types.UnbondingDelegation
	for _, validator := range genState.Validators {
		valUD := findUnbondingDelegations(genState.UnbondingDelegations, validator.OperatorAddress)
		for _, ud := range valUD {
			for _, entry := range ud.Entries {
				unbondingDelegations = append(unbondingDelegations, types.NewUnbondingDelegation(
					ud.DelegatorAddress,
					validator.OperatorAddress,
					sdk.NewCoin(genState.Params.BondDenom, entry.InitialBalance),
					entry.CompletionTime,
					doc.InitialHeight,
				))
			}
		}
	}

	return db.SaveUnbondingDelegations(unbondingDelegations)
}

// saveRedelegations stores the redelegations data present inside the given genesis state
func saveRedelegations(doc *tmtypes.GenesisDoc, genState stakingtypes.GenesisState, db DB) error {
	var redelegations []types.Redelegation
	for _, redelegation := range genState.Redelegations {
		for _, entry := range redelegation.Entries {
			redelegations = append(redelegations, types.NewRedelegation(
				redelegation.DelegatorAddress,
				redelegation.ValidatorSrcAddress,
				redelegation.ValidatorDstAddress,
				sdk.NewCoin(genState.Params.BondDenom, entry.InitialBalance),
				entry.CompletionTime,
				doc.InitialHeight,
			))
		}
	}

	return db.SaveRedelegations(redelegations)
}

// findDelegations returns the list of all the delegations that are
// related to the validator having the given validator address
func findDelegations(genData stakingtypes.Delegations, valAddr string) stakingtypes.Delegations {
	var delegations stakingtypes.Delegations
	for _, delegation := range genData {
		if delegation.ValidatorAddress == valAddr {
			delegations = append(delegations, delegation)
		}
	}
	return delegations
}

// findUnbondingDelegations returns the list of all the unbonding delegations
// that are related to the validator having the given validator address
func findUnbondingDelegations(genData stakingtypes.UnbondingDelegations, valAddr string) stakingtypes.UnbondingDelegations {
	var unbondingDelegations stakingtypes.UnbondingDelegations
	for _, unbondingDelegation := range genData {
		if unbondingDelegation.ValidatorAddress == valAddr {
			unbondingDelegations = append(unbondingDelegations, unbondingDelegation)
		}
	}
	return unbondingDelegations
}
