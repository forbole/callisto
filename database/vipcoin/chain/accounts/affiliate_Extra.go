package accounts

import (
	accountstypes "git.ooo.ua/vipcoin/chain/x/accounts/types"
	"git.ooo.ua/vipcoin/lib/filter"
	"github.com/forbole/bdjuno/v2/database/types"
)

// SaveAffiliateExtra - saves the given affiliate extra inside the database
func (r Repository) SaveAffiliateExtra(msg ...*accountstypes.MsgSetAffiliateExtra) error {
	if len(msg) == 0 {
		return nil
	}

	query := `INSERT INTO vipcoin_chain_accounts_set_affiliate_extra 
			(creator, account_hash, affiliation_hash, extras) 
		VALUES 
			(:creator, :account_hash, :affiliation_hash, :extras)`

	if _, err := r.db.NamedExec(query, toSetAffiliatesExtraDatabase(msg...)); err != nil {
		return err
	}

	return nil
}

// GetAffiliateExtra - get the given affiliate extra from database
func (r Repository) GetAffiliateExtra(accfilter filter.Filter) ([]*accountstypes.MsgSetAffiliateExtra, error) {
	query, args := accfilter.Build("vipcoin_chain_accounts_set_affiliate_extra",
		`creator, account_hash, affiliation_hash, extras`)

	var result []types.DBSetAffiliateExtra
	if err := r.db.Select(&result, query, args...); err != nil {
		return []*accountstypes.MsgSetAffiliateExtra{}, err
	}

	affiliates := make([]*accountstypes.MsgSetAffiliateExtra, 0, len(result))
	for _, affiliate := range result {
		affiliates = append(affiliates, toSetAffiliateExtraDomain(affiliate))
	}

	return affiliates, nil
}
