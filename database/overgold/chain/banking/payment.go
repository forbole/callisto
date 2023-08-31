package banking

import (
	"context"
	"database/sql"
	"errors"

	bankingtypes "git.ooo.ua/vipcoin/chain/x/banking/types"
	"git.ooo.ua/vipcoin/lib/errs"
	"git.ooo.ua/vipcoin/lib/filter"
	"github.com/lib/pq"

	"github.com/forbole/bdjuno/v3/database/types"
)

// SavePayments - method that create payments to the "overgold_chain_banking_payment" table
func (r Repository) SavePayments(payments ...*bankingtypes.Payment) error {
	if len(payments) == 0 {
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
       (:id,:asset, :amount, :kind, :extras, :timestamp, :tx_hash)`

	queryPayment := `INSERT INTO overgold_chain_banking_payment
			("id", "wallet_from", "wallet_to", "fee")
			VALUES
			(:id, :wallet_from, :wallet_to, :fee)`

	var pgErr *pq.Error
	for _, payment := range payments {
		paymentDB := toPaymentDatabase(payment)

		if _, err := tx.NamedExec(queryBaseTransfer, paymentDB); err != nil {
			if errors.As(err, &pgErr) {
				if pgErr.Code == "23505" {
					return errs.AlreadyExists{}
				}
			}
			return errs.Internal{Cause: err.Error()}
		}

		if _, err := tx.NamedExec(queryPayment, paymentDB); err != nil {
			return errs.Internal{Cause: err.Error()}
		}
	}

	return tx.Commit()
}

// GetPayments - method that get payments from the "overgold_chain_banking_payment" table
func (r Repository) GetPayments(filter filter.Filter) ([]*bankingtypes.Payment, error) {
	query, args := filter.ToJoiner().
		PrepareTable(tableTransfers, types.FieldID, types.FieldAsset, types.FieldAmount, types.FieldKind, types.FieldExtras, types.FieldTimestamp, types.FieldTxHash).
		PrepareTable(tablePayment, types.FieldID, types.FieldWalletFrom, types.FieldWalletTo, types.FieldFee).
		PrepareJoinStatement("INNER JOIN overgold_chain_banking_base_transfers on overgold_chain_banking_base_transfers.id = overgold_chain_banking_payment.id").
		Build(tablePayment)

	var paymentsDB []types.DBPayment
	if err := r.db.Select(&paymentsDB, query, args...); err != nil {
		return []*bankingtypes.Payment{}, errs.Internal{Cause: err.Error()}
	}

	result := make([]*bankingtypes.Payment, 0, len(paymentsDB))
	for _, payment := range paymentsDB {
		result = append(result, toPaymentDomain(payment))
	}

	return result, nil
}

// UpdatePayments - method that update the payment in the "overgold_chain_banking_payment" table
func (r Repository) UpdatePayments(payments ...*bankingtypes.Payment) error {
	if len(payments) == 0 {
		return nil
	}

	tx, err := r.db.BeginTxx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	defer tx.Rollback()

	queryBaseTransfer := `UPDATE overgold_chain_banking_base_transfers SET
	asset =:asset, amount =:amount, kind =:kind, extras =:extras, timestamp =:timestamp, tx_hash =:tx_hash
	WHERE id =:id;
	`

	queryPayment := `UPDATE overgold_chain_banking_payment SET
	wallet_from =:wallet_from, wallet_to =:wallet_to, fee =:fee
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
