package database

import (
	"encoding/json"
	"fmt"

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

// SaveCcvProviderParams saves ccv provider params for the given height
func (db *Db) SaveCcvProviderParams(params *types.CcvProviderParams) error {
	paramsBz, err := json.Marshal(&params.Params)
	if err != nil {
		return err
	}

	stmt := `
INSERT INTO ccv_provider_params (params, height) 
VALUES ($1, $2)
ON CONFLICT (one_row_id) DO UPDATE 
    SET params = excluded.params, 
        height = excluded.height
WHERE ccv_provider_params.height <= excluded.height`

	_, err = db.SQL.Exec(stmt, string(paramsBz), params.Height)
	if err != nil {
		return fmt.Errorf("error while storing ccv provider params: %s", err)
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

// SaveCcvProviderChain saves the ccv provider chain state for the given height
func (db *Db) SaveCcvProviderChain(providerChain *types.CcvProviderChain) error {
	consumerStates, err := json.Marshal(&providerChain.ConsumerStates)
	if err != nil {
		return err
	}
	unbondingOps, err := json.Marshal(&providerChain.UnbondingOps)
	if err != nil {
		return err
	}
	matureUnbondingOps, err := json.Marshal(&providerChain.MatureUnbondingOps)
	if err != nil {
		return err
	}
	valsetUpdateIdToHeight, err := json.Marshal(&providerChain.ValsetUpdateIdToHeight)
	if err != nil {
		return err
	}
	consumerAdditionProposals, err := json.Marshal(&providerChain.ConsumerAdditionProposals)
	if err != nil {
		return err
	}
	consumerRemovalProposals, err := json.Marshal(&providerChain.ConsumerRemovalProposals)
	if err != nil {
		return err
	}
	validatorConsumerPubkeys, err := json.Marshal(&providerChain.ValidatorConsumerPubkeys)
	if err != nil {
		return err
	}
	validatorsByConsumerAddr, err := json.Marshal(&providerChain.ValidatorsByConsumerAddr)
	if err != nil {
		return err
	}
	consumerAddrsToPrune, err := json.Marshal(&providerChain.ConsumerAddrsToPrune)
	if err != nil {
		return err
	}

	stmt := `
INSERT INTO ccv_provider_chain (valset_update_id, consumer_states,unbonding_ops, 
	mature_unbonding_ops, valset_update_id_to_height, consumer_addition_proposals,
	consumer_removal_proposals, validator_consumer_pubkeys, 
	validators_by_consumer_addr, consumer_addrs_to_prune,  height) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
ON CONFLICT DO NOTHING`

	_, err = db.SQL.Exec(stmt, providerChain.ValsetUpdateID, string(consumerStates),
		string(unbondingOps), string(matureUnbondingOps), string(valsetUpdateIdToHeight),
		string(consumerAdditionProposals), string(consumerRemovalProposals),
		string(validatorConsumerPubkeys), string(validatorsByConsumerAddr),
		string(consumerAddrsToPrune), providerChain.Height)

	if err != nil {
		return fmt.Errorf("error while storing ccv provider chain info: %s", err)
	}

	return nil
}
