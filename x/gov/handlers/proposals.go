package handlers

import (
	"fmt"
	"time"

	gov "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/desmos-labs/juno/parse/client"
	juno "github.com/desmos-labs/juno/types"
	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/x/gov/types"
)

// HandleMsgSubmitProposal allows to properly handle a HandleMsgSubmitProposal
func HandleMsgSubmitProposal(tx juno.Tx, msg gov.MsgSubmitProposal, db database.BigDipperDb, cp client.ClientProxy) error {
	timestamp, err := time.Parse(time.RFC3339, tx.Timestamp)
	if err != nil {
		return err
	}

	//get proposal ID
	var s gov.Proposals
	_, err = cp.QueryLCDWithHeight("/gov/proposals/", &s)
	if err != nil {
		return err
	}
	var proposal gov.Proposal
	for _, p := range s {
		if p.Content.GetTitle() == msg.Content.GetTitle() {
			proposal = p
			break
		}
	}

	submitTime, err := time.Parse(time.RFC3339, proposal.SubmitTime.String())
	if err != nil {
		return err
	}
	depositEndTime, err := time.Parse(time.RFC3339, proposal.DepositEndTime.String())
	if err != nil {
		return err
	}
    votingStartTime, err := time.Parse(time.RFC3339, proposal.VotingStartTime.String())
	if err != nil {
		return err
	}
	votingEndTime, err := time.Parse(time.RFC3339, proposal.VotingEndTime.String())
	if err != nil {
		return err
	} 

	db.SaveProposal(types.NewProposal(proposal.GetTitle(),proposal.GetDescription(),proposal.ProposalRoute(),proposal.ProposalType(),proposal.ProposalID,proposal.Status.String(),
	submitTime,depositEndTime,proposal.TotalDeposit,votingStartTime,votingEndTime,msg.Proposer))

	db.SaveDeposit(types.NewDeposit(proposal.ProposalID, msg.Proposer, msg.InitialDeposit, tx.Height, timestamp))
	return nil
}

// HandleMsgDeposit allows to properly handle a HandleMsgDeposit
//refresh the proposal and record the deposit
func HandleMsgDeposit(tx juno.Tx, msg gov.MsgDeposit, db database.BigDipperDb, cp client.ClientProxy) error {
	timestamp, err := time.Parse(time.RFC3339, tx.Timestamp)
	if err != nil {
		return err
	}

	return db.SaveDeposit(types.NewDeposit(msg.ProposalID, msg.Depositor, msg.Amount, tx.Height, timestamp))
}

// HandleMsgVote allows to properly handle a HandleMsgVote
func HandleMsgVote(tx juno.Tx, msg gov.MsgVote, db database.BigDipperDb, cp client.ClientProxy) error {
	//each vote voted
	timestamp, err := time.Parse(time.RFC3339, tx.Timestamp)
	if err != nil {
		return err
	}

	//fetch from lcd & store voter in specific time
	var s gov.TallyResult
	_, err = cp.QueryLCDWithHeight(fmt.Sprintf("/gov/proposals/%s/tally", string(msg.ProposalID)), &s)
	if err != nil {
		return err
	}

	if err = db.SaveVote(types.NewVote(msg.ProposalID, msg.Voter, msg.Option, tx.Height, timestamp)); err != nil {
		return err
	}

	return db.SaveTallyResult(types.NewTallyResult(msg.ProposalID, s.Yes.Int64(), s.Abstain.Int64(), s.No.Int64(), s.NoWithVeto.Int64(),
		tx.Height, timestamp))
}
