package banking

import (
	"context"
	"database/sql"

	bankingtypes "git.ooo.ua/vipcoin/chain/x/banking/types"
	"git.ooo.ua/vipcoin/lib/filter"

	"github.com/forbole/bdjuno/v2/database/types"
)

// SavePayments - method that create payments to the "vipcoin_chain_banking_payment" table
func (r Repository) SavePayments(payments ...*bankingtypes.Payment) error {
	if len(payments) == 0 {
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

	queryPayment := `INSERT INTO vipcoin_chain_banking_payment
			("id", "wallet_from", "wallet_to", "fee")
			VALUES
			(:id,:wallet_from,:wallet_to,:fee)`

	paymentDB := toPaymentsDatabase(payments...)

	if _, err := tx.NamedExec(queryBaseTransfer, paymentDB); err != nil {
		return err
	}

	if _, err := tx.NamedExec(queryPayment, paymentDB); err != nil {
		return err
	}

	return tx.Commit()
}

// GetPayments - method that get payments from the "vipcoin_chain_banking_payment" table
func (r Repository) GetPayments(filter filter.Filter) ([]*bankingtypes.Payment, error) {
	query, args := filter.ToJoiner().
		PrepareTable(tableTransfers, types.FieldID, types.FieldAsset, types.FieldAmount, types.FieldKind, types.FieldExtras, types.FieldTimestamp, types.FieldTxHash).
		PrepareTable(tablePayment, types.FieldID, types.FieldWalletFrom, types.FieldWalletTo, types.FieldFee).
		PrepareJoinStatement("INNER JOIN vipcoin_chain_banking_base_transfers on vipcoin_chain_banking_base_transfers.id = vipcoin_chain_banking_payment.id").
		Build(tablePayment)

	var paymentsDB []types.DBPayment
	if err := r.db.Select(&paymentsDB, query, args...); err != nil {
		return []*bankingtypes.Payment{}, err
	}

	result := make([]*bankingtypes.Payment, 0, len(paymentsDB))
	for _, payment := range paymentsDB {
		result = append(result, toPaymentDomain(payment))
	}

	return result, nil
}

// UpdatePayments - method that update the payment in the "vipcoin_chain_banking_payment" table
func (r Repository) UpdatePayments(payments ...*bankingtypes.Payment) error {
	if len(payments) == 0 {
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

	queryPayment := `UPDATE vipcoin_chain_banking_payment SET
	id =:id, wallet_from =:wallet_from, wallet_to =:wallet_to, fee =:fee
	WHERE id =:id;
	`

	for _, payment := range payments {
		paymentDB := toPaymentDatabase(payment)

		if _, err := tx.NamedExec(queryBaseTransfer, paymentDB); err != nil {
			return err
		}

		if _, err := tx.NamedExec(queryPayment, paymentDB); err != nil {
			return err
		}
	}

	return tx.Commit()
}
