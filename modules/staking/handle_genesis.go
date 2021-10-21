package staking

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx"

	"github.com/forbole/bdjuno/v2/types"

	tmtypes "github.com/tendermint/tendermint/types"

	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/rs/zerolog/log"
)

// HandleGenesis implements GenesisModule
func (m *Module) HandleGenesis(doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	log.Debug().Str("module", "staking").Msg("parsing genesis")

	// Read the genesis state
	var genState stakingtypes.GenesisState
	err := m.cdc.UnmarshalJSON(appState[stakingtypes.ModuleName], &genState)
	if err != nil {
		return fmt.Errorf("error while unmarshaling staking state: %s", err)
	}

	// Save the params
	err = m.saveParams(doc.InitialHeight, genState.Params)
	if err != nil {
		return fmt.Errorf("error while storing staking genesis params: %s", err)
	}

	// Parse genesis transactions
	err = m.parseGenesisTransactions(doc, appState)
	if err != nil {
		return fmt.Errorf("error while storing genesis transactions: %s", err)
	}

	// Save the validators
	err = m.saveValidators(doc, genState.Validators)
	if err != nil {
		return fmt.Errorf("error while storing staking genesis validators: %s", err)
	}

	// Save the delegations
	err = m.saveDelegations(doc, genState)
	if err != nil {
		return fmt.Errorf("error while storing staking genesis delegations: %s", err)
	}

	// Save the unbonding delegations
	err = m.saveUnbondingDelegations(doc, genState)
	if err != nil {
		return fmt.Errorf("error while storing staking genesis unbonding delegations: %s", err)
	}

	// Save the re-delegations
	err = m.saveRedelegations(doc, genState)
	if err != nil {
		return fmt.Errorf("error while storing staking genesis redelegations: %s", err)
	}

	// Save the description
	err = m.saveValidatorDescription(doc, genState.Validators)
	if err != nil {
		return fmt.Errorf("error while storing staking genesis validator descriptions: %s", err)
	}

	err = m.saveValidatorsCommissions(doc.InitialHeight, genState.Validators)
	if err != nil {
		return fmt.Errorf("error while storing staking genesis validators commissions: %s", err)
	}

	return nil
}

func (m *Module) parseGenesisTransactions(doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	var genUtilState genutiltypes.GenesisState
	err := m.cdc.UnmarshalJSON(appState[genutiltypes.ModuleName], &genUtilState)
	if err != nil {
		return fmt.Errorf("error while unmarhsaling genutil state: %s", err)
	}

	for _, genTxBz := range genUtilState.GetGenTxs() {
		// Unmarshal the transaction
		var genTx tx.Tx
		err = m.cdc.UnmarshalJSON(genTxBz, &genTx)
		if err != nil {
			return fmt.Errorf("error while unmashasling genesis tx: %s", err)
		}

		for _, msg := range genTx.GetMsgs() {
			// Handle the message properly
			createValMsg, ok := msg.(*stakingtypes.MsgCreateValidator)
			if !ok {
				continue
			}

			err = m.handleMsgCreateValidator(doc.InitialHeight, createValMsg)
			if err != nil {
				return fmt.Errorf("error while storing validators from MsgCreateValidator: %s", err)
			}
		}
	}

	return nil
}

// -------------------------------------------------------------------------------------------------------------------

// saveParams saves the given params into the database
func (m *Module) saveParams(height int64, params stakingtypes.Params) error {
	return m.db.SaveStakingParams(types.NewStakingParams(params, height))
}

// --------------------------------------------------------------------------------------------------------------------

// saveValidators stores the validators data present inside the given genesis state
func (m *Module) saveValidators(doc *tmtypes.GenesisDoc, validators stakingtypes.Validators) error {
	vals := make([]types.Validator, len(validators))
	for i, val := range validators {
		validator, err := m.convertValidator(doc.InitialHeight, val)
		if err != nil {
			return err
		}

		vals[i] = validator
	}

	return m.db.SaveValidatorsData(vals)
}

// saveValidatorDescription saves the description for the given validators
func (m *Module) saveValidatorDescription(doc *tmtypes.GenesisDoc, validators stakingtypes.Validators) error {
	for _, account := range validators {
		description := m.convertValidatorDescription(
			doc.InitialHeight,
			account.OperatorAddress,
			account.Description,
		)

		err := m.db.SaveValidatorDescription(description)
		if err != nil {
			return err
		}
	}

	return nil
}

// --------------------------------------------------------------------------------------------------------------------

// saveDelegations stores the delegations data present inside the given genesis state
func (m *Module) saveDelegations(doc *tmtypes.GenesisDoc, genState stakingtypes.GenesisState) error {
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

	return m.db.SaveDelegations(delegations)
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
func (m *Module) saveUnbondingDelegations(doc *tmtypes.GenesisDoc, genState stakingtypes.GenesisState) error {
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

	return m.db.SaveUnbondingDelegations(unbondingDelegations)
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
func (m *Module) saveRedelegations(doc *tmtypes.GenesisDoc, genState stakingtypes.GenesisState) error {
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

	return m.db.SaveRedelegations(redelegations)
}

// --------------------------------------------------------------------------------------------------------------------

// saveValidatorsCommissions save the initial commission for each validator
func (m *Module) saveValidatorsCommissions(height int64, validators stakingtypes.Validators) error {
	for _, account := range validators {
		err := m.db.SaveValidatorCommission(types.NewValidatorCommission(
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
