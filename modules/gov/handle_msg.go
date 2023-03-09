package gov

import (
	"fmt"
	"time"

	"strconv"

	"github.com/forbole/bdjuno/v4/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	ibcclienttypes "github.com/cosmos/ibc-go/v4/modules/core/02-client/types"
	ccvprovidertypes "github.com/cosmos/interchain-security/x/ccv/provider/types"
	juno "github.com/forbole/juno/v4/types"
)

// HandleMsg implements modules.MessageModule
func (m *Module) HandleMsg(index int, msg sdk.Msg, tx *juno.Tx) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	switch cosmosMsg := msg.(type) {
	case *govtypes.MsgSubmitProposal:
		switch c := cosmosMsg.GetContent().(type) {
		case *ccvprovidertypes.ConsumerAdditionProposal:
			return m.handleConsumerAdditionProposal(c, tx, 0, cosmosMsg)
		case *ccvprovidertypes.ConsumerRemovalProposal:
			return m.handleConsumerRemovalProposal(c, tx, 0, cosmosMsg)
		default:
			return m.handleMsgSubmitProposal(tx, index, cosmosMsg)
		}
	case *govtypes.MsgDeposit:
		return m.handleMsgDeposit(tx, cosmosMsg)

	case *govtypes.MsgVote:
		return m.handleMsgVote(tx, cosmosMsg)
	}

	return nil
}

// handleMsgSubmitProposal allows to properly handle MsgSubmitProposal
func (m *Module) handleMsgSubmitProposal(tx *juno.Tx, index int, msg *govtypes.MsgSubmitProposal) error {
	// Get the proposal id
	event, err := tx.FindEventByType(index, govtypes.EventTypeSubmitProposal)
	if err != nil {
		return fmt.Errorf("error while searching for EventTypeSubmitProposal: %s", err)
	}

	id, err := tx.FindAttributeByKey(event, govtypes.AttributeKeyProposalID)
	if err != nil {
		return fmt.Errorf("error while searching for AttributeKeyProposalID: %s", err)
	}

	proposalID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return fmt.Errorf("error while parsing proposal id: %s", err)
	}

	// Get the proposal
	proposal, err := m.source.Proposal(tx.Height, proposalID)
	if err != nil {
		return fmt.Errorf("error while getting proposal: %s", err)
	}

	// Unpack the content
	var content govtypes.Content
	err = m.cdc.UnpackAny(proposal.Content, &content)
	if err != nil {
		return fmt.Errorf("error while unpacking proposal content: %s", err)
	}

	// Store the proposal
	proposalObj := types.NewProposal(
		proposal.ProposalId,
		proposal.ProposalRoute(),
		proposal.ProposalType(),
		proposal.GetContent(),
		proposal.Status.String(),
		proposal.SubmitTime,
		proposal.DepositEndTime,
		proposal.VotingStartTime,
		proposal.VotingEndTime,
		msg.Proposer,
	)
	err = m.db.SaveProposals([]types.Proposal{proposalObj})
	if err != nil {
		return err
	}

	// Store the deposit
	deposit := types.NewDeposit(proposal.ProposalId, msg.Proposer, msg.InitialDeposit, tx.Height)
	return m.db.SaveDeposits([]types.Deposit{deposit})
}

// handleMsgDeposit allows to properly handle a MsgDeposit
func (m *Module) handleMsgDeposit(tx *juno.Tx, msg *govtypes.MsgDeposit) error {
	deposit, err := m.source.ProposalDeposit(tx.Height, msg.ProposalId, msg.Depositor)
	if err != nil {
		return fmt.Errorf("error while getting proposal deposit: %s", err)
	}

	return m.db.SaveDeposits([]types.Deposit{
		types.NewDeposit(msg.ProposalId, msg.Depositor, deposit.Amount, tx.Height),
	})
}

// handleMsgVote allows to properly handle a MsgVote
func (m *Module) handleMsgVote(tx *juno.Tx, msg *govtypes.MsgVote) error {
	vote := types.NewVote(msg.ProposalId, msg.Voter, msg.Option, tx.Height)
	return m.db.SaveVote(vote)
}

// handleConsumerAdditionProposal allows to properly handle a ConsumerAdditionProposal
func (m *Module) handleConsumerAdditionProposal(consAddProposal *ccvprovidertypes.ConsumerAdditionProposal, tx *juno.Tx, index int, msg *govtypes.MsgSubmitProposal) error {
	// Get the proposal id
	event, err := tx.FindEventByType(index, govtypes.EventTypeSubmitProposal)
	if err != nil {
		return fmt.Errorf("error while searching for EventTypeSubmitProposal: %s", err)
	}

	id, err := tx.FindAttributeByKey(event, govtypes.AttributeKeyProposalID)
	if err != nil {
		return fmt.Errorf("error while searching for AttributeKeyProposalID: %s", err)
	}

	proposalID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return fmt.Errorf("error while parsing proposal id: %s", err)
	}

	// Get the proposal
	proposal, err := m.source.Proposal(tx.Height, proposalID)
	if err != nil {
		return fmt.Errorf("error while getting proposal: %s", err)
	}

	// Store the ccv ConsumerAdditionProposal proposal
	ccvProposalContent := types.NewCcvProposalContent(consAddProposal.Title, consAddProposal.Description, consAddProposal.ChainId, string(consAddProposal.GenesisHash), string(consAddProposal.BinaryHash),
		consAddProposal.SpawnTime, time.Time{}, consAddProposal.InitialHeight, consAddProposal.UnbondingPeriod, consAddProposal.CcvTimeoutPeriod,
		consAddProposal.TransferTimeoutPeriod, consAddProposal.ConsumerRedistributionFraction, consAddProposal.BlocksPerDistributionTransmission,
		consAddProposal.HistoricalEntries)

	ccvProposalObj := types.NewCcvProposal(
		proposal.ProposalId,
		consAddProposal.ProposalRoute(),
		consAddProposal.ProposalType(),
		ccvProposalContent,
		proposal.Status.String(),
		proposal.SubmitTime,
		proposal.DepositEndTime,
		proposal.VotingStartTime,
		proposal.VotingEndTime,
		msg.Proposer,
	)

	err = m.db.SaveCcvProposals(tx.Height, []types.CcvProposal{ccvProposalObj})
	if err != nil {
		return err
	}

	// Store the deposit
	deposit := types.NewDeposit(proposal.ProposalId, msg.Proposer, msg.InitialDeposit, tx.Height)
	return m.db.SaveDeposits([]types.Deposit{deposit})
}

// handleConsumerRemovalProposal allows to properly handle a ConsumerRemovalProposal
func (m *Module) handleConsumerRemovalProposal(consRemProposal *ccvprovidertypes.ConsumerRemovalProposal, tx *juno.Tx, index int, msg *govtypes.MsgSubmitProposal) error {
	// Get the proposal id
	event, err := tx.FindEventByType(index, govtypes.EventTypeSubmitProposal)
	if err != nil {
		return fmt.Errorf("error while searching for EventTypeSubmitProposal: %s", err)
	}

	id, err := tx.FindAttributeByKey(event, govtypes.AttributeKeyProposalID)
	if err != nil {
		return fmt.Errorf("error while searching for AttributeKeyProposalID: %s", err)
	}

	proposalID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return fmt.Errorf("error while parsing proposal id: %s", err)
	}

	// Get the proposal
	proposal, err := m.source.Proposal(tx.Height, proposalID)
	if err != nil {
		return fmt.Errorf("error while getting proposal: %s", err)
	}

	// Store the ccv ConsumerRemovalProposal proposal
	ccvProposalContent := types.NewCcvProposalContent(consRemProposal.Title, consRemProposal.Description,
		consRemProposal.ChainId, "", "", time.Time{}, consRemProposal.StopTime, ibcclienttypes.ZeroHeight(),
		0, 0, 0, "", 0, 0)

	ccvProposalObj := types.NewCcvProposal(
		proposal.ProposalId,
		consRemProposal.ProposalRoute(),
		consRemProposal.ProposalType(),
		ccvProposalContent,
		proposal.Status.String(),
		proposal.SubmitTime,
		proposal.DepositEndTime,
		proposal.VotingStartTime,
		proposal.VotingEndTime,
		msg.Proposer,
	)

	err = m.db.SaveCcvProposals(tx.Height, []types.CcvProposal{ccvProposalObj})
	if err != nil {
		return err
	}

	// Store the deposit
	deposit := types.NewDeposit(proposal.ProposalId, msg.Proposer, msg.InitialDeposit, tx.Height)
	return m.db.SaveDeposits([]types.Deposit{deposit})
}
