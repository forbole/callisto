package database

import (
	"fmt"

	"github.com/lib/pq"

	dbtypes "github.com/forbole/bdjuno/database/types"
	"github.com/forbole/bdjuno/x/gov/types"
)

// SaveProposals allows to save for the given height the given total amount of coins
func (db *BigDipperDb) SaveProposals(proposals []types.Proposal) error {
	//do nothing if empty
	if len(proposals) == 0 {
		return nil
	}

	query := `
INSERT INTO proposal(
	title, description, proposer, proposal_route, proposal_type, proposal_id, status, 
    submit_time, deposit_end_time, voting_start_time, voting_end_time
) VALUES`
	var param []interface{}
	for i, proposal := range proposals {
		vi := i * 11
		query += fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d),",
			vi+1, vi+2, vi+3, vi+4, vi+5, vi+6, vi+7, vi+8, vi+9, vi+10, vi+11)
		param = append(param, proposal.Title,
			proposal.Description,
			proposal.Proposer,
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
	query += " ON CONFLICT DO NOTHING"
	_, err := db.Sql.Exec(query, param...)
	return err
}

//SaveProposal save a single proposal
func (db *BigDipperDb) SaveProposal(proposal types.Proposal) error {
	query := `
INSERT INTO proposal(
	title, description, proposer, proposal_route, proposal_type, proposal_id, status, 
    submit_time, deposit_end_time, voting_start_time, voting_end_time
) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) ON CONFLICT DO NOTHING`

	_, err := db.Sql.Exec(query, proposal.Title,
		proposal.Description,
		proposal.Proposer,
		proposal.ProposalRoute,
		proposal.ProposalType,
		proposal.ProposalID,
		proposal.Status.String(),
		proposal.SubmitTime,
		proposal.DepositEndTime,
		proposal.VotingStartTime,
		proposal.VotingEndTime)
	return err
}

// SaveTallyResults allows to save for the given height the given total amount of coins
func (db *BigDipperDb) SaveTallyResults(tallys []types.TallyResult) error {
	if len(tallys) == 0 {
		return nil
	}
	query := `INSERT INTO tally_result(proposal_id, yes, abstain, no, no_with_veto, height) VALUES`
	var param []interface{}
	for i, tally := range tallys {
		vi := i * 6
		query += fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d),", vi+1, vi+2, vi+3, vi+4, vi+5, vi+6)
		param = append(param, tally.ProposalID,
			tally.Yes,
			tally.Abstain,
			tally.No,
			tally.NoWithVeto,
			tally.Height,
		)
	}
	query = query[:len(query)-1] // Remove trailing ","
	query += " ON CONFLICT DO NOTHING"
	_, err := db.Sql.Exec(query, param...)
	return err
}

// SaveTallyResult insert a single row into tally_result table
func (db *BigDipperDb) SaveTallyResult(tally types.TallyResult) error {
	query := `
INSERT INTO tally_result(proposal_id, yes, abstain, no, no_with_veto, height)
VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT DO NOTHING`
	_, err := db.Sql.Exec(query, tally.ProposalID,
		tally.Yes,
		tally.Abstain,
		tally.No,
		tally.NoWithVeto,
		tally.Height,
	)
	return err
}

// SaveVote allows to save for the given height and the message vote
func (db *BigDipperDb) SaveVote(vote types.Vote) error {
	query := `INSERT INTO vote(proposal_id, voter, option, height) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	_, err := db.Sql.Exec(query,
		vote.ProposalID,
		vote.Voter,
		vote.Option.String(),
		vote.Height,
	)
	return err
}

// SaveDeposit allows to save for the given message deposit and height
func (db *BigDipperDb) SaveDeposit(deposit types.Deposit) error {
	query := `
INSERT INTO deposit(proposal_id, depositor, amount, total_deposit, height) 
VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`
	_, err := db.Sql.Exec(query,
		deposit.ProposalID,
		deposit.Depositor,
		pq.Array(dbtypes.NewDbCoins(deposit.Amount)),
		pq.Array(dbtypes.NewDbCoins(deposit.TotalDeposit)),
		deposit.Height,
	)
	return err
}

// SaveDeposits allows to save multiple deposits
func (db *BigDipperDb) SaveDeposits(deposits []types.Deposit) error {
	if len(deposits) == 0 {
		return nil
	}
	query := `INSERT INTO deposit(proposal_id, depositor, amount, total_deposit, height) VALUES `
	var param []interface{}

	for i, deposit := range deposits {
		vi := i * 5
		query += fmt.Sprintf("($%d,$%d,$%d,$%d,$%d),", vi+1, vi+2, vi+3, vi+4, vi+5)
		param = append(param, deposit.ProposalID,
			deposit.Depositor,
			pq.Array(dbtypes.NewDbCoins(deposit.Amount)),
			pq.Array(dbtypes.NewDbCoins(deposit.TotalDeposit)),
			deposit.Height,
		)
	}
	query = query[:len(query)-1] // Remove trailing ","
	query += " ON CONFLICT DO NOTHING"
	_, err := db.Sql.Exec(query, param...)
	return err
}

func (db *BigDipperDb) UpdateProposal(proposal types.Proposal) error {
	query := `UPDATE proposal SET status = $1, voting_start_time = $2, voting_end_time = $3 where proposal_id = $4`
	_, err := db.Sql.Exec(query,
		proposal.Status.String(),
		proposal.VotingStartTime,
		proposal.VotingEndTime,
		proposal.ProposalID,
	)
	return err
}
