package accounts

import (
	accountstypes "git.ooo.ua/vipcoin/chain/x/accounts/types"
	"git.ooo.ua/vipcoin/lib/errs"
	"git.ooo.ua/vipcoin/lib/filter"

	"github.com/forbole/bdjuno/v3/database/types"
)

// SaveKinds - saves the given kinds inside the database
func (r Repository) SaveKinds(msg *accountstypes.MsgSetKinds, transactionHash string) error {
	query := `INSERT INTO overgold_chain_accounts_set_kinds 
			(transaction_hash, creator, hash, kinds) 
			VALUES 
			(:transaction_hash, :creator, :hash, :kinds)`

	if _, err := r.db.NamedExec(query, toSetKindsDatabase(msg, transactionHash)); err != nil {
		return errs.Internal{Cause: err.Error()}
	}

	return nil
}

// GetKinds - get the given kinds from database
func (r Repository) GetKinds(accountFilter filter.Filter) ([]*accountstypes.MsgSetKinds, error) {
	query, args := accountFilter.Build(
		tableKinds,
		types.FieldCreator, types.FieldHash, types.FieldKinds,
	)

	var result []types.DBSetKinds
	if err := r.db.Select(&result, query, args...); err != nil {
		return []*accountstypes.MsgSetKinds{}, errs.Internal{Cause: err.Error()}
	}

	kinds := make([]*accountstypes.MsgSetKinds, 0, len(result))
	for _, kind := range result {
		kinds = append(kinds, toSetKindsDomain(kind))
	}

	return kinds, nil
}
