package gov

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/x/authz"
	"google.golang.org/grpc/codes"

	"github.com/forbole/bdjuno/v4/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govtypesv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"

	gov "github.com/cosmos/cosmos-sdk/x/gov/types"
	juno "github.com/forbole/juno/v5/types"
)

// HandleMsgExec implements modules.AuthzMessageModule
func (m *Module) HandleMsgExec(index int, _ *authz.MsgExec, _ int, executedMsg sdk.Msg, tx *juno.Tx) error {
	return m.HandleMsg(index, executedMsg, tx)
}

// HandleMsg implements modules.MessageModule
func (m *Module) HandleMsg(index int, msg sdk.Msg, tx *juno.Tx) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	switch cosmosMsg := msg.(type) {
	case *govtypesv1beta1.MsgSubmitProposal:
		return m.handleMsgLegacySubmitProposal(tx, index, cosmosMsg)

	case *govtypesv1.MsgSubmitProposal:
		return m.handleMsgSubmitProposal(tx, index, cosmosMsg)

	case *govtypesv1.MsgDeposit:
		return m.handleMsgDeposit(tx, cosmosMsg)

	case *govtypesv1.MsgVote:
		return m.handleMsgVote(tx, cosmosMsg)

	case *govtypesv1.MsgVoteWeighted:
		return m.handleMsgVoteWeighted(tx, cosmosMsg)
	}

	return nil
}

// saveProposalAndDeposit allows to properly get and store a proposal and its initial deposit
func (m *Module) saveProposalAndDeposit(tx *juno.Tx, index int, proposer string, initialDeposit []sdk.Coin) (types.Proposal, error) {
	// Get the proposal id
	event, err := tx.FindEventByType(index, gov.EventTypeSubmitProposal)
	if err != nil {
		return types.Proposal{}, fmt.Errorf("error while searching for EventTypeSubmitProposal: %s", err)
	}

	id, err := tx.FindAttributeByKey(event, gov.AttributeKeyProposalID)
	if err != nil {
		return types.Proposal{}, fmt.Errorf("error while searching for AttributeKeyProposalID: %s", err)
	}

	proposalID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return types.Proposal{}, fmt.Errorf("error while parsing proposal id: %s", err)
	}

	// Get the proposal
	proposal, err := m.source.Proposal(tx.Height, proposalID)
	if err != nil {
		if strings.Contains(err.Error(), codes.NotFound.String()) {
			// query the proposal details using the latest height stored in db
			// to fix the rpc error returning code = NotFound desc = proposal x doesn't exist
			block, err := m.db.GetLastBlockHeightAndTimestamp()
			if err != nil {
				return types.Proposal{}, fmt.Errorf("error while getting latest block height: %s", err)
			}
			proposal, err = m.source.Proposal(block.Height, proposalID)
			if err != nil {
				return types.Proposal{}, fmt.Errorf("error while getting proposal: %s", err)
			}
		} else {
			return types.Proposal{}, fmt.Errorf("error while getting proposal: %s", err)
		}
	}

	// For backward-compatibility we use the summary as the proposal description, and we don't consider the summary
	var summary = ""
	var description = proposal.Summary

	// If the user has set a long-text description, then we use the metadata description
	// as the proposal description, and the summary as the summary (as it should be)
	metadataDescription, err := GetDescriptionFromMetadata(proposal.Metadata)
	if err != nil {
		return types.Proposal{}, fmt.Errorf("error while getting proposal metadata: %s", err)
	}

	if metadataDescription != "" {
		summary = proposal.Summary
		description = metadataDescription
	}

	// Store the proposal
	proposalObj := types.NewProposal(
		proposal.Id,
		proposal.Title,
		summary,
		description,
		proposal.Metadata,
		proposal.Messages,
		proposal.Status.String(),
		*proposal.SubmitTime,
		*proposal.DepositEndTime,
		proposal.VotingStartTime,
		proposal.VotingEndTime,
		proposal.Proposer,
	)

	err = m.db.SaveProposals([]types.Proposal{proposalObj})
	if err != nil {
		return types.Proposal{}, err
	}

	txTimestamp, err := time.Parse(time.RFC3339, tx.Timestamp)
	if err != nil {
		return types.Proposal{}, fmt.Errorf("error while parsing time: %s", err)
	}

	// Store the deposit
	deposit := types.NewDeposit(proposal.Id, proposer, initialDeposit, txTimestamp, tx.TxHash, tx.Height)
	return proposalObj, m.db.SaveDeposits([]types.Deposit{deposit})
}

// handleMsgSubmitProposal allows to properly handle a v1beta1.MsgSubmitProposal
func (m *Module) handleMsgLegacySubmitProposal(tx *juno.Tx, index int, msg *govtypesv1beta1.MsgSubmitProposal) error {
	_, err := m.saveProposalAndDeposit(tx, index, msg.Proposer, msg.InitialDeposit)
	return err
}

// handleMsgSubmitProposal allows to properly handle a v1.MsgSubmitProposal
func (m *Module) handleMsgSubmitProposal(tx *juno.Tx, index int, msg *govtypesv1.MsgSubmitProposal) error {
	proposal, err := m.saveProposalAndDeposit(tx, index, msg.Proposer, msg.InitialDeposit)
	if err != nil {
		return err
	}

	var addresses []types.Account
	for _, msg := range proposal.Messages {
		var sdkMsg sdk.Msg
		err := m.cdc.UnpackAny(msg, &sdkMsg)
		if err != nil {
			return fmt.Errorf("error while unpacking proposal message: %s", err)
		}

		switch msg := sdkMsg.(type) {
		case *distrtypes.MsgCommunityPoolSpend:
			addresses = append(addresses, types.NewAccount(msg.Recipient))
		case *govtypesv1.MsgExecLegacyContent:
			content, ok := msg.Content.GetCachedValue().(*distrtypes.CommunityPoolSpendProposal)
			if ok {
				addresses = append(addresses, types.NewAccount(content.Recipient))
			}
		}
	}

	return m.db.SaveAccounts(addresses)
}

// handleMsgDeposit allows to properly handle a MsgDeposit
func (m *Module) handleMsgDeposit(tx *juno.Tx, msg *govtypesv1.MsgDeposit) error {
	govDeposit, err := m.source.ProposalDeposit(tx.Height, msg.ProposalId, msg.Depositor)
	if err != nil {
		return fmt.Errorf("error while getting proposal deposit: %s", err)
	}
	txTimestamp, err := time.Parse(time.RFC3339, tx.Timestamp)
	if err != nil {
		return fmt.Errorf("error while parsing time: %s", err)
	}

	// Save the deposit
	deposit := types.NewDeposit(govDeposit.ProposalId, govDeposit.Depositor, govDeposit.Amount, txTimestamp, tx.TxHash, tx.Height)
	err = m.db.SaveDeposits([]types.Deposit{deposit})
	if err != nil {
		return err
	}

	// Update the proposal status (in case the deposit made the proposal enter the voting phase)
	return m.UpdateProposalStatus(tx.Height, msg.ProposalId)
}

// handleMsgVote allows to properly handle a MsgVote
func (m *Module) handleMsgVote(tx *juno.Tx, msg *govtypesv1.MsgVote) error {
	txTimestamp, err := time.Parse(time.RFC3339, tx.Timestamp)
	if err != nil {
		return fmt.Errorf("error while parsing time: %s", err)
	}

	vote := types.NewVote(msg.ProposalId, msg.Voter, msg.Option, "1.0", txTimestamp, tx.Height)

	err = m.db.SaveVote(vote)
	if err != nil {
		return fmt.Errorf("error while saving vote: %s", err)
	}

	// update tally result for given proposal
	return m.UpdateProposalTallyResult(msg.ProposalId, tx.Height)
}

// handleMsgVoteWeighted allows to properly handle a MsgVoteWeighted
func (m *Module) handleMsgVoteWeighted(tx *juno.Tx, msg *govtypesv1.MsgVoteWeighted) error {
	txTimestamp, err := time.Parse(time.RFC3339, tx.Timestamp)
	if err != nil {
		return fmt.Errorf("error while parsing time: %s", err)
	}

	for _, option := range msg.Options {
		vote := types.NewVote(msg.ProposalId, msg.Voter, option.Option, option.Weight, txTimestamp, tx.Height)
		err = m.db.SaveVote(vote)
		if err != nil {
			return fmt.Errorf("error while saving weighted vote for address %s: %s", msg.Voter, err)
		}
	}

	// update tally result for given proposal
	return m.UpdateProposalTallyResult(msg.ProposalId, tx.Height)
}
