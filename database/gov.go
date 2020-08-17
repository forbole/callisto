package database

import (
	"fmt"

	dbtypes "github.com/forbole/bdjuno/database/types"
	types "github.com/forbole/bdjuno/x/gov/types"
	"github.com/lib/pq"
)

// SaveProposals allows to save for the given height the given total amount of coins
func (db BigDipperDb) SaveProposals(proposals []types.Proposal) error {
	query := `INSERT INTO proposal(title,description ,proposer,proposal_route ,proposal_type,proposal_id,
		status,submit_time ,deposit_end_time ,voting_start_time,voting_end_time) VALUES`
	var param []interface{}
	for i, proposal := range proposals {
		vi := i * 11
		query += fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d),",
			vi+1, vi+2, vi+3, vi+4, vi+5, vi+6, vi+7, vi+8, vi+9, vi+10, vi+11)
		param = append(param, proposal.Title,
			proposal.Description,
			proposal.Proposer.String(),
			proposal.ProposalRoute,
			proposal.ProposalType,
			proposal.ProposalID,
			proposal.Status.String(),
			proposal.SubmitTime,
			proposal.DepositEndTime,
			proposal.VotingStartTime,
			proposal.VotingEndTime)
	}
	query = query[:len(query)-1] // Remove trailing ","
	query += " ON CONFLICT DO UPDATE"
	_, err := db.Sql.Exec(query, param...)
	if err != nil {
		return err
	}
	return nil
}

//SaveProposal save a single proposal
func (db BigDipperDb) SaveProposal(proposal types.Proposal) error {
	query := `INSERT INTO proposal(title,description ,proposer,proposal_route ,proposal_type,proposal_id,
		status,submit_time ,deposit_end_time,voting_start_time,voting_end_time)
		 VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) ON CONFLICT DO NOTHING`

	_, err := db.Sql.Exec(query, proposal.Title,
		proposal.Description,
		proposal.Proposer.String(),
		proposal.ProposalRoute,
		proposal.ProposalType,
		proposal.ProposalID,
		proposal.Status.String(),
		proposal.SubmitTime,
		proposal.DepositEndTime,
		proposal.VotingStartTime,
		proposal.VotingEndTime)
	if err != nil {
		return err
	}
	return nil
}

// SaveTallyResults allows to save for the given height the given total amount of coins
func (db BigDipperDb) SaveTallyResults(tallys []types.TallyResult) error {
	query := `INSERT INTO tally_result(proposal_id,yes,abstain,no,no_with_veto,height,timestamp) VALUES`
	var param []interface{}
	for i, tally := range tallys {
		vi := i * 7
		query += fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d,$%d),", vi+1, vi+2, vi+3, vi+4, vi+5, vi+6, vi+7)
		param = append(param, tally.ProposalID,
			tally.Yes,
			tally.Abstain,
			tally.No,
			tally.NoWithVeto,
			tally.Height,
			tally.Timestamp)
	}
	query = query[:len(query)-1] // Remove trailing ","
	query += " ON CONFLICT DO UPDATE"
	_, err := db.Sql.Exec(query, param...)
	if err != nil {
		return err
	}
	return nil
}

// SaveTallyResult insert a single row into tally_result table
func (db BigDipperDb) SaveTallyResult(tally types.TallyResult) error {
	query := `INSERT INTO tally_result(proposal_id,yes,abstain,no,no_with_veto,height,timestamp) VALUES
	($1,$2,$3,$4,$5,$6,$7) ON CONFLICT DO NOTHING`
	_, err := db.Sql.Exec(query, tally.ProposalID,
		tally.Yes,
		tally.Abstain,
		tally.No,
		tally.NoWithVeto,
		tally.Height,
		tally.Timestamp)
	if err != nil {
		return err
	}
	return nil
}

// SaveVote allows to save for the given height and the message vote
func (db BigDipperDb) SaveVote(vote types.Vote) error {
	query := `INSERT INTO vote(proposal_id,voter,option,height,timestamp) VALUES
	($1,$2,$3,$4,$5) ON CONFLICT DO NOTHING`
	_, err := db.Sql.Exec(query,
		vote.ProposalID,
		vote.Voter.String(),
		vote.Option.String(),
		vote.Height,
		vote.Timestamp)
	if err != nil {
		return err
	}
	return nil
}

// SaveDeposit allows to save for the given message deposit and height
func (db BigDipperDb) SaveDeposit(deposit types.Deposit) error {
	query := `INSERT INTO deposit(proposal_id,depositor,amount,total_deposit,height,timestamp) VALUES
	($1,$2,$3,$4,$5,$6) ON CONFLICT DO NOTHING`
	_, err := db.Sql.Exec(query,
		deposit.ProposalID,
		deposit.Depositor.String(),
		pq.Array(dbtypes.NewDbCoins(deposit.Amount)),
		pq.Array(dbtypes.NewDbCoins(deposit.TotalDeposit)),
		deposit.Height,
		deposit.Timestamp)
	if err != nil {
		return err
	}
	return nil
}

// SaveDeposits allows to save multiple deposits
func (db BigDipperDb) SaveDeposits(deposits []types.Deposit) error {
	query := `INSERT INTO deposit(proposal_id,depositor,amount,total_deposit,height,timestamp) VALUES `
	var param []interface{}

	for i, deposit := range deposits {
		vi := i * 6
		query += fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d),", vi+1, vi+2, vi+3, vi+4, vi+5, vi+6)
		param = append(param, deposit.ProposalID,
			deposit.Depositor.String(),
			pq.Array(dbtypes.NewDbCoins(deposit.Amount)),
			pq.Array(dbtypes.NewDbCoins(deposit.TotalDeposit)),
			deposit.Height,
			deposit.Timestamp)
	}
	query = query[:len(query)-1] // Remove trailing ","
	query += " ON CONFLICT DO NOTHING"
	_, err := db.Sql.Exec(query, param...)
	if err != nil {
		return err
	}
	return nil
}

func (db BigDipperDb) UpdateProposal(proposal types.Proposal) error {
	query := `UPDATE proposal SET (status,voting_start_time,voting_end_time)
		 = ($1,$2,$3) where proposal_id=$4`

	_, err := db.Sql.Exec(query,
		proposal.Status.String(),
		proposal.VotingStartTime,
		proposal.VotingEndTime,
		proposal.ProposalID)
	if err != nil {
		return err
	}
	return nil
}
