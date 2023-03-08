package database

import (
	"encoding/json"
	"fmt"

	ccvconsumertypes "github.com/cosmos/interchain-security/x/ccv/consumer/types"
	dbtypes "github.com/forbole/bdjuno/v4/database/types"
	"github.com/forbole/bdjuno/v4/types"
	"github.com/lib/pq"
)

// SaveCcvProviderParams saves the ccv provider params for the given height
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

// SaveCcvConsumerParams saves the ccv consumer params for the given height
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

// SaveCcvConsumerChain saves the ccv consumer chain state for the given height
func (db *Db) SaveCcvConsumerChain(consumerChain *types.CcvConsumerChain) error {
	providerClientState, err := json.Marshal(&consumerChain.ProviderClientState)
	if err != nil {
		return err
	}
	providerConsensusState, err := json.Marshal(&consumerChain.ProviderConsensusState)
	if err != nil {
		return err
	}
	maturingPackets, err := json.Marshal(&consumerChain.MaturingPackets)
	if err != nil {
		return err
	}
	initialValSet, err := json.Marshal(&consumerChain.InitialValSet)
	if err != nil {
		return err
	}
	heightToValsetUpdateID, err := json.Marshal(&consumerChain.HeightToValsetUpdateID)
	if err != nil {
		return err
	}
	outstandingDowntimeSlashing, err := json.Marshal(&consumerChain.OutstandingDowntimeSlashing)
	if err != nil {
		return err
	}
	pendingConsumerPackets, err := json.Marshal(&consumerChain.PendingConsumerPackets)
	if err != nil {
		return err
	}
	lastTransmissionBlockHeight, err := json.Marshal(&consumerChain.LastTransmissionBlockHeight)
	if err != nil {
		return err
	}

	stmt := `
INSERT INTO ccv_consumer_chain (provider_client_id, provider_channel_id, new_chain, provider_client_state,
	provider_consensus_state, maturing_packets, initial_val_set, height_to_valset_update_id, 
	outstanding_downtime_slashing, pending_consumer_packets, last_transmission_block_height, height) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
ON CONFLICT (provider_client_id) DO UPDATE 
	SET provider_channel_id = excluded.provider_channel_id,
		new_chain = excluded.new_chain,
		provider_client_state = excluded.provider_client_state,
		provider_consensus_state = excluded.provider_consensus_state,
		maturing_packets = excluded.maturing_packets,
		initial_val_set = excluded.initial_val_set,
		height_to_valset_update_id = excluded.height_to_valset_update_id,
		outstanding_downtime_slashing = excluded.outstanding_downtime_slashing,
		pending_consumer_packets = excluded.pending_consumer_packets,
		last_transmission_block_height = excluded.last_transmission_block_height,
		height = excluded.height
WHERE ccv_consumer_chain.height <= excluded.height`

	_, err = db.SQL.Exec(stmt, consumerChain.ProviderClientID, consumerChain.ProviderChannelID,
		consumerChain.NewChain, string(providerClientState), string(providerConsensusState),
		string(maturingPackets), string(initialValSet), string(heightToValsetUpdateID),
		string(outstandingDowntimeSlashing), string(pendingConsumerPackets),
		string(lastTransmissionBlockHeight), consumerChain.Height)

	if err != nil {
		return fmt.Errorf("error while storing ccv consumer chain state info: %s", err)
	}

	return nil
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

	_, err = db.SQL.Exec(stmt, providerChain.ValsetUpdateID, string(consumerStates), string(unbondingOps), string(matureUnbondingOps),
		string(valsetUpdateIdToHeight), string(consumerAdditionProposals), string(consumerRemovalProposals),
		string(validatorConsumerPubkeys), string(validatorsByConsumerAddr),
		string(consumerAddrsToPrune), providerChain.Height)

	if err != nil {
		return fmt.Errorf("error while storing ccv provider chain state info: %s", err)
	}

	return nil
}

// SaveNextFeeDistributionEstimate allows to store the next fee distribution estimate
func (db *Db) SaveNextFeeDistributionEstimate(height int64, fee ccvconsumertypes.NextFeeDistributionEstimate) error {
	// Store the accounts
	var accounts []types.Account
	accounts = append(accounts, types.NewAccount(fee.ToProvider), types.NewAccount(fee.ToConsumer))
	err := db.SaveAccounts(accounts)
	if err != nil {
		return fmt.Errorf("error while storing provider and consumer fee distr accounts: %s", err)
	}

	stmt := `
INSERT INTO ccv_fee_distribution(current_height, last_height, next_height, distribution_fraction,
	total, to_provider, to_consumer, height) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
ON CONFLICT ON CONSTRAINT unique_provider_consumer_fee_distribution DO UPDATE 
    SET current_height = excluded.current_height,
		last_height = excluded.last_height,
		next_height = excluded.next_height,
		distribution_fraction = excluded.distribution_fraction,
		total = excluded.total,
        height = excluded.height
WHERE ccv_fee_distribution.height <= excluded.height`

	_, err = db.SQL.Exec(stmt, fee.CurrentHeight, fee.LastHeight, fee.NextHeight,
		fee.DistributionFraction, fee.Total, fee.ToProvider, fee.ToConsumer, height)
	if err != nil {
		return fmt.Errorf("error while saving next fee distribution estimate: %s", err)
	}

	return nil
}

// SaveCcvProposals allows to save ccv proposals for the given height
func (db *Db) SaveCcvProposals(height int64, proposals []types.CcvProposal) error {
	if len(proposals) == 0 {
		return nil
	}

	var accounts []types.Account

	proposalsQuery := `
INSERT INTO ccv_proposal(
	id, title, description, chain_id, genesis_hash, binary_hash, proposal_type, proposal_route, 
    spawn_time, stop_time, initial_height, unbonding_period, ccv_timeout_period, transfer_timeout_period, 
	consumer_redistribution_fraction, blocks_per_distribution_transmission, historical_entries, status, 
	submit_time, proposer_address, height
) VALUES`
	var proposalsParams []interface{}

	for i, proposal := range proposals {
		// Prepare the account query
		accounts = append(accounts, types.NewAccount(proposal.Proposer))
		initialHeight, err := json.Marshal(&proposal.InitialHeight)
		if err != nil {
			return err
		}
		// Prepare the proposal query
		vi := i * 21
		proposalsQuery += fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d),",
			vi+1, vi+2, vi+3, vi+4, vi+5, vi+6, vi+7, vi+8, vi+9, vi+10, vi+11, vi+12, vi+13,
			vi+14, vi+15, vi+16, vi+17, vi+18, vi+19, vi+20, vi+21)

		proposalsParams = append(proposalsParams,
			proposal.ProposalID,
			proposal.Title,
			proposal.Description,
			proposal.ChainID,
			proposal.GenesisHash,
			proposal.BinaryHash,
			proposal.ProposalType,
			proposal.ProposalRoute,
			proposal.SpawnTime,
			proposal.StopTime,
			string(initialHeight),
			proposal.UnbondingPeriod,
			proposal.CcvTimeoutPeriod,
			proposal.TransferTimeoutPeriod,
			proposal.ConsumerRedistributionFraction,
			proposal.BlocksPerDistributionTransmission,
			proposal.HistoricalEntries,
			proposal.Status,
			proposal.SubmitTime,
			proposal.Proposer,
			height,
		)
	}

	// Store the accounts
	err := db.SaveAccounts(accounts)
	if err != nil {
		return fmt.Errorf("error while storing ccv proposal proposers accounts: %s", err)
	}

	// Store the proposals
	proposalsQuery = proposalsQuery[:len(proposalsQuery)-1] // Remove trailing ","
	proposalsQuery += " ON CONFLICT DO NOTHING"
	_, err = db.SQL.Exec(proposalsQuery, proposalsParams...)
	if err != nil {
		return fmt.Errorf("error while storing ccv proposals: %s", err)
	}

	return nil
}

// SaveCcvDeposits allows to save multiple ccv proposal deposits
func (db *Db) SaveCcvDeposits(deposits []types.Deposit) error {
	if len(deposits) == 0 {
		return nil
	}

	query := `INSERT INTO ccv_proposal_deposit (proposal_id, depositor_address, amount, height) VALUES `
	var param []interface{}

	for i, deposit := range deposits {
		vi := i * 4
		query += fmt.Sprintf("($%d,$%d,$%d,$%d),", vi+1, vi+2, vi+3, vi+4)
		param = append(param, deposit.ProposalID,
			deposit.Depositor,
			pq.Array(dbtypes.NewDbCoins(deposit.Amount)),
			deposit.Height,
		)
	}
	query = query[:len(query)-1] // Remove trailing ","
	query += `
ON CONFLICT ON CONSTRAINT unique_deposit DO UPDATE
	SET amount = excluded.amount,
		height = excluded.height
WHERE ccv_proposal_deposit.height <= excluded.height`
	_, err := db.SQL.Exec(query, param...)
	if err != nil {
		return fmt.Errorf("error while storing ccv proposal deposits: %s", err)
	}

	return nil
}
