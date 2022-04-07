package banking

import (
	bankingtypes "git.ooo.ua/vipcoin/chain/x/banking/types"
	"git.ooo.ua/vipcoin/lib/filter"

	"github.com/forbole/bdjuno/v2/database/types"
)

// SaveMsgWithdraw - method that create withdraw to the "vipcoin_chain_banking_msg_withdraw" table
func (r Repository) SaveMsgWithdraw(withdraws ...*bankingtypes.MsgWithdraw) error {
	if len(withdraws) == 0 {
		return nil
	}

	query := `INSERT INTO vipcoin_chain_banking_msg_withdraw 
		(creator, wallet, asset, amount, extras) 
		VALUES 
		(:creator, :wallet, :asset, :amount, :extras)`

	if _, err := r.db.NamedExec(query, toMsgWithdrawsDatabase(withdraws...)); err != nil {
		return err
	}

	return nil
}

// GetSystemTransfers - method that get withdraw from the "vipcoin_chain_banking_msg_withdraw" table
func (r Repository) GetMsgWithdraw(filter filter.Filter) ([]*bankingtypes.MsgWithdraw, error) {
	query, args := filter.Build(
		tableMsgWithdraw,
		types.FieldCreator, types.FieldWallet,
		types.FieldAsset, types.FieldAmount, types.FieldExtras,
	)

	var result []types.DBMsgWithdraw
	if err := r.db.Select(&result, query, args...); err != nil {
		return []*bankingtypes.MsgWithdraw{}, err
	}

	withdraws := make([]*bankingtypes.MsgWithdraw, 0, len(result))
	for _, withdraw := range result {
		withdraws = append(withdraws, toMsgWithdrawDomain(withdraw))
	}

	return withdraws, nil
}
