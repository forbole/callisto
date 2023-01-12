package banking

import (
	bankingtypes "git.ooo.ua/vipcoin/chain/x/banking/types"
	"git.ooo.ua/vipcoin/lib/errs"
	"git.ooo.ua/vipcoin/lib/filter"

	"github.com/forbole/bdjuno/v2/database/types"
)

// SaveMsgWithdraw - method that create withdraw to the "overgold_chain_banking_msg_withdraw" table
func (r Repository) SaveMsgWithdraw(withdraws *bankingtypes.MsgWithdraw, transactionHash string) error {
	query := `INSERT INTO overgold_chain_banking_msg_withdraw 
		(transaction_hash, creator, wallet, asset, amount, extras) 
		VALUES 
		(:transaction_hash, :creator, :wallet, :asset, :amount, :extras)`

	if _, err := r.db.NamedExec(query, toMsgWithdrawDatabase(withdraws, transactionHash)); err != nil {
		return errs.Internal{Cause: err.Error()}
	}

	return nil
}

// GetSystemTransfers - method that get withdraw from the "overgold_chain_banking_msg_withdraw" table
func (r Repository) GetMsgWithdraw(filter filter.Filter) ([]*bankingtypes.MsgWithdraw, error) {
	query, args := filter.Build(
		tableMsgWithdraw,
		types.FieldCreator, types.FieldWallet,
		types.FieldAsset, types.FieldAmount, types.FieldExtras,
	)

	var result []types.DBMsgWithdraw
	if err := r.db.Select(&result, query, args...); err != nil {
		return []*bankingtypes.MsgWithdraw{}, errs.Internal{Cause: err.Error()}
	}

	withdraws := make([]*bankingtypes.MsgWithdraw, 0, len(result))
	for _, withdraw := range result {
		withdraws = append(withdraws, toMsgWithdrawDomain(withdraw))
	}

	return withdraws, nil
}
