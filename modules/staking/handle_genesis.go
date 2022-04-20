package staking

import (
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/types/tx"

	"github.com/forbole/bdjuno/v3/types"

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
	err = m.db.SaveStakingParams(types.NewStakingParams(genState.Params, doc.InitialHeight))
	if err != nil {
		return fmt.Errorf("error while storing genesis staking params: %s", err)
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

			err = m.StoreValidatorsFromMsgCreateValidator(doc.InitialHeight, createValMsg)
			if err != nil {
				return fmt.Errorf("error while storing validators from MsgCreateValidator: %s", err)
			}
		}

	}

	return nil
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
