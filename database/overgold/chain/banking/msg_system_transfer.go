package banking

import (
	bankingtypes "git.ooo.ua/vipcoin/chain/x/banking/types"
	"git.ooo.ua/vipcoin/lib/errs"
	"git.ooo.ua/vipcoin/lib/filter"

	"github.com/forbole/bdjuno/v3/database/types"
)

// SaveMsgSystemTransfers - method that create system transfers to the "overgold_chain_banking_msg_system_transfer" table
func (r Repository) SaveMsgSystemTransfers(transfers *bankingtypes.MsgSystemTransfer, transactionHash string) error {
	query := `INSERT INTO overgold_chain_banking_msg_system_transfer 
		(transaction_hash, creator, wallet_from, wallet_to, asset, amount, extras) 
		VALUES 
		(:transaction_hash, :creator, :wallet_from, :wallet_to, :asset, :amount, :extras)`

	if _, err := r.db.NamedExec(query, toMsgSystemTransferDatabase(transfers, transactionHash)); err != nil {
		return errs.Internal{Cause: err.Error()}
	}

	return nil
}

// GetMsgSystemTransfers - method that get system transfers from the "overgold_chain_banking_msg_system_transfer" table
func (r Repository) GetMsgSystemTransfers(filter filter.Filter) ([]*bankingtypes.MsgSystemTransfer, error) {
	query, args := filter.Build(
		tableMsgSystemTransfer,
		types.FieldCreator, types.FieldWalletFrom, types.FieldWalletTo,
		types.FieldAsset, types.FieldAmount, types.FieldExtras,
	)

	var transfersDB []types.DBMsgSystemTransfer
	if err := r.db.Select(&transfersDB, query, args...); err != nil {
		return []*bankingtypes.MsgSystemTransfer{}, errs.Internal{Cause: err.Error()}
	}

	transfers := make([]*bankingtypes.MsgSystemTransfer, 0, len(transfersDB))
	for _, transfer := range transfersDB {
		transfers = append(transfers, toMsgSystemTransferDomain(transfer))
	}

	return transfers, nil
}
