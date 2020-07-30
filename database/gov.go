package database

import (
	"fmt"
	"time"

	dbtypes "github.com/forbole/bdjuno/database/types"
	types "github.com/forbole/bdjuno/x/gov/types"
	api "github.com/forbole/bdjuno/x/pricefeed/apiTypes"
	"github.com/lib/pq"
)

// SaveProposals allows to save for the given height the given total amount of coins
func (db BigDipperDb) SaveProposals(proposals []types.Proposal) error {
	query := `INSERT INTO proposal(title,description ,proposal_route ,proposal_type,proposal_ID,
		status,submit_time ,deposit_end_time ,total_deposit,voting_start_time,voting_end_time) VALUES`
	var param []interface{}
	for i, proposal := range proposals {
		vi := i * 11
		query += fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d),",
			vi+1, vi+2, vi+3, vi+4, vi+5, vi+6, vi+7, vi+8, vi+9, vi+10, vi+11)
		param = append(param, proposal.Title,
			proposal.Description,
			proposal.ProposalRoute,
			proposal.ProposalType,
			proposal.ProposalID,
			proposal.Status,
			proposal.SubmitTime,
			proposal.DepositEndTime,
			pq.Array(dbtypes.NewDbCoins(proposal.TotalDeposit)),
			proposal.VotingStartTime,
			proposal.VotingEndTime)
	}
	query = query[:len(query)-1] // Remove trailing ","
	_, err := db.Sql.Exec(query, param...)
	if err != nil {
		return err
	}
	return nil
}

// SaveTallyResult allows to save for the given height the given total amount of coins
func (db BigDipperDb) SaveTallyResults(pricefeeds api.MarketTickers, timestamp time.Time) error {
	query := `INSERT INTO token_price(denom,price,market_cap,timestamp) VALUES`
	var param []interface{}
	for i, pricefeed := range pricefeeds {
		vi := i * 4
		query += fmt.Sprintf("($%d,$%d,$%d,$%d),", vi+1, vi+2, vi+3, vi+4)
		param = append(param, pricefeed.ID, pricefeed.CurrentPrice, pricefeed.MarketCap, timestamp)
	}
	query = query[:len(query)-1] // Remove trailing ","
	_, err := db.Sql.Exec(query, param...)
	if err != nil {
		return err
	}
	return nil
}
