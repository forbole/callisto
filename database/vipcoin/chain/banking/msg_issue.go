package banking

import (
	bankingtypes "git.ooo.ua/vipcoin/chain/x/banking/types"
	"git.ooo.ua/vipcoin/lib/filter"

	"github.com/forbole/bdjuno/v2/database/types"
)

// SaveMsgIssue - method that create issue to the "vipcoin_chain_banking_msg_issue" table
func (r Repository) SaveMsgIssue(issue ...*bankingtypes.MsgIssue) error {
	if len(issue) == 0 {
		return nil
	}

	query := `INSERT INTO vipcoin_chain_banking_msg_issue 
		(creator, wallet, asset, amount, extras) 
		VALUES 
		(:creator, :wallet, :asset, :amount, :extras)`

	if _, err := r.db.NamedExec(query, toMsgIssuesDatabase(issue...)); err != nil {
		return err
	}

	return nil
}

// GetMsgIssue - method that get issue from the "vipcoin_chain_banking_msg_issue" table
func (r Repository) GetMsgIssue(filter filter.Filter) ([]*bankingtypes.MsgIssue, error) {
	query, args := filter.Build(
		tableMsgIssue,
		types.FieldCreator, types.FieldWallet,
		types.FieldAsset, types.FieldAmount, types.FieldExtras,
	)

	var result []types.DBMsgIssue
	if err := r.db.Select(&result, query, args...); err != nil {
		return []*bankingtypes.MsgIssue{}, err
	}

	issues := make([]*bankingtypes.MsgIssue, 0, len(result))
	for _, issue := range result {
		issues = append(issues, toMsgIssueDomain(issue))
	}

	return issues, nil
}
