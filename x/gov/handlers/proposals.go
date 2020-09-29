package handlers

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	gov "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/desmos-labs/juno/parse/client"
	juno "github.com/desmos-labs/juno/types"
	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/x/auth"
	ops "github.com/forbole/bdjuno/x/gov/operations"
	"github.com/forbole/bdjuno/x/gov/types"
)

// HandleMsgSubmitProposal allows to properly handle a HandleMsgSubmitProposal
func HandleMsgSubmitProposal(tx juno.Tx, msg gov.MsgSubmitProposal, db database.BigDipperDb, cp client.ClientProxy) error {
	timestamp, err := time.Parse(time.RFC3339, tx.Timestamp)
	if err != nil {
		return err
	}

	// Get proposal ID
	var restProposals gov.Proposals
	_, err = cp.QueryLCDWithHeight(fmt.Sprintf("/gov/proposals?height=%d", tx.Height), &restProposals)
	if err != nil {
		return err
	}

	var proposal gov.Proposal
	for _, p := range restProposals {
		if p.Content.GetTitle() == msg.Content.GetTitle() {
			proposal = p
			break
		}
	}

	if err = auth.RefreshAccounts([]sdk.AccAddress{msg.Proposer}, tx.Height, timestamp, cp, db); err != nil {
		return err
	}

	proposalObj := types.NewProposal(
		proposal.GetTitle(), proposal.GetDescription(), proposal.ProposalRoute(), proposal.ProposalType(),
		proposal.ProposalID, proposal.Status, proposal.SubmitTime, proposal.DepositEndTime,
		proposal.VotingStartTime, proposal.VotingEndTime, msg.Proposer,
	)
	if err = db.SaveProposal(proposalObj); err != nil {
		return err
	}

	deposit := types.NewDeposit(
		proposal.ProposalID, msg.Proposer, msg.InitialDeposit, msg.InitialDeposit, tx.Height, timestamp,
	)
	if err = db.SaveDeposit(deposit); err != nil {
		return err
	}

	// Watch the proposal and renew the database when deposit end and voting end in the future
	if proposal.Status.String() == "VotingPeriod" && proposal.VotingEndTime.After(time.Now()) {
		time.AfterFunc(time.Since(proposal.VotingEndTime), ops.UpdateProposal(proposal.ProposalID, cp, db))
	} else if proposal.Status.String() == "DepositPeriod" && proposal.DepositEndTime.After(time.Now()) {
		time.AfterFunc(time.Since(proposal.DepositEndTime), ops.UpdateProposal(proposal.ProposalID, cp, db))
	}

	return nil
}

// HandleMsgDeposit allows to properly handle a HandleMsgDeposit
//refresh the proposal and record the deposit
func HandleMsgDeposit(tx juno.Tx, msg gov.MsgDeposit, db database.BigDipperDb, cp client.ClientProxy) error {
	timestamp, err := time.Parse(time.RFC3339, tx.Timestamp)
	if err != nil {
		return err
	}

	if err = auth.RefreshAccounts([]sdk.AccAddress{msg.Depositor}, tx.Height, timestamp, cp, db); err != nil {
		return err
	}

	//getTotalDeposit
	var s gov.Proposals
	_, err = cp.QueryLCDWithHeight(fmt.Sprintf("/gov/proposals?height=%d/%d", tx.Height, msg.ProposalID), &s)
	if err != nil {
		return err
	}
	for _, proposal := range s {
		if err = db.SaveDeposit(types.NewDeposit(msg.ProposalID, msg.Depositor, msg.Amount, proposal.TotalDeposit,
			tx.Height, timestamp)); err != nil {
			return err
		}
	}
	return nil
}

// HandleMsgVote allows to properly handle a HandleMsgVote
func HandleMsgVote(tx juno.Tx, msg gov.MsgVote, db database.BigDipperDb, cp client.ClientProxy) error {
	//each vote voted
	timestamp, err := time.Parse(time.RFC3339, tx.Timestamp)
	if err != nil {
		return err
	}

	if err = auth.RefreshAccounts([]sdk.AccAddress{msg.Voter}, tx.Height, timestamp, cp, db); err != nil {
		return err
	}

	//fetch from lcd & store voter in specific time
	var s gov.TallyResult
	_, err = cp.QueryLCDWithHeight(fmt.Sprintf("/gov/proposals?height=%d/%d/tally", tx.Height, msg.ProposalID), &s)
	if err != nil {
		return err
	}

	if err = db.SaveVote(types.NewVote(msg.ProposalID, msg.Voter, msg.Option, tx.Height, timestamp)); err != nil {
		return err
	}

	return db.SaveTallyResult(types.NewTallyResult(msg.ProposalID, s.Yes.Int64(), s.Abstain.Int64(), s.No.Int64(), s.NoWithVeto.Int64(),
		tx.Height, timestamp))
}
