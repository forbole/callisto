package database

import (
	"encoding/json"
	"fmt"

	dbtypes "github.com/forbole/bdjuno/v4/database/types"
	"github.com/forbole/bdjuno/v4/types"
)

// SaveCcvConsumerParams saves ccv consumer params for the given height
func (db *Db) SaveCcvConsumerParams(params *types.CcvConsumerParams) error {
	paramsBz, err := json.Marshal(&params.Params)
	if err != nil {
		return err
	}

	stmt := `
INSERT INTO ccv_consumer_params (params, height) 
VALUES ($1, $2)
ON CONFLICT (one_row_id) DO UPDATE 
    SET params = excluded.params, 
        height = excluded.height
WHERE ccv_consumer_params.height <= excluded.height`

	_, err = db.SQL.Exec(stmt, string(paramsBz), params.Height)
	if err != nil {
		return fmt.Errorf("error while storing ccv consumer params: %s", err)
	}

	return nil
}

// SaveCcvConsumerChains saves ccv consumer chains for the given height
func (db *Db) SaveCcvConsumerChains(consumerChains []*types.CcvConsumerChain) error {
	if len(consumerChains) == 0 {
		return nil
	}

	stmt := `
INSERT INTO ccv_consumer_chain (provider_client_id, provider_channel_id, chain_id, provider_client_state,
	provider_consensus_state, initial_val_set, height) 
VALUES `
	var consumerChainsList []interface{}

	for i, consumerChain := range consumerChains {

		// Prepare the consumer chains query
		vi := i * 7
		stmt += fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d,$%d),",
			vi+1, vi+2, vi+3, vi+4, vi+5, vi+6, vi+7)

		providerClientState, err := json.Marshal(&consumerChain.ProviderClientState)
		if err != nil {
			return err
		}
		providerConsensusState, err := json.Marshal(&consumerChain.ProviderConsensusState)
		if err != nil {
			return err
		}
		initialValSet, err := json.Marshal(&consumerChain.InitialValSet)
		if err != nil {
			return err
		}

		consumerChainsList = append(consumerChainsList,
			consumerChain.ProviderClientID, consumerChain.ProviderChannelID,
			consumerChain.ChainID, string(providerClientState), string(providerConsensusState),
			string(initialValSet), consumerChain.Height,
		)
	}

	// Store the consumer chains
	stmt = stmt[:len(stmt)-1] // Remove trailing ","
	stmt += `
ON CONFLICT ON CONSTRAINT unique_provider_id DO UPDATE
	SET provider_channel_id = excluded.provider_channel_id,
		chain_id = excluded.chain_id,
		provider_client_state = excluded.provider_client_state,
		provider_consensus_state = excluded.provider_consensus_state,
		initial_val_set = excluded.initial_val_set,
		height = excluded.height
WHERE ccv_consumer_chain.height <= excluded.height`
	_, err := db.SQL.Exec(stmt, consumerChainsList...)
	if err != nil {
		return fmt.Errorf("error while storing ccv consumer chain info: %s", err)
	}

	return nil
}

// DeleteConsumerChainFromDB removes given ccv consumer chain from database
func (db *Db) DeleteConsumerChainFromDB(height int64, chainID string) error {
	_, err := db.SQL.Exec(`DELETE FROM ccv_consumer_chain WHERE chain_id = $1`, chainID)
	return fmt.Errorf("error while removing consumer chain with chain_id %s: %s", chainID, err)
}

// GetValidatorsConsensusAddress returns all validators consensus address stored inside the database.
func (db *Db) GetValidatorsConsensusAddress() ([]dbtypes.ValidatorAddressRow, error) {
	stmt := `SELECT consensus_address FROM validator`

	var rows []dbtypes.ValidatorAddressRow
	err := db.Sqlx.Select(&rows, stmt)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

// StoreCCvValidators stores ccv validators operators address inside the database
// To create relationship between provider and consumer chain
func (db *Db) StoreCCvValidators(ccvValidators []types.CCVValidator) error {
	if len(ccvValidators) == 0 {
		return nil
	}

	stmt := `
INSERT INTO ccv_validator (consumer_consensus_address, provider_consensus_address, height) 
VALUES `

	var ccvValidatorsList []interface{}
	for i, ccvValidator := range ccvValidators {
		vi := i * 3
		stmt += fmt.Sprintf("($%d,$%d,$%d),", vi+1, vi+2, vi+3)

		ccvValidatorsList = append(ccvValidatorsList,
			ccvValidator.ConsumerConsensusAddress,
			ccvValidator.ProviderConsensusAddress,
			ccvValidator.Height,
		)
	}

	// Store the ccv validators
	stmt = stmt[:len(stmt)-1] // Remove trailing ","
	stmt += " ON CONFLICT DO NOTHING"
	_, err := db.SQL.Exec(stmt, ccvValidatorsList...)
	if err != nil {
		return fmt.Errorf("error while storing ccv validators info: %s", err)
	}

	return nil
}
