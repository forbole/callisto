package banking

import (
	bankingtypes "git.ooo.ua/vipcoin/chain/x/banking/types"
	"git.ooo.ua/vipcoin/lib/filter"

	"github.com/forbole/bdjuno/v2/database/types"
)

// SaveMsgSetRewardMgrAddress - method that save to the "vipcoin_chain_banking_set_reward_manager_address" table
func (r Repository) SaveMsgSetRewardMgrAddress(addresses ...*bankingtypes.MsgSetRewardManagerAddress) error {
	if len(addresses) == 0 {
		return nil
	}

	query := `INSERT INTO vipcoin_chain_banking_set_reward_manager_address 
		(creator, address) 
		VALUES 
		(:creator, :address)`

	if _, err := r.db.NamedExec(query, toMsgSetRewardMgrAddressesDB(addresses...)); err != nil {
		return err
	}

	return nil
}

// GetMsgSetRewardMgrAddress - method that gets from the "vipcoin_chain_banking_set_reward_manager_address" table
func (r Repository) GetMsgSetRewardMgrAddress(filter filter.Filter) ([]*bankingtypes.MsgSetRewardManagerAddress, error) {
	query, args := filter.Build(
		tableMsgSetRewardManagerAddress,
		types.FieldCreator, types.FieldAddress,
	)

	var result []types.DBSetRewardManagerAddress
	if err := r.db.Select(&result, query, args...); err != nil {
		return []*bankingtypes.MsgSetRewardManagerAddress{}, err
	}

	addresses := make([]*bankingtypes.MsgSetRewardManagerAddress, 0, len(result))
	for _, address := range result {
		addresses = append(addresses, toMsgSetRewardMgrAddressDomain(address))
	}

	return addresses, nil
}
