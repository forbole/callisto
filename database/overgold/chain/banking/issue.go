package banking

import (
	"context"
	"database/sql"

	bankingtypes "git.ooo.ua/vipcoin/chain/x/banking/types"
	"git.ooo.ua/vipcoin/lib/errs"
	"git.ooo.ua/vipcoin/lib/filter"

	"github.com/forbole/bdjuno/v2/database/types"
)

// SaveIssues - method that create issues to the "overgold_chain_banking_issue" table
func (r Repository) SaveIssues(issues ...*bankingtypes.Issue) error {
	if len(issues) == 0 {
		return nil
	}

	tx, err := r.db.BeginTxx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return errs.Internal{Cause: err.Error()}
	}

	defer tx.Rollback()

	queryBaseTransfer := `INSERT INTO overgold_chain_banking_base_transfers 
       ("id", "asset", "amount", "kind", "extras", "timestamp", "tx_hash") 
     VALUES 
       (:id, :asset, :amount, :kind, :extras, :timestamp, :tx_hash)
       ON CONFLICT (id) DO NOTHING`

	queryIssue := `INSERT INTO overgold_chain_banking_issue
			("id", "wallet")
			VALUES
			(:id, :wallet)
			ON CONFLICT (id) DO NOTHING`

	for _, issue := range issues {
		issueDB := toIssueDatabase(issue)

		if _, err := tx.NamedExec(queryBaseTransfer, issueDB); err != nil {
			return errs.Internal{Cause: err.Error()}
		}

		if _, err := tx.NamedExec(queryIssue, issueDB); err != nil {
			return errs.Internal{Cause: err.Error()}
		}
	}

	return tx.Commit()
}

// GetIssues - method that get issues from the "overgold_chain_banking_issue" table
func (r Repository) GetIssues(filter filter.Filter) ([]*bankingtypes.Issue, error) {
	query, args := filter.ToJoiner().
		PrepareTable(tableTransfers, types.FieldID, types.FieldAsset, types.FieldAmount, types.FieldKind, types.FieldExtras, types.FieldTimestamp, types.FieldTxHash).
		PrepareTable(tableIssue, types.FieldID, types.FieldWallet).
		PrepareJoinStatement("INNER JOIN overgold_chain_banking_base_transfers on overgold_chain_banking_base_transfers.id = overgold_chain_banking_issue.id").
		Build(tableIssue)

	var issuesDB []types.DBIssue
	if err := r.db.Select(&issuesDB, query, args...); err != nil {
		return []*bankingtypes.Issue{}, errs.Internal{Cause: err.Error()}
	}

	result := make([]*bankingtypes.Issue, 0, len(issuesDB))
	for _, issue := range issuesDB {
		result = append(result, toIssueDomain(issue))
	}

	return result, nil
}

// UpdateIssues - method that update the issue in the "overgold_chain_banking_issue" table
func (r Repository) UpdateIssues(issues ...*bankingtypes.Issue) error {
	if len(issues) == 0 {
		return nil
	}

	tx, err := r.db.BeginTxx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return errs.Internal{Cause: err.Error()}
	}

	defer tx.Rollback()

	queryBaseTransfer := `UPDATE overgold_chain_banking_base_transfers SET
	asset =:asset, amount =:amount, kind =:kind, extras =:extras, timestamp =:timestamp, tx_hash =:tx_hash
	WHERE id =:id;
	`

	queryIssue := `UPDATE overgold_chain_banking_issue SET
	wallet =:wallet
	WHERE id =:id;
	`

	for _, issue := range issues {
		issueDB := toIssuesDatabase(issue)

		if _, err := tx.NamedExec(queryBaseTransfer, issueDB); err != nil {
			return errs.Internal{Cause: err.Error()}
		}

		if _, err := tx.NamedExec(queryIssue, issueDB); err != nil {
			return errs.Internal{Cause: err.Error()}
		}
	}

	return tx.Commit()
}
