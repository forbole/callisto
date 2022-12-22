package banking

import (
	bankingtypes "git.ooo.ua/vipcoin/chain/x/banking/types"
	"git.ooo.ua/vipcoin/lib/errs"
	"git.ooo.ua/vipcoin/lib/filter"

	"github.com/forbole/bdjuno/v2/database/types"
)

// SaveMsgIssue - method that create issue to the "overgold_chain_banking_msg_issue" table
func (r Repository) SaveMsgIssue(issue *bankingtypes.MsgIssue, transactionHash string) error {
	query := `INSERT INTO overgold_chain_banking_msg_issue 
		(transaction_hash, creator, wallet, asset, amount, extras) 
		VALUES 
		(:transaction_hash, :creator, :wallet, :asset, :amount, :extras)`

	if _, err := r.db.NamedExec(query, toMsgIssueDatabase(issue, transactionHash)); err != nil {
		return errs.Internal{Cause: err.Error()}
	}

	return nil
}

// GetMsgIssue - method that get issue from the "overgold_chain_banking_msg_issue" table
func (r Repository) GetMsgIssue(filter filter.Filter) ([]*bankingtypes.MsgIssue, error) {
	query, args := filter.Build(
		tableMsgIssue,
		types.FieldCreator, types.FieldWallet,
		types.FieldAsset, types.FieldAmount, types.FieldExtras,
	)

	var result []types.DBMsgIssue
	if err := r.db.Select(&result, query, args...); err != nil {
		return []*bankingtypes.MsgIssue{}, errs.Internal{Cause: err.Error()}
	}

	issues := make([]*bankingtypes.MsgIssue, 0, len(result))
	for _, issue := range result {
		issues = append(issues, toMsgIssueDomain(issue))
	}

	return issues, nil
}
