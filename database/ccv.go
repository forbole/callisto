package database

import (
	"encoding/json"
	"fmt"

	ccvconsumertypes "github.com/cosmos/interchain-security/x/ccv/consumer/types"
	"github.com/forbole/bdjuno/v4/types"
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

// SaveProposals allows to save for the given height the given total amount of coins
func (db *Db) SaveCcvProposals(height int64, proposals []types.CcvProposal) error {
	if len(proposals) == 0 {
		return nil
	}

	var accounts []types.Account

	proposalsQuery := `
INSERT INTO proposal(
	id, title, description, content, proposer_address, proposal_route, proposal_type, status, 
    submit_time, deposit_end_time, voting_start_time, voting_end_time
) VALUES`
	var proposalsParams []interface{}

	for i, proposal := range proposals {
		// Prepare the account query
		accounts = append(accounts, types.NewAccount(proposal.Proposer))

		// Prepare the proposal query
		vi := i * 12
		proposalsQuery += fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d),",
			vi+1, vi+2, vi+3, vi+4, vi+5, vi+6, vi+7, vi+8, vi+9, vi+10, vi+11, vi+12)

		contentBz, err := json.Marshal(&proposal.Content)
		if err != nil {
			return err
		}

		proposalsParams = append(proposalsParams,
			proposal.ProposalID,
			proposal.Content.Title,
			proposal.Content.Description,
			string(contentBz),
			proposal.Proposer,
			proposal.ProposalRoute,
			proposal.ProposalType,
			proposal.Status,
			proposal.SubmitTime,
			proposal.DepositEndTime,
			proposal.VotingStartTime,
			proposal.VotingEndTime,
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
