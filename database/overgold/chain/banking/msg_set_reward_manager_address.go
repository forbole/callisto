package banking

import (
	bankingtypes "git.ooo.ua/vipcoin/chain/x/banking/types"
	"git.ooo.ua/vipcoin/lib/errs"
	"git.ooo.ua/vipcoin/lib/filter"

	"github.com/forbole/bdjuno/v2/database/types"
)

// SaveMsgSetRewardMgrAddress - method that save to the "overgold_chain_banking_set_reward_manager_address" table
func (r Repository) SaveMsgSetRewardMgrAddress(addresses *bankingtypes.MsgSetRewardManagerAddress, transactionHash string) error {
	query := `INSERT INTO overgold_chain_banking_set_reward_manager_address 
		(transaction_hash, creator, address) 
		VALUES 
		(:transaction_hash, :creator, :address)`

	if _, err := r.db.NamedExec(query, toMsgSetRewardMgrAddressDB(addresses, transactionHash)); err != nil {
		return errs.Internal{Cause: err.Error()}
	}

	return nil
}

// GetMsgSetRewardMgrAddress - method that gets from the "overgold_chain_banking_set_reward_manager_address" table
func (r Repository) GetMsgSetRewardMgrAddress(filter filter.Filter) ([]*bankingtypes.MsgSetRewardManagerAddress, error) {
	query, args := filter.Build(
		tableMsgSetRewardManagerAddress,
		types.FieldCreator, types.FieldAddress,
	)

	var result []types.DBSetRewardManagerAddress
	if err := r.db.Select(&result, query, args...); err != nil {
		return []*bankingtypes.MsgSetRewardManagerAddress{}, errs.Internal{Cause: err.Error()}
	}

	addresses := make([]*bankingtypes.MsgSetRewardManagerAddress, 0, len(result))
	for _, address := range result {
		addresses = append(addresses, toMsgSetRewardMgrAddressDomain(address))
	}

	return addresses, nil
}
