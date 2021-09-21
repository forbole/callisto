package staking

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx"

	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/modules/staking/utils"
	"github.com/forbole/bdjuno/types"

	tmtypes "github.com/tendermint/tendermint/types"

	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/rs/zerolog/log"
)

func HandleGenesis(
	doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage, cdc codec.Marshaler, db *database.Db,
) error {
	log.Debug().Str("module", "staking").Msg("parsing genesis")

	// Read the genesis state
	var genState stakingtypes.GenesisState
	err := cdc.UnmarshalJSON(appState[stakingtypes.ModuleName], &genState)
	if err != nil {
		return fmt.Errorf("error while unmarshaling staking state: %s", err)
	}

	// Save the params
	err = saveParams(doc.InitialHeight, genState.Params, db)
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

	// Save the description
	err = saveValidatorDescription(doc, genState.Validators, db)
	if err != nil {
		return fmt.Errorf("error while storing staking genesis validator descriptions: %s", err)
	}

	err = saveValidatorsCommissions(doc.InitialHeight, genState.Validators, db)
	if err != nil {
		return fmt.Errorf("error while storing staking genesis validators commissions: %s", err)
	}

	return nil
}

func parseGenesisTransactions(
	doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage, cdc codec.Marshaler, db *database.Db,
) error {
	var genUtilState genutiltypes.GenesisState
	err := cdc.UnmarshalJSON(appState[genutiltypes.ModuleName], &genUtilState)
	if err != nil {
		return fmt.Errorf("error while unmarhsaling genutil state: %s", err)
	}

	for _, genTxBz := range genUtilState.GetGenTxs() {
		// Unmarshal the transaction
		var genTx tx.Tx
		err = cdc.UnmarshalJSON(genTxBz, &genTx)
		if err != nil {
			return fmt.Errorf("error while unmashasling genesis tx: %s", err)
		}

		for _, msg := range genTx.GetMsgs() {
			// Handle the message properly
			createValMsg, ok := msg.(*stakingtypes.MsgCreateValidator)
			if !ok {
				continue
			}

			err = utils.StoreValidatorFromMsgCreateValidator(doc.InitialHeight, createValMsg, cdc, db)
			if err != nil {
				return fmt.Errorf("error while storing validators from MsgCreateValidator: %s", err)
			}
		}
	}

	return nil
}

// -------------------------------------------------------------------------------------------------------------------

// saveParams saves the given params into the database
func saveParams(height int64, params stakingtypes.Params, db *database.Db) error {
	return db.SaveStakingParams(types.NewStakingParams(params, height))
}

// --------------------------------------------------------------------------------------------------------------------

// saveValidators stores the validators data present inside the given genesis state
func saveValidators(
	doc *tmtypes.GenesisDoc, validators stakingtypes.Validators, cdc codec.Marshaler, db *database.Db,
) error {
	vals := make([]types.Validator, len(validators))
	for i, val := range validators {
		validator, err := utils.ConvertValidator(cdc, val, doc.InitialHeight)
		if err != nil {
			return err
		}

		vals[i] = validator
	}

	return db.SaveValidatorsData(vals)
}

// saveValidatorDescription saves the description for the given validators
func saveValidatorDescription(doc *tmtypes.GenesisDoc, validators stakingtypes.Validators, db *database.Db) error {
	for _, account := range validators {
		description, err := utils.ConvertValidatorDescription(
			account.OperatorAddress,
			account.Description,
			doc.InitialHeight,
		)
		if err != nil {
			return fmt.Errorf("error while converting validator description: %s", err)
		}

		err = db.SaveValidatorDescription(description)
		if err != nil {
			return err
		}
	}

	return nil
}

// --------------------------------------------------------------------------------------------------------------------

// saveDelegations stores the delegations data present inside the given genesis state
func saveDelegations(doc *tmtypes.GenesisDoc, genState stakingtypes.GenesisState, db *database.Db) error {
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

	return db.SaveDelegations(delegations)
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

// --------------------------------------------------------------------------------------------------------------------

// saveUnbondingDelegations stores the unbonding delegations data present inside the given genesis state
func saveUnbondingDelegations(doc *tmtypes.GenesisDoc, genState stakingtypes.GenesisState, db *database.Db) error {
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

// --------------------------------------------------------------------------------------------------------------------

// saveRedelegations stores the redelegations data present inside the given genesis state
func saveRedelegations(doc *tmtypes.GenesisDoc, genState stakingtypes.GenesisState, db *database.Db) error {
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

// --------------------------------------------------------------------------------------------------------------------

// saveValidatorsCommissions save the initial commission for each validator
func saveValidatorsCommissions(height int64, validators stakingtypes.Validators, db *database.Db) error {
	for _, account := range validators {
		err := db.SaveValidatorCommission(types.NewValidatorCommission(
			account.OperatorAddress,
			&account.Commission.Rate,
			&account.MinSelfDelegation,
			height,
		))
		if err != nil {
			return err
		}
	}

	return nil
}
