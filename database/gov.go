package database

import (
	"fmt"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/gogo/protobuf/proto"

	"github.com/forbole/bdjuno/types"

	dbtypes "github.com/forbole/bdjuno/database/types"

	"github.com/lib/pq"
)

// SaveProposals allows to save for the given height the given total amount of coins
func (db *Db) SaveProposals(proposals []types.Proposal) error {
	//do nothing if empty
	if len(proposals) == 0 {
		return nil
	}

	query := `
INSERT INTO proposal(
	title, description, content, proposer_address, proposal_route, proposal_type, proposal_id, status, 
    submit_time, deposit_end_time, voting_start_time, voting_end_time
) VALUES`
	var param []interface{}
	for i, proposal := range proposals {
		vi := i * 12
		query += fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d),",
			vi+1, vi+2, vi+3, vi+4, vi+5, vi+6, vi+7, vi+8, vi+9, vi+10, vi+11, vi+12)

		// Encode the content properly
		protoContent, ok := proposal.Content.(proto.Message)
		if !ok {
			return fmt.Errorf("invalid proposal content types: %T", proposal.Content)
		}

		anyContent, err := codectypes.NewAnyWithValue(protoContent)
		if err != nil {
			return err
		}

		contentBz, err := db.EncodingConfig.Marshaler.MarshalJSON(anyContent)
		if err != nil {
			return err
		}

		param = append(param,
			proposal.Content.GetTitle(),
			proposal.Content.GetDescription(),
			string(contentBz),
			proposal.Proposer,
			proposal.ProposalRoute,
			proposal.ProposalType,
			proposal.ProposalID,
			proposal.Status.String(),
			proposal.SubmitTime,
			proposal.DepositEndTime,
			proposal.VotingStartTime,
			proposal.VotingEndTime,
		)
	}
	query = query[:len(query)-1] // Remove trailing ","
	query += " ON CONFLICT DO NOTHING"
	_, err := db.Sql.Exec(query, param...)
	return err
}

// SaveTallyResults allows to save for the given height the given total amount of coins
func (db *Db) SaveTallyResults(tallys []types.TallyResult) error {
	if len(tallys) == 0 {
		return nil
	}
	query := `INSERT INTO proposal_tally_result(proposal_id, yes, abstain, no, no_with_veto, height) VALUES`
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

// SaveVote allows to save for the given height and the message vote
func (db *Db) SaveVote(vote types.Vote) error {
	query := `INSERT INTO proposal_vote(proposal_id, voter_address, option, height) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	_, err := db.Sql.Exec(query,
		vote.ProposalID,
		vote.Voter,
		vote.Option.String(),
		vote.Height,
	)
	return err
}

// SaveDeposits allows to save multiple deposits
func (db *Db) SaveDeposits(deposits []types.Deposit) error {
	if len(deposits) == 0 {
		return nil
	}

	query := `INSERT INTO proposal_deposit(proposal_id, depositor_address, amount, height) VALUES `
	var param []interface{}

	for i, deposit := range deposits {
		vi := i * 4
		query += fmt.Sprintf("($%d,$%d,$%d,$%d),", vi+1, vi+2, vi+3, vi+4)
		param = append(param, deposit.ProposalID,
			deposit.Depositor,
			pq.Array(dbtypes.NewDbCoins(deposit.Amount)),
			deposit.Height,
		)
	}
	query = query[:len(query)-1] // Remove trailing ","
	query += " ON CONFLICT DO NOTHING"
	_, err := db.Sql.Exec(query, param...)
	return err
}

// UpdateProposal updates a proposal stored inside the database
func (db *Db) UpdateProposal(update types.ProposalUpdate) error {
	query := `UPDATE proposal SET status = $1, voting_start_time = $2, voting_end_time = $3 where proposal_id = $4`
	_, err := db.Sql.Exec(query,
		update.Status.String(),
		update.VotingStartTime,
		update.VotingEndTime,
		update.ProposalID,
	)
	return err
}
