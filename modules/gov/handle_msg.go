package gov

import (
	"fmt"
	"strings"
	"time"

	"github.com/forbole/bdjuno/v4/types"
	"google.golang.org/grpc/codes"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/authz"

	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govtypesv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"

	juno "github.com/forbole/juno/v4/types"
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
	case *govtypesv1.MsgSubmitProposal:
		return m.handleSubmitProposalEvent(tx, cosmosMsg.Proposer, tx.Logs[index].Events)
	case *govtypesv1beta1.MsgSubmitProposal:
		return m.handleSubmitProposalEvent(tx, cosmosMsg.Proposer, tx.Logs[index].Events)

	case *govtypesv1.MsgDeposit:
		return m.handleDepositEvent(tx, cosmosMsg.Depositor, tx.Logs[index].Events)
	case *govtypesv1beta1.MsgDeposit:
		return m.handleDepositEvent(tx, cosmosMsg.Depositor, tx.Logs[index].Events)

	case *govtypesv1.MsgVote:
		return m.handleVoteEvent(tx, cosmosMsg.Voter, tx.Logs[index].Events)
	case *govtypesv1beta1.MsgVote:
		return m.handleVoteEvent(tx, cosmosMsg.Voter, tx.Logs[index].Events)

	case *govtypesv1.MsgVoteWeighted:
		return m.handleVoteEvent(tx, cosmosMsg.Voter, tx.Logs[index].Events)
	case *govtypesv1beta1.MsgVoteWeighted:
		return m.handleVoteEvent(tx, cosmosMsg.Voter, tx.Logs[index].Events)
	}

	return nil
}

// handleSubmitProposalEvent allows to properly handle a handleSubmitProposalEvent
func (m *Module) handleSubmitProposalEvent(tx *juno.Tx, proposer string, events sdk.StringEvents) error {
	// Get the proposal id
	proposalID, err := ProposalIDFromEvents(events)
	if err != nil {
		return fmt.Errorf("error while getting proposal id: %s", err)
	}

	// Get the proposal
	proposal, err := m.source.Proposal(tx.Height, proposalID)
	if err != nil {
		if strings.Contains(err.Error(), codes.NotFound.String()) {
			// query the proposal details using the latest height stored in db
			// to fix the rpc error returning code = NotFound desc = proposal x doesn't exist
			block, err := m.db.GetLastBlockHeightAndTimestamp()
			if err != nil {
				return fmt.Errorf("error while getting latest block height: %s", err)
			}
			proposal, err = m.source.Proposal(block.Height, proposalID)
			if err != nil {
				return fmt.Errorf("error while getting proposal: %s", err)
			}
		} else {
			return fmt.Errorf("error while getting proposal: %s", err)
		}
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

	err = m.db.SaveAccounts(addresses)
	if err != nil {
		return fmt.Errorf("error while storing proposal recipient: %s", err)
	}

	// Unpack the proposal interfaces
	err = proposal.UnpackInterfaces(m.cdc)
	if err != nil {
		return fmt.Errorf("error while unpacking proposal interfaces: %s", err)
	}

	// Store the proposal
	proposalObj := types.NewProposal(
		proposal.ProposalId,
		proposal.GetContent().ProposalRoute(),
		proposal.GetContent().ProposalType(),
		proposal.GetContent(),
		proposal.Status.String(),
		*proposal.SubmitTime,
		*proposal.DepositEndTime,
		proposal.VotingStartTime,
		proposal.VotingEndTime,
		proposer,
	)

	err = m.db.SaveProposals([]types.Proposal{proposalObj})
	if err != nil {
		return fmt.Errorf("error while saving proposal: %s", err)
	}

	// Submit proposal must have a deposit event with depositor equal to the proposer
	return m.handleDepositEvent(tx, proposer, events)
}

// handleDepositEvent allows to properly handle a handleDepositEvent
func (m *Module) handleDepositEvent(tx *juno.Tx, depositor string, events sdk.StringEvents) error {
	// Get the proposal id
	proposalID, err := ProposalIDFromEvents(events)
	if err != nil {
		return fmt.Errorf("error while getting proposal id: %s", err)
	}

	deposit, err := m.source.ProposalDeposit(tx.Height, proposalID, depositor)
	if err != nil {
		return fmt.Errorf("error while getting proposal deposit: %s", err)
	}
	txTimestamp, err := time.Parse(time.RFC3339, tx.Timestamp)
	if err != nil {
		return fmt.Errorf("error while parsing time: %s", err)
	}

	return m.db.SaveDeposits([]types.Deposit{
		types.NewDeposit(proposalID, depositor, deposit.Amount, txTimestamp, tx.Height),
	})
}

// handleVoteEvent allows to properly handle a handleVoteEvent
func (m *Module) handleVoteEvent(tx *juno.Tx, voter string, events sdk.StringEvents) error {
	// Get the proposal id
	proposalID, err := ProposalIDFromEvents(events)
	if err != nil {
		return fmt.Errorf("error while getting proposal id: %s", err)
	}

	txTimestamp, err := time.Parse(time.RFC3339, tx.Timestamp)
	if err != nil {
		return fmt.Errorf("error while parsing time: %s", err)
	}

	// Get the vote option
	voteOption, err := VoteOptionFromEvents(events)
	if err != nil {
		return fmt.Errorf("error while getting vote option: %s", err)
	}

	vote := types.NewVote(proposalID, voter, voteOption, txTimestamp, tx.Height)

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
