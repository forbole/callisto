package banking

import (
	bankingtypes "git.ooo.ua/vipcoin/chain/x/banking/types"
	"git.ooo.ua/vipcoin/lib/filter"

	"github.com/forbole/bdjuno/v2/database/types"
)

// SaveMsgSystemRewardTransfers - method that create transfers to the "vipcoin_chain_banking_system_reward_transfer"
func (r Repository) SaveMsgSystemRewardTransfers(transfers ...*bankingtypes.MsgSystemRewardTransfer) error {
	if len(transfers) == 0 {
		return nil
	}

	query := `INSERT INTO vipcoin_chain_banking_system_reward_transfer 
		(creator, wallet_from, wallet_to, asset, amount, extras) 
		VALUES 
		(:creator, :wallet_from, :wallet_to, :asset, :amount, :extras)`

	if _, err := r.db.NamedExec(query, toMsgSystemRewardTransfersDatabase(transfers...)); err != nil {
		return err
	}

	return nil
}

// GetMsgSystemRewardTransfers - method that get transfers from the "vipcoin_chain_banking_system_reward_transfer"
func (r Repository) GetMsgSystemRewardTransfers(filter filter.Filter) ([]*bankingtypes.MsgSystemRewardTransfer, error) {
	query, args := filter.Build(
		tableMsgSystemRewardTransfer,
		types.FieldCreator, types.FieldWalletFrom, types.FieldWalletTo,
		types.FieldAsset, types.FieldAmount, types.FieldExtras,
	)

	var result []types.DBSystemRewardTransfer
	if err := r.db.Select(&result, query, args...); err != nil {
		return []*bankingtypes.MsgSystemRewardTransfer{}, err
	}

	transfers := make([]*bankingtypes.MsgSystemRewardTransfer, 0, len(result))
	for _, transfer := range result {
		transfers = append(transfers, toMsgSystemRewardTransferDomain(transfer))
	}

	return transfers, nil
}
