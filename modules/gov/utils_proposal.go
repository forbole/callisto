package gov

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govtypesv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/rs/zerolog/log"

	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	proposaltypes "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"google.golang.org/grpc/codes"

	"github.com/forbole/bdjuno/v4/types"

	gov "github.com/cosmos/cosmos-sdk/x/gov/types"
)

// UpdateProposalStatus queries the latest details of given proposal ID, updates it's status
// in database and handles changes if the proposal has been passed.
func (m *Module) UpdateProposalStatus(height int64, id uint64) error {
	// Get the proposal
	proposal, err := m.source.Proposal(height, id)
	if err != nil {
		// Check if proposal exist on the chain
		if strings.Contains(err.Error(), codes.NotFound.String()) && strings.Contains(err.Error(), "doesn't exist") {
			// Handle case when a proposal is deleted from the chain (did not pass deposit period)
			return m.updateDeletedProposalStatus(id)
		}

		return fmt.Errorf("error while getting proposal: %s", err)
	}

	err = m.updateProposalStatus(proposal)
	if err != nil {
		return fmt.Errorf("error while updating proposal status: %s", err)
	}

	err = m.handlePassedProposal(proposal, height)
	if err != nil {
		return fmt.Errorf("error while handling passed proposals: %s", err)
	}

	return nil
}

// updateProposalStatus updates given proposal status
func (m *Module) updateProposalStatus(proposal *govtypesv1.Proposal) error {
	return m.db.UpdateProposal(
		types.NewProposalUpdate(
			proposal.Id,
			proposal.Status.String(),
			proposal.VotingStartTime,
			proposal.VotingEndTime,
		),
	)
}

// UpdateProposalsStakingPoolSnapshot updates
// staking pool snapshots for active proposals
func (m *Module) UpdateProposalsStakingPoolSnapshot() error {
	log.Debug().Str("module", "gov").Msg("refreshing proposal staking pool snapshots")
	block, err := m.db.GetLastBlockHeightAndTimestamp()
	if err != nil {
		return err
	}

	ids, err := m.db.GetOpenProposalsIds(block.BlockTimestamp)
	if err != nil {
		log.Error().Err(err).Str("module", "gov").Msg("error while getting open proposals ids")
	}

	for _, proposalID := range ids {
		err = m.UpdateProposalStakingPoolSnapshot(block.Height, proposalID)
		if err != nil {
			return fmt.Errorf("error while updating proposal %d staking pool snapshots: %s", proposalID, err)
		}
	}

	return nil
}

// UpdateProposalStakingPoolSnapshot updates the staking pool snapshot associated with the gov
// proposal having the provided id
func (m *Module) UpdateProposalStakingPoolSnapshot(height int64, proposalID uint64) error {
	pool, err := m.stakingModule.GetStakingPoolSnapshot(height)
	if err != nil {
		return fmt.Errorf("error while getting staking pool: %s", err)
	}

	return m.db.SaveProposalStakingPoolSnapshot(
		types.NewProposalStakingPoolSnapshot(proposalID, pool),
	)
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
			stored.ID,
			types.ProposalStatusInvalid,
			stored.VotingStartTime,
			stored.VotingEndTime,
		),
	)
}

// handleParamChangeProposal updates params to the corresponding modules if a ParamChangeProposal has passed
func (m *Module) handleParamChangeProposal(height int64, moduleName string) (err error) {
	switch moduleName {
	case distrtypes.ModuleName:
		err = m.distrModule.UpdateParams(height)
		if err != nil {
			return fmt.Errorf("error while updating ParamChangeProposal %s params : %s", distrtypes.ModuleName, err)
		}
	case gov.ModuleName:
		err = m.UpdateParams(height)
		if err != nil {
			return fmt.Errorf("error while updating ParamChangeProposal %s params : %s", gov.ModuleName, err)
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

	return nil
}

// UpdateProposalsTallyResults updates the tally for active proposals
func (m *Module) UpdateProposalsTallyResults() error {
	log.Debug().Str("module", "gov").Msg("refreshing proposal tally results")
	block, err := m.db.GetLastBlockHeightAndTimestamp()
	if err != nil {
		return err
	}

	ids, err := m.db.GetOpenProposalsIds(block.BlockTimestamp)
	if err != nil {
		log.Error().Err(err).Str("module", "gov").Msg("error while getting open proposals ids")
	}

	for _, proposalID := range ids {
		err = m.UpdateProposalTallyResult(proposalID, block.Height)
		if err != nil {
			return fmt.Errorf("error while updating proposal %d tally result : %s", proposalID, err)
		}
	}

	return nil
}

// UpdateProposalTallyResult updates the tally result associated with the given proposal ID
func (m *Module) UpdateProposalTallyResult(proposalID uint64, height int64) error {
	result, err := m.source.TallyResult(height, proposalID)
	if err != nil {
		return fmt.Errorf("error while getting tally result: %s", err)
	}

	return m.db.SaveTallyResults([]types.TallyResult{
		types.NewTallyResult(
			proposalID,
			result.YesCount,
			result.AbstainCount,
			result.NoCount,
			result.NoWithVetoCount,
			height,
		),
	})
}

func (m *Module) handlePassedProposal(proposal *govtypesv1.Proposal, height int64) error {
	if proposal.Status != govtypesv1.StatusPassed {
		// If proposal status is not passed, do nothing
		return nil
	}

	for _, msg := range proposal.Messages {
		var sdkMsg sdk.Msg
		err := m.cdc.UnpackAny(msg, &sdkMsg)
		if err != nil {
			return fmt.Errorf("error while unpacking proposal message: %s", err)
		}

		switch msg := sdkMsg.(type) {
		case *govtypesv1.MsgExecLegacyContent:
			err := m.handlePassedV1Beta1Proposal(proposal, msg, height)
			if err != nil {
				return err
			}

		default:
			err := m.handlePassedV1Proposal(proposal, msg, height)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// handlePassedV1Proposal handles a passed proposal that contains a v1 message (new version)
func (m *Module) handlePassedV1Proposal(proposal *govtypesv1.Proposal, msg sdk.Msg, height int64) error {
	switch msg := msg.(type) {
	case *upgradetypes.MsgSoftwareUpgrade:
		// Store software upgrade plan while SoftwareUpgradeProposal passed
		err := m.db.SaveSoftwareUpgradePlan(proposal.Id, msg.Plan, height)
		if err != nil {
			return fmt.Errorf("error while storing software upgrade plan: %s", err)
		}

	case *upgradetypes.MsgCancelUpgrade:
		// Delete software upgrade plan while CancelSoftwareUpgradeProposal passed
		err := m.db.DeleteSoftwareUpgradePlan(proposal.Id)
		if err != nil {
			return fmt.Errorf("error while deleting software upgrade plan: %s", err)
		}

	default:
		// Try to see if it's a param change proposal. This should be handled as last case
		// because it's the most generic one
		subspace, ok := getParamChangeSubspace(msg)
		if ok {
			err := m.handleParamChangeProposal(height, subspace)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// getParamChangeSubspace returns the subspace of the param change proposal, if any.
// If the message is not a param change proposal, it returns false
func getParamChangeSubspace(msg sdk.Msg) (string, bool) {
	switch msg.(type) {
	case *distrtypes.MsgUpdateParams:
		return distrtypes.ModuleName, true
	case *govtypesv1.MsgUpdateParams:
		return gov.ModuleName, true
	case *minttypes.MsgUpdateParams:
		return minttypes.ModuleName, true
	case *slashingtypes.MsgUpdateParams:
		return slashingtypes.ModuleName, true
	case *stakingtypes.MsgUpdateParams:
		return stakingtypes.ModuleName, true

	default:
		return "", false
	}
}

// handlePassedV1Beta1Proposal handles a passed proposal with a v1beta1 message (legacy)
func (m *Module) handlePassedV1Beta1Proposal(proposal *govtypesv1.Proposal, msg *govtypesv1.MsgExecLegacyContent, height int64) error {
	// Unpack proposal
	var content govtypesv1beta1.Content
	var protoCodec codec.ProtoCodec
	err := protoCodec.UnpackAny(msg.Content, &content)
	if err != nil {
		return fmt.Errorf("error while handling ParamChangeProposal: %s", err)
	}

	switch p := content.(type) {
	case *proposaltypes.ParameterChangeProposal:
		// Update params while ParameterChangeProposal passed
		for _, change := range p.Changes {
			err = m.handleParamChangeProposal(height, change.Subspace)
			if err != nil {
				return fmt.Errorf("error while updating params from ParamChangeProposal: %s", err)
			}
		}
	case *upgradetypes.SoftwareUpgradeProposal:
		// Store software upgrade plan while SoftwareUpgradeProposal passed
		err = m.db.SaveSoftwareUpgradePlan(proposal.Id, p.Plan, height)
		if err != nil {
			return fmt.Errorf("error while storing software upgrade plan: %s", err)
		}
	case *upgradetypes.CancelSoftwareUpgradeProposal:
		// Delete software upgrade plan while CancelSoftwareUpgradeProposal passed
		err = m.db.DeleteSoftwareUpgradePlan(proposal.Id)
		if err != nil {
			return fmt.Errorf("error while deleting software upgrade plan: %s", err)
		}
	}
	return nil
}
