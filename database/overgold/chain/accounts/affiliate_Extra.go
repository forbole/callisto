package accounts

import (
	accountstypes "git.ooo.ua/vipcoin/chain/x/accounts/types"
	"git.ooo.ua/vipcoin/lib/errs"
	"git.ooo.ua/vipcoin/lib/filter"

	"github.com/forbole/bdjuno/v3/database/types"
)

// SaveAffiliateExtra - saves the given affiliate extra inside the database
func (r Repository) SaveAffiliateExtra(msg *accountstypes.MsgSetAffiliateExtra, transactionHash string) error {
	query := `INSERT INTO overgold_chain_accounts_set_affiliate_extra 
			(transaction_hash, creator, account_hash, affiliation_hash, extras) 
		VALUES 
			(:transaction_hash, :creator, :account_hash, :affiliation_hash, :extras)`

	if _, err := r.db.NamedExec(query, toSetAffiliateExtraDatabase(msg, transactionHash)); err != nil {
		return errs.Internal{Cause: err.Error()}
	}

	return nil
}

// GetAffiliateExtra - get the given affiliate extra from database
func (r Repository) GetAffiliateExtra(accountFilter filter.Filter) ([]*accountstypes.MsgSetAffiliateExtra, error) {
	query, args := accountFilter.Build(
		tableAffiliateExtra,
		types.FieldCreator, types.FieldAccountHash,
		types.FieldAffiliationHash, types.FieldExtras,
	)

	var result []types.DBSetAffiliateExtra
	if err := r.db.Select(&result, query, args...); err != nil {
		return []*accountstypes.MsgSetAffiliateExtra{}, errs.Internal{Cause: err.Error()}
	}

	affiliates := make([]*accountstypes.MsgSetAffiliateExtra, 0, len(result))
	for _, affiliate := range result {
		affiliates = append(affiliates, toSetAffiliateExtraDomain(affiliate))
	}

	return affiliates, nil
}
