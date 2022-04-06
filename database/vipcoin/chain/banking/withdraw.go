package banking

import (
	"context"
	"database/sql"

	bankingtypes "git.ooo.ua/vipcoin/chain/x/banking/types"
	"git.ooo.ua/vipcoin/lib/filter"

	"github.com/forbole/bdjuno/v2/database/types"
)

// SaveWithdraws - method that create withdraws to the "vipcoin_chain_banking_withdraw" table
func (r Repository) SaveWithdraws(withdraws ...*bankingtypes.Withdraw) error {
	if len(withdraws) == 0 {
		return nil
	}

	tx, err := r.db.BeginTxx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	defer tx.Rollback()

	queryBaseTransfer := `INSERT INTO vipcoin_chain_banking_base_transfers 
       ("id", "asset", "amount", "kind", "extras", "timestamp", "tx_hash") 
     VALUES 
       (:id, :asset, :amount, :kind, :extras, :timestamp, :tx_hash)`

	queryWithdraw := `INSERT INTO vipcoin_chain_banking_withdraw
			("id", "wallet")
			VALUES
			(:id, :wallet)`

	withdrawDB := toWithdrawsDatabase(withdraws...)

	if _, err := tx.NamedExec(queryBaseTransfer, withdrawDB); err != nil {
		return err
	}

	if _, err := tx.NamedExec(queryWithdraw, withdrawDB); err != nil {
		return err
	}

	return tx.Commit()
}

// GetWithdraws - method that get withdraws from the "vipcoin_chain_banking_withdraw" table
func (r Repository) GetWithdraws(filter filter.Filter) ([]*bankingtypes.Withdraw, error) {
	query, args := filter.ToJoiner().
		PrepareTable(tableTransfers, types.FieldID, types.FieldAsset, types.FieldAmount, types.FieldKind, types.FieldExtras, types.FieldTimestamp, types.FieldTxHash).
		PrepareTable(tableWithdraw, types.FieldID, types.FieldWallet).
		PrepareJoinStatement("INNER JOIN vipcoin_chain_banking_base_transfers on vipcoin_chain_banking_base_transfers.id = vipcoin_chain_banking_withdraw.id").
		Build(tableWithdraw)

	var withdrawsDB []types.DBWithdraw
	if err := r.db.Select(&withdrawsDB, query, args...); err != nil {
		return []*bankingtypes.Withdraw{}, err
	}

	result := make([]*bankingtypes.Withdraw, 0, len(withdrawsDB))
	for _, withdraw := range withdrawsDB {
		result = append(result, toWithdrawDomain(withdraw))
	}

	return result, nil
}

// UpdateWithdraws - method that update the withdraw in the "vipcoin_chain_banking_withdraw" table
func (r Repository) UpdateWithdraws(withdraws ...*bankingtypes.Withdraw) error {
	if len(withdraws) == 0 {
		return nil
	}

	tx, err := r.db.BeginTxx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	defer tx.Rollback()

	queryBaseTransfer := `UPDATE vipcoin_chain_banking_base_transfers SET
	id =:id, asset =:asset, amount =:amount, kind =:kind, extras =:extras, timestamp =:timestamp, tx_hash =:tx_hash
	WHERE id =:id;
	`

	queryWithdraw := `UPDATE vipcoin_chain_banking_withdraw SET
	id =:id, wallet =:wallet
	WHERE id =:id;
	`

	for _, withdraw := range withdraws {
		withdrawDB := toWithdrawDatabase(withdraw)

		if _, err := tx.NamedExec(queryBaseTransfer, withdrawDB); err != nil {
			return err
		}

		if _, err := tx.NamedExec(queryWithdraw, withdrawDB); err != nil {
			return err
		}
	}

	return tx.Commit()
}
