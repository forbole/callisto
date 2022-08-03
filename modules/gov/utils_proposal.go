package gov

import (
	"fmt"
	"strings"

	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	proposaltypes "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"

	"google.golang.org/grpc/codes"

	"github.com/forbole/bdjuno/v3/types"

	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

func (m *Module) UpdateProposal(height int64, id uint64) error {
	// Get the proposal
	proposal, err := m.source.Proposal(height, id)
	if err != nil {
		if strings.Contains(err.Error(), codes.NotFound.String()) {
			// Handle case when a proposal is deleted from the chain (did not pass deposit period)
			return m.updateDeletedProposalStatus(id)
		}

		return fmt.Errorf("error while getting proposal: %s", err)
	}

	err = m.handleParamChangeProposal(height, proposal)
	if err != nil {
		return fmt.Errorf("error while updating params from ParamChangeProposal: %s", err)
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
	return nil
}

func (m *Module) UpdateProposalSnapshots(height int64, blockVals *tmctypes.ResultValidators, id uint64) error {
	err := m.updateProposalStakingPoolSnapshot(height, id)
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

// handleParamChangeProposal updates params to the corresponding modules if a ParamChangeProposal has passed
func (m *Module) handleParamChangeProposal(height int64, proposal govtypes.Proposal) error {
	if proposal.Status != govtypes.StatusPassed {
		// If the status of ParamChangeProposal is not passed, do nothing
		return nil
	}

	var content govtypes.Content
	err := m.db.EncodingConfig.Marshaler.UnpackAny(proposal.Content, &content)
	if err != nil {
		return fmt.Errorf("error while handling ParamChangeProposal: %s", err)
	}

	paramChangeProposal, ok := content.(*proposaltypes.ParameterChangeProposal)
	if !ok {
		return nil
	}

	for _, change := range paramChangeProposal.Changes {
		// Update the params for corresponding modules
		switch change.Subspace {
		case distrtypes.ModuleName:
			err = m.distrModule.UpdateParams(height)
			if err != nil {
				return fmt.Errorf("error while updating ParamChangeProposal %s params : %s", distrtypes.ModuleName, err)
			}
		case govtypes.ModuleName:
			err = m.UpdateParams(height)
			if err != nil {
				return fmt.Errorf("error while updating ParamChangeProposal %s params : %s", govtypes.ModuleName, err)
			}
		case minttypes.ModuleName:
			err = m.mintModule.UpdateParams(height)
			if err != nil {
				return fmt.Errorf("error while updating ParamChangeProposal %s params : %s", minttypes.ModuleName, err)
			}

			// Update the inflation
			err = m.mintModule.UpdateInflation()
			if err != nil {
				return fmt.Errorf("error while updating inflation with ParamChangeProposal: %s", err)
			}
		case slashingtypes.ModuleName:
			err = m.slashingModule.UpdateParams(height)
			if err != nil {
				return fmt.Errorf("error while updating ParamChangeProposal %s params : %s", slashingtypes.ModuleName, err)
			}
		case stakingtypes.ModuleName:
			err = m.stakingModule.UpdateParams(height)
			if err != nil {
				return fmt.Errorf("error while updating ParamChangeProposal %s params : %s", stakingtypes.ModuleName, err)
			}
		}
	}
	return nil
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
			result.Yes.String(),
			result.Abstain.String(),
			result.No.String(),
			result.NoWithVeto.String(),
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

		return nil
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

	votingPowers, err := m.stakingModule.GetValidatorsVotingPowers(height, blockVals)
	if err != nil {
		return fmt.Errorf("error while getting validators voting powers: %s", err)
	}

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
