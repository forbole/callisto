package gov

import (
	"context"
	"time"

	bigdipperdb "github.com/forbole/bdjuno/database/bigdipper"
	"github.com/forbole/bdjuno/modules/bigdipper/gov/types"

	"github.com/cosmos/cosmos-sdk/codec"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/cosmos/cosmos-sdk/types/query"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	juno "github.com/desmos-labs/juno/types"
)

// HandleMsg allows to handle the different utils related to the staking module
func HandleMsg(
	tx *juno.Tx, msg sdk.Msg,
	govClient govtypes.QueryClient, authClient authtypes.QueryClient, bankClient banktypes.QueryClient,
	cdc codec.Marshaler, db *bigdipperdb.Db,
) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	switch cosmosMsg := msg.(type) {
	case *govtypes.MsgSubmitProposal:
		return handleMsgSubmitProposal(tx, cosmosMsg, govClient, authClient, bankClient, cdc, db)

	case *govtypes.MsgDeposit:
		return handleMsgDeposit(tx, cosmosMsg, db)

	case *govtypes.MsgVote:
		return handleMsgVote(tx, cosmosMsg, db)
	}

	return nil
}

// handleMsgSubmitProposal allows to properly handle a handleMsgSubmitProposal
func handleMsgSubmitProposal(
	tx *juno.Tx, msg *govtypes.MsgSubmitProposal,
	govClient govtypes.QueryClient, authClient authtypes.QueryClient, bankClient banktypes.QueryClient,
	cdc codec.Marshaler, db *bigdipperdb.Db,
) error {
	// Get proposals
	res, err := govClient.Proposals(
		context.Background(),
		&govtypes.QueryProposalsRequest{
			Pagination: &query.PageRequest{
				Limit: 10000, // Query 10.000 proposals
			},
		},
	)
	if err != nil {
		return err
	}

	// Get the specific proposal
	var proposal govtypes.Proposal
	for _, p := range res.Proposals {
		// Unmarshal the content properly
		var content govtypes.Content
		err = cdc.UnpackAny(p.Content, &content)
		if err != nil {
			return err
		}

		if p.GetContent().GetTitle() == msg.GetContent().GetTitle() {
			proposal = p
			break
		}
	}

	// Store the proposal
	proposalObj := types.NewProposal(
		proposal.GetTitle(),
		proposal.GetContent().GetDescription(),
		proposal.ProposalRoute(),
		proposal.ProposalType(),
		proposal.ProposalId,
		proposal.Status,
		proposal.SubmitTime,
		proposal.DepositEndTime,
		proposal.VotingStartTime,
		proposal.VotingEndTime,
		msg.Proposer,
	)
	err = db.SaveProposal(proposalObj)
	if err != nil {
		return err
	}

	// Store the deposit
	deposit := types.NewDeposit(
		proposal.ProposalId,
		msg.Proposer,
		msg.InitialDeposit,
		tx.Height,
	)
	err = db.SaveDeposit(deposit)
	if err != nil {
		return err
	}

	// Watch the proposal and renew the BigDipper when deposit end and voting end in the future
	update := UpdateProposal(proposal.ProposalId, govClient, authClient, bankClient, cdc, db)
	if proposal.Status == govtypes.StatusVotingPeriod && proposal.VotingEndTime.After(time.Now()) {
		time.AfterFunc(time.Until(proposal.VotingEndTime), update)
	} else if proposal.Status == govtypes.StatusDepositPeriod && proposal.DepositEndTime.After(time.Now()) {
		time.AfterFunc(time.Until(proposal.DepositEndTime), update)
	}

	return nil
}

// handleMsgDeposit allows to properly handle a handleMsgDeposit
func handleMsgDeposit(
	tx *juno.Tx, msg *govtypes.MsgDeposit, db *bigdipperdb.Db,
) error {
	// Save the deposits
	deposit := types.NewDeposit(msg.ProposalId, msg.Depositor, msg.Amount, tx.Height)
	return db.SaveDeposit(deposit)
}

// handleMsgVote allows to properly handle a handleMsgVote
func handleMsgVote(
	tx *juno.Tx, msg *govtypes.MsgVote, db *bigdipperdb.Db,
) error {
	// Save the vote
	vote := types.NewVote(msg.ProposalId, msg.Voter, msg.Option, tx.Height)
	return db.SaveVote(vote)
}
