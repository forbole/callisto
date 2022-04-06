package banking

import (
	bankingtypes "git.ooo.ua/vipcoin/chain/x/banking/types"
	"git.ooo.ua/vipcoin/lib/filter"

	"github.com/forbole/bdjuno/v2/database/types"
)

// SaveMsgPayments - method that create payments to the "vipcoin_chain_banking_msg_payment" table
func (r Repository) SaveMsgPayments(payments ...*bankingtypes.MsgPayment) error {
	if len(payments) == 0 {
		return nil
	}

	query := `INSERT INTO vipcoin_chain_banking_msg_payment 
		(creator, wallet_from, wallet_to, asset, amount, extras) 
		VALUES 
		(:creator, :wallet_from, :wallet_to, :asset, :amount, :extras)`

	if _, err := r.db.NamedExec(query, toMsgPaymentsDatabase(payments...)); err != nil {
		return err
	}

	return nil
}

// GetMsgPayments - method that get payments from the "vipcoin_chain_banking_msg_payment" table
func (r Repository) GetMsgPayments(filter filter.Filter) ([]*bankingtypes.MsgPayment, error) {
	query, args := filter.Build(
		tableMsgPayment,
		types.FieldCreator, types.FieldWalletFrom, types.FieldWalletTo,
		types.FieldAsset, types.FieldAmount, types.FieldExtras,
	)

	var result []types.DBMsgPayment
	if err := r.db.Select(&result, query, args...); err != nil {
		return []*bankingtypes.MsgPayment{}, err
	}

	payments := make([]*bankingtypes.MsgPayment, 0, len(result))
	for _, payment := range result {
		payments = append(payments, toMsgPaymentDomain(payment))
	}

	return payments, nil
}
