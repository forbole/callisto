package gov

import (
	"fmt"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"

	"google.golang.org/grpc/codes"

	"github.com/forbole/bdjuno/types"

	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

const (
	ErrProposalNotFound = "rpc error: code = %s desc = rpc error: code = %s desc = proposal %d doesn't exist: key not found"
)

func (m *Module) UpdateProposal(height int64, blockVals *tmctypes.ResultValidators, id uint64) error {
	// Get the proposal
	proposal, err := m.source.Proposal(height, id)
	if err != nil {
		// Get the error code
		var code string
		_, err := fmt.Sscanf(err.Error(), ErrProposalNotFound, &code, &code, &id)
		if err != nil {
			return err
		}

		if code == codes.NotFound.String() {
			// Handle case when a proposal is deleted from the chain (did not pass deposit period)
			return m.updateDeletedProposalStatus(id)
		}

		return fmt.Errorf("error while getting proposal: %s", err)
	}

	err = m.updateProposalStatus(proposal)
	if err != nil {
		return fmt.Errorf("error while updating proposal status: %s", err)
	}

	err = m.updateProposalTallyResult(proposal)
	if err != nil {
		return fmt.Errorf("error while updating proposal tally result: %s", err)
	}

	err = m.updateAccounts(proposal)
	if err != nil {
		return fmt.Errorf("error while updating account: %s", err)
	}

	err = m.updateProposalStakingPoolSnapshot(height, id)
	if err != nil {
		return fmt.Errorf("error while updating proposal staking pool snapshot: %s", err)
	}

	err = m.updateProposalValidatorStatusesSnapshot(height, id, blockVals)
	if err != nil {
		return fmt.Errorf("error while updating proposal validator statuses snapshot: %s", err)
	}

	return nil
}

// updateDeletedProposalStatus updates the proposal having the given id by setting its status
// to the one that represents a deleted proposal
func (m *Module) updateDeletedProposalStatus(id uint64) error {
	stored, err := m.db.GetProposal(id)
	if err != nil {
		return err
	}

	return m.db.UpdateProposal(
		types.NewProposalUpdate(
			stored.ProposalID,
			types.ProposalStatusInvalid,
			stored.VotingStartTime,
			stored.VotingEndTime,
		),
	)
}

// updateProposalStatus updates the given proposal status
func (m *Module) updateProposalStatus(proposal govtypes.Proposal) error {
	return m.db.UpdateProposal(
		types.NewProposalUpdate(
			proposal.ProposalId,
			proposal.Status.String(),
			proposal.VotingStartTime,
			proposal.VotingEndTime,
		),
	)
}

// updateProposalTallyResult updates the tally result associated with the given proposal
func (m *Module) updateProposalTallyResult(proposal govtypes.Proposal) error {
	height, err := m.db.GetLastBlockHeight()
	if err != nil {
		return err
	}

	result, err := m.source.TallyResult(height, proposal.ProposalId)
	if err != nil {
		return fmt.Errorf("error while getting tally result: %s", err)
	}

	return m.db.SaveTallyResults([]types.TallyResult{
		types.NewTallyResult(
			proposal.ProposalId,
			result.Yes.Int64(),
			result.Abstain.Int64(),
			result.No.Int64(),
			result.NoWithVeto.Int64(),
			height,
		),
	})
}

// updateAccounts updates any account that might be involved in the proposal (eg. fund community recipient)
func (m *Module) updateAccounts(proposal govtypes.Proposal) error {
	content, ok := proposal.Content.GetCachedValue().(*distrtypes.CommunityPoolSpendProposal)
	if ok {
		height, err := m.db.GetLastBlockHeight()
		if err != nil {
			return fmt.Errorf("error while getting last block height: %s", err)
		}

		addresses := []string{content.Recipient}

		err = m.authModule.RefreshAccounts(height, addresses)
		if err != nil {
			return err
		}

		return m.bankModule.RefreshBalances(height, addresses)
	}
	return nil
}

// updateProposalStakingPoolSnapshot updates the staking pool snapshot associated with the gov
// proposal having the provided id
func (m *Module) updateProposalStakingPoolSnapshot(height int64, proposalID uint64) error {
	pool, err := m.stakingModule.GetStakingPool(height)
	if err != nil {
		return fmt.Errorf("error while getting staking pool: %s", err)
	}

	return m.db.SaveProposalStakingPoolSnapshot(
		types.NewProposalStakingPoolSnapshot(proposalID, pool),
	)
}

// updateProposalValidatorStatusesSnapshot updates the snapshots of the various validators for
// the proposal having the given id
func (m *Module) updateProposalValidatorStatusesSnapshot(
	height int64, proposalID uint64, blockVals *tmctypes.ResultValidators,
) error {
	validators, _, err := m.stakingModule.GetValidatorsWithStatus(height, stakingtypes.Bonded.String())
	if err != nil {
		return fmt.Errorf("error while getting validators with bonded status: %s", err)
	}

	votingPowers := m.stakingModule.GetValidatorsVotingPowers(height, blockVals)

	statuses, err := m.stakingModule.GetValidatorsStatuses(height, validators)
	if err != nil {
		return fmt.Errorf("error while getting validator statuses: %s", err)
	}

	var snapshots = make([]types.ProposalValidatorStatusSnapshot, len(validators))
	for index, validator := range validators {
		consAddr, err := validator.GetConsAddr()
		if err != nil {
			return err
		}

		status, err := findStatus(consAddr.String(), statuses)
		if err != nil {
			return fmt.Errorf("error while searching for status: %s", err)
		}

		votingPower, err := findVotingPower(consAddr.String(), votingPowers)
		if err != nil {
			return fmt.Errorf("error while searching for voting power: %s", err)
		}

		snapshots[index] = types.NewProposalValidatorStatusSnapshot(
			proposalID,
			consAddr.String(),
			votingPower.VotingPower,
			status.Status,
			status.Jailed,
			height,
		)
	}

	return m.db.SaveProposalValidatorsStatusesSnapshots(snapshots)
}

func findVotingPower(consAddr string, powers []types.ValidatorVotingPower) (types.ValidatorVotingPower, error) {
	for _, votingPower := range powers {
		if votingPower.ConsensusAddress == consAddr {
			return votingPower, nil
		}
	}
	return types.ValidatorVotingPower{}, fmt.Errorf("voting power not found for validator with consensus address %s", consAddr)
}

func findStatus(consAddr string, statuses []types.ValidatorStatus) (types.ValidatorStatus, error) {
	for _, status := range statuses {
		if status.ConsensusAddress == consAddr {
			return status, nil
		}
	}
	return types.ValidatorStatus{}, fmt.Errorf("cannot find status for validator with consensus address %s", consAddr)
}
