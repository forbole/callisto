package accounts

import (
	accountstypes "git.ooo.ua/vipcoin/chain/x/accounts/types"
	"git.ooo.ua/vipcoin/lib/filter"

	"github.com/forbole/bdjuno/v2/database/types"
)

// SaveAffiliateAddress - saves the given affiliate address inside the database
func (r Repository) SaveAffiliateAddress(msg ...*accountstypes.MsgSetAffiliateAddress) error {
	if len(msg) == 0 {
		return nil
	}

	query := `INSERT INTO vipcoin_chain_accounts_set_affiliate_address 
			(creator, hash, old_address, new_address) 
		VALUES 
			(:creator, :hash, :old_address, :new_address)`

	if _, err := r.db.NamedExec(query, toSetAffiliatesAddressDatabase(msg...)); err != nil {
		return err
	}

	return nil
}

// GetAffiliateAddress - get the given affiliate address from database
func (r Repository) GetAffiliateAddress(accfilter filter.Filter) ([]*accountstypes.MsgSetAffiliateAddress, error) {
	query, args := accfilter.Build(
		tableAffiliateAddress,
		types.FieldCreator, types.FieldHash,
		types.FieldOldAddress, types.FieldNewAddress,
	)

	var result []types.DBSetAffiliateAddress
	if err := r.db.Select(&result, query, args...); err != nil {
		return []*accountstypes.MsgSetAffiliateAddress{}, err
	}

	affiliates := make([]*accountstypes.MsgSetAffiliateAddress, 0, len(result))
	for _, affiliate := range result {
		affiliates = append(affiliates, toSetAffiliateAddressDomain(affiliate))
	}

	return affiliates, nil
}
