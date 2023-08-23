package accounts

import (
	accountstypes "git.ooo.ua/vipcoin/chain/x/accounts/types"
	"git.ooo.ua/vipcoin/lib/errs"
	"git.ooo.ua/vipcoin/lib/filter"

	"github.com/forbole/bdjuno/v3/database/types"
)

// SaveAffiliateAddress - saves the given affiliate address inside the database
func (r Repository) SaveAffiliateAddress(msg *accountstypes.MsgSetAffiliateAddress, transactionHash string) error {
	query := `INSERT INTO overgold_chain_accounts_set_affiliate_address 
			(transaction_hash, creator, hash, old_address, new_address) 
		VALUES 
			(:transaction_hash, :creator, :hash, :old_address, :new_address)
			ON CONFLICT (id) DO NOTHING`

	if _, err := r.db.NamedExec(query, toSetAffiliateAddressDatabase(msg, transactionHash)); err != nil {
		return errs.Internal{Cause: err.Error()}
	}

	return nil
}

// GetAffiliateAddress - get the given affiliate address from database
func (r Repository) GetAffiliateAddress(accountFilter filter.Filter) ([]*accountstypes.MsgSetAffiliateAddress, error) {
	query, args := accountFilter.Build(
		tableAffiliateAddress,
		types.FieldCreator, types.FieldHash,
		types.FieldOldAddress, types.FieldNewAddress,
	)

	var result []types.DBSetAffiliateAddress
	if err := r.db.Select(&result, query, args...); err != nil {
		return []*accountstypes.MsgSetAffiliateAddress{}, errs.Internal{Cause: err.Error()}
	}

	affiliates := make([]*accountstypes.MsgSetAffiliateAddress, 0, len(result))
	for _, affiliate := range result {
		affiliates = append(affiliates, toSetAffiliateAddressDomain(affiliate))
	}

	return affiliates, nil
}
