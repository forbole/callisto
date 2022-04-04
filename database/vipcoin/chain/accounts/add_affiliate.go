package accounts

import (
	accountstypes "git.ooo.ua/vipcoin/chain/x/accounts/types"
	"git.ooo.ua/vipcoin/lib/filter"

	"github.com/forbole/bdjuno/v2/database/types"
)

// SaveAddAffiliate - saves the given affiliate inside the database
func (r Repository) SaveAddAffiliate(msg ...*accountstypes.MsgAddAffiliate) error {
	if len(msg) == 0 {
		return nil
	}

	query := `INSERT INTO vipcoin_chain_accounts_add_affiliate 
			(creator, account_hash, affiliation_hash, affiliation, extras) 
		VALUES 
			(:creator, :account_hash, :affiliation_hash, :affiliation, :extras)`

	if _, err := r.db.NamedExec(query, toAddAffiliatesDatabase(msg...)); err != nil {
		return err
	}

	return nil
}

// GetAddAffiliate - get the given affiliate from database
func (r Repository) GetAddAffiliate(accountFilter filter.Filter) ([]*accountstypes.MsgAddAffiliate, error) {
	query, args := accountFilter.Build(
		tableAddAffiliate,
		types.FieldCreator, types.FieldAccountHash, types.FieldAffiliationHash,
		types.FieldAffiliation, types.FieldExtras,
	)

	var result []types.DBAddAffiliate
	if err := r.db.Select(&result, query, args...); err != nil {
		return []*accountstypes.MsgAddAffiliate{}, err
	}

	affiliates := make([]*accountstypes.MsgAddAffiliate, 0, len(result))
	for _, affiliate := range result {
		affiliates = append(affiliates, toAddAffiliateDomain(affiliate))
	}

	return affiliates, nil
}
