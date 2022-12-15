package accounts

import (
	accountstypes "git.ooo.ua/vipcoin/chain/x/accounts/types"
	"git.ooo.ua/vipcoin/lib/errs"
	"git.ooo.ua/vipcoin/lib/filter"

	"github.com/forbole/bdjuno/v2/database/types"
)

// SaveAddAffiliate - saves the given affiliate inside the database
func (r Repository) SaveAddAffiliate(msg *accountstypes.MsgAddAffiliate, transactionHash string) error {
	query := `INSERT INTO overgold_chain_accounts_add_affiliate 
			(transaction_hash, creator, account_hash, affiliation_hash, affiliation, extras) 
		VALUES 
			(:transaction_hash, :creator, :account_hash, :affiliation_hash, :affiliation, :extras)`

	if _, err := r.db.NamedExec(query, toAddAffiliateDatabase(msg, transactionHash)); err != nil {
		return errs.Internal{Cause: err.Error()}
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
		return []*accountstypes.MsgAddAffiliate{}, errs.Internal{Cause: err.Error()}
	}

	affiliates := make([]*accountstypes.MsgAddAffiliate, 0, len(result))
	for _, affiliate := range result {
		affiliates = append(affiliates, toAddAffiliateDomain(affiliate))
	}

	return affiliates, nil
}
