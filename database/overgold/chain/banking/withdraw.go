package banking

import (
	"context"
	"database/sql"

	bankingtypes "git.ooo.ua/vipcoin/chain/x/banking/types"
	"git.ooo.ua/vipcoin/lib/errs"
	"git.ooo.ua/vipcoin/lib/filter"

	"github.com/forbole/bdjuno/v2/database/types"
)

// SaveWithdraws - method that create withdraws to the "overgold_chain_banking_withdraw" table
func (r Repository) SaveWithdraws(withdraws ...*bankingtypes.Withdraw) error {
	if len(withdraws) == 0 {
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
       (:id, :asset, :amount, :kind, :extras, :timestamp, :tx_hash)`

	queryWithdraw := `INSERT INTO overgold_chain_banking_withdraw
			("id", "wallet")
			VALUES
			(:id, :wallet)`

	for _, withdraw := range withdraws {
		withdrawDB := toWithdrawDatabase(withdraw)

		if _, err := tx.NamedExec(queryBaseTransfer, withdrawDB); err != nil {
			return errs.Internal{Cause: err.Error()}
		}

		if _, err := tx.NamedExec(queryWithdraw, withdrawDB); err != nil {
			return errs.Internal{Cause: err.Error()}
		}
	}

	return tx.Commit()
}

// GetWithdraws - method that get withdraws from the "overgold_chain_banking_withdraw" table
func (r Repository) GetWithdraws(filter filter.Filter) ([]*bankingtypes.Withdraw, error) {
	query, args := filter.ToJoiner().
		PrepareTable(tableTransfers, types.FieldID, types.FieldAsset, types.FieldAmount, types.FieldKind, types.FieldExtras, types.FieldTimestamp, types.FieldTxHash).
		PrepareTable(tableWithdraw, types.FieldID, types.FieldWallet).
		PrepareJoinStatement("INNER JOIN overgold_chain_banking_base_transfers on overgold_chain_banking_base_transfers.id = overgold_chain_banking_withdraw.id").
		Build(tableWithdraw)

	var withdrawsDB []types.DBWithdraw
	if err := r.db.Select(&withdrawsDB, query, args...); err != nil {
		return []*bankingtypes.Withdraw{}, errs.Internal{Cause: err.Error()}
	}

	result := make([]*bankingtypes.Withdraw, 0, len(withdrawsDB))
	for _, withdraw := range withdrawsDB {
		result = append(result, toWithdrawDomain(withdraw))
	}

	return result, nil
}

// UpdateWithdraws - method that update the withdraw in the "overgold_chain_banking_withdraw" table
func (r Repository) UpdateWithdraws(withdraws ...*bankingtypes.Withdraw) error {
	if len(withdraws) == 0 {
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

	queryWithdraw := `UPDATE overgold_chain_banking_withdraw SET
	wallet =:wallet
	WHERE id =:id;
	`

	for _, withdraw := range withdraws {
		withdrawDB := toWithdrawDatabase(withdraw)

		if _, err := tx.NamedExec(queryBaseTransfer, withdrawDB); err != nil {
			return errs.Internal{Cause: err.Error()}
		}

		if _, err := tx.NamedExec(queryWithdraw, withdrawDB); err != nil {
			return errs.Internal{Cause: err.Error()}
		}
	}

	return tx.Commit()
}
