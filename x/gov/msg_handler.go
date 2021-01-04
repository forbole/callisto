package gov

import (
	"fmt"
	"time"

	"github.com/desmos-labs/juno/client"

	"github.com/forbole/bdjuno/x/auth"
	govtypes "github.com/forbole/bdjuno/x/gov/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	juno "github.com/desmos-labs/juno/types"

	"github.com/forbole/bdjuno/database"
)

// HandleMsg allows to handle the different messages related to the staking module
func HandleMsg(tx *juno.Tx, msg sdk.Msg, cp *client.Proxy, db *database.BigDipperDb) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	switch cosmosMsg := msg.(type) {
	case gov.MsgSubmitProposal:
		return handleMsgSubmitProposal(tx, cosmosMsg, db, cp)

	case gov.MsgDeposit:
		return handleMsgDeposit(tx, cosmosMsg, db, cp)

	case gov.MsgVote:
		return handleMsgVote(tx, cosmosMsg, db, cp)
	}

	return nil
}

// handleMsgSubmitProposal allows to properly handle a handleMsgSubmitProposal
func handleMsgSubmitProposal(tx *juno.Tx, msg gov.MsgSubmitProposal, db *database.BigDipperDb, cp *client.Proxy) error {
	// Get proposals
	var restProposals gov.Proposals
	_, err := cp.QueryLCDWithHeight(fmt.Sprintf("/gov/proposals?height=%d", tx.Height), &restProposals)
	if err != nil {
		return err
	}

	// Get the specific proposal
	var proposal gov.Proposal
	for _, p := range restProposals {
		if p.Content.GetTitle() == msg.Content.GetTitle() {
			proposal = p
			break
		}
	}

	// Refresh the accounts
	err = auth.RefreshAccounts([]string{msg.Proposer.String()}, tx.Height, cp, db)
	if err != nil {
		return err
	}

	// Store the proposal
	proposalObj := govtypes.NewProposal(
		proposal.GetTitle(), proposal.GetDescription(), proposal.ProposalRoute(), proposal.ProposalType(),
		proposal.ProposalID, proposal.Status, proposal.SubmitTime, proposal.DepositEndTime,
		proposal.VotingStartTime, proposal.VotingEndTime, msg.Proposer.String(),
	)
	err = db.SaveProposal(proposalObj)
	if err != nil {
		return err
	}

	// Store the deposit
	deposit := govtypes.NewDeposit(
		proposal.ProposalID, msg.Proposer.String(), msg.InitialDeposit, msg.InitialDeposit, tx.Height,
	)
	err = db.SaveDeposit(deposit)
	if err != nil {
		return err
	}

	// Watch the proposal and renew the database when deposit end and voting end in the future
	if proposal.Status.String() == "VotingPeriod" && proposal.VotingEndTime.After(time.Now()) {
		time.AfterFunc(time.Since(proposal.VotingEndTime), UpdateProposal(proposal.ProposalID, cp, db))
	} else if proposal.Status.String() == "DepositPeriod" && proposal.DepositEndTime.After(time.Now()) {
		time.AfterFunc(time.Since(proposal.DepositEndTime), UpdateProposal(proposal.ProposalID, cp, db))
	}

	return nil
}

// handleMsgDeposit allows to properly handle a handleMsgDeposit
func handleMsgDeposit(tx *juno.Tx, msg gov.MsgDeposit, db *database.BigDipperDb, cp *client.Proxy) error {
	// Refresh the accounts
	err := auth.RefreshAccounts([]string{msg.Depositor.String()}, tx.Height, cp, db)
	if err != nil {
		return err
	}

	// Get proposals
	var s gov.Proposals
	_, err = cp.QueryLCDWithHeight(fmt.Sprintf("/gov/proposals?height=%d/%d", tx.Height, msg.ProposalID), &s)
	if err != nil {
		return err
	}

	// Save the deposits
	for _, proposal := range s {
		deposit := govtypes.NewDeposit(
			msg.ProposalID, msg.Depositor.String(), msg.Amount, proposal.TotalDeposit, tx.Height,
		)
		if err = db.SaveDeposit(deposit); err != nil {
			return err
		}
	}

	return nil
}

// handleMsgVote allows to properly handle a handleMsgVote
func handleMsgVote(tx *juno.Tx, msg gov.MsgVote, db *database.BigDipperDb, cp *client.Proxy) error {
	// Refresh accounts
	err := auth.RefreshAccounts([]string{msg.Voter.String()}, tx.Height, cp, db)
	if err != nil {
		return err
	}

	// Get the rally result
	var s gov.TallyResult
	_, err = cp.QueryLCDWithHeight(fmt.Sprintf("/gov/proposals?height=%d/%d/tally", tx.Height, msg.ProposalID), &s)
	if err != nil {
		return err
	}

	vote := govtypes.NewVote(msg.ProposalID, msg.Voter.String(), msg.Option, tx.Height)
	err = db.SaveVote(vote)
	if err != nil {
		return err
	}

	// Save tally result
	tallyResult := govtypes.NewTallyResult(
		msg.ProposalID, s.Yes.Int64(), s.Abstain.Int64(), s.No.Int64(), s.NoWithVeto.Int64(), tx.Height,
	)
	return db.SaveTallyResult(tallyResult)
}
