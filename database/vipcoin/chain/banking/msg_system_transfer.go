package banking

import (
	bankingtypes "git.ooo.ua/vipcoin/chain/x/banking/types"
	"git.ooo.ua/vipcoin/lib/filter"

	"github.com/forbole/bdjuno/v2/database/types"
)

// SaveSystemTransfers - method that create system transfers to the "vipcoin_chain_banking_msg_system_transfer" table
func (r Repository) SaveMsgSystemTransfers(transfers ...*bankingtypes.MsgSystemTransfer) error {
	if len(transfers) == 0 {
		return nil
	}

	query := `INSERT INTO vipcoin_chain_banking_msg_system_transfer 
		(creator, wallet_from, wallet_to, asset, amount, extras) 
		VALUES 
		(:creator, :wallet_from, :wallet_to, :asset, :amount, :extras)`

	if _, err := r.db.NamedExec(query, toMsgSystemTransfersDatabase(transfers...)); err != nil {
		return err
	}

	return nil
}

// GetSystemTransfers - method that get system transfers from the "vipcoin_chain_banking_msg_system_transfer" table
func (r Repository) GetMsgSystemTransfers(filter filter.Filter) ([]*bankingtypes.MsgSystemTransfer, error) {
	query, args := filter.Build(
		tableMsgSystemTransfer,
		types.FieldCreator, types.FieldWalletFrom, types.FieldWalletTo,
		types.FieldAsset, types.FieldAmount, types.FieldExtras,
	)

	var result []types.DBMsgSystemTransfer
	if err := r.db.Select(&result, query, args...); err != nil {
		return []*bankingtypes.MsgSystemTransfer{}, err
	}

	transfers := make([]*bankingtypes.MsgSystemTransfer, 0, len(result))
	for _, transfer := range result {
		transfers = append(transfers, toMsgSystemTransferDomain(transfer))
	}

	return transfers, nil
}
