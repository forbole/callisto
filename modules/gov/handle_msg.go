package gov

import (
	"context"
	"strconv"

	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	juno "github.com/desmos-labs/juno/types"
)

// HandleMsg allows to handle the different utils related to the staking module
func HandleMsg(
	tx *juno.Tx, index int, msg sdk.Msg,
	govClient govtypes.QueryClient,
	cdc codec.Marshaler, db *database.Db,
) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	switch cosmosMsg := msg.(type) {
	case *govtypes.MsgSubmitProposal:
		return handleMsgSubmitProposal(tx, index, cosmosMsg, govClient, cdc, db)

	case *govtypes.MsgDeposit:
		return handleMsgDeposit(tx, cosmosMsg, db)

	case *govtypes.MsgVote:
		return handleMsgVote(tx, cosmosMsg, db)
	}

	return nil
}

// handleMsgSubmitProposal allows to properly handle a handleMsgSubmitProposal
func handleMsgSubmitProposal(
	tx *juno.Tx, index int, msg *govtypes.MsgSubmitProposal,
	govClient govtypes.QueryClient, cdc codec.Marshaler, db *database.Db,
) error {
	// Get the proposal id
	event, err := tx.FindEventByType(index, govtypes.EventTypeSubmitProposal)
	if err != nil {
		return err
	}

	id, err := tx.FindAttributeByKey(event, govtypes.AttributeKeyProposalID)
	if err != nil {
		return err
	}

	proposalID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return err
	}

	// Get the proposal
	res, err := govClient.Proposal(
		context.Background(),
		&govtypes.QueryProposalRequest{ProposalId: proposalID},
	)
	if err != nil {
		return err
	}

	proposal := res.Proposal

	// Unpack the content
	var content govtypes.Content
	err = cdc.UnpackAny(proposal.Content, &content)
	if err != nil {
		return err
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
	err = db.SaveProposals([]types.Proposal{proposalObj})
	if err != nil {
		return err
	}

	// Store the deposit
	deposit := types.NewDeposit(proposal.ProposalId, msg.Proposer, msg.InitialDeposit, tx.Height)
	return db.SaveDeposits([]types.Deposit{deposit})
}

// handleMsgDeposit allows to properly handle a handleMsgDeposit
func handleMsgDeposit(tx *juno.Tx, msg *govtypes.MsgDeposit, db *database.Db) error {
	deposit := types.NewDeposit(msg.ProposalId, msg.Depositor, msg.Amount, tx.Height)
	return db.SaveDeposits([]types.Deposit{deposit})
}

// handleMsgVote allows to properly handle a handleMsgVote
func handleMsgVote(tx *juno.Tx, msg *govtypes.MsgVote, db *database.Db) error {
	vote := types.NewVote(msg.ProposalId, msg.Voter, msg.Option, tx.Height)
	return db.SaveVote(vote)
}
