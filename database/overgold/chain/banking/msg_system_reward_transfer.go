package banking

import (
	bankingtypes "git.ooo.ua/vipcoin/chain/x/banking/types"
	"git.ooo.ua/vipcoin/lib/errs"
	"git.ooo.ua/vipcoin/lib/filter"

	"github.com/forbole/bdjuno/v3/database/types"
)

// SaveMsgSystemRewardTransfers - method that create transfers to the "overgold_chain_banking_system_msg_reward_transfer"
func (r Repository) SaveMsgSystemRewardTransfers(transfers *bankingtypes.MsgSystemRewardTransfer, transactionHash string) error {
	query := `INSERT INTO overgold_chain_banking_system_msg_reward_transfer 
		(transaction_hash, creator, wallet_from, wallet_to, asset, amount, extras) 
		VALUES 
		(:transaction_hash, :creator, :wallet_from, :wallet_to, :asset, :amount, :extras)
		 ON CONFLICT (id) DO NOTHING`

	if _, err := r.db.NamedExec(query, toMsgSystemRewardTransferDatabase(transfers, transactionHash)); err != nil {
		return errs.Internal{Cause: err.Error()}
	}

	return nil
}

// GetMsgSystemRewardTransfers - method that get transfers from the "overgold_chain_banking_system_msg_reward_transfer"
func (r Repository) GetMsgSystemRewardTransfers(filter filter.Filter) ([]*bankingtypes.MsgSystemRewardTransfer, error) {
	query, args := filter.Build(
		tableMsgSystemRewardTransfer,
		types.FieldCreator, types.FieldWalletFrom, types.FieldWalletTo,
		types.FieldAsset, types.FieldAmount, types.FieldExtras,
	)

	var result []types.DBMsgSystemRewardTransfer
	if err := r.db.Select(&result, query, args...); err != nil {
		return []*bankingtypes.MsgSystemRewardTransfer{}, errs.Internal{Cause: err.Error()}
	}

	transfers := make([]*bankingtypes.MsgSystemRewardTransfer, 0, len(result))
	for _, transfer := range result {
		transfers = append(transfers, toMsgSystemRewardTransferDomain(transfer))
	}

	return transfers, nil
}
