package accounts

import (
	accountstypes "git.ooo.ua/vipcoin/chain/x/accounts/types"
	"git.ooo.ua/vipcoin/lib/errs"
	"git.ooo.ua/vipcoin/lib/filter"

	"github.com/forbole/bdjuno/v2/database/types"
)

// SaveExtra - saves the given extra inside the database
func (r Repository) SaveExtra(msg *accountstypes.MsgSetExtra, transactionHash string) error {
	query := `INSERT INTO vipcoin_chain_accounts_set_extra 
			(transaction_hash, creator, hash, extras) 
		VALUES 
			(:transaction_hash, :creator, :hash, :extras)`

	if _, err := r.db.NamedExec(query, toSetExtraDatabase(msg, transactionHash)); err != nil {
		return errs.Internal{Cause: err.Error()}
	}

	return nil
}

// GetExtra - get the given extra from database
func (r Repository) GetExtra(accountFilter filter.Filter) ([]*accountstypes.MsgSetExtra, error) {
	query, args := accountFilter.Build(
		tableExtra,
		types.FieldCreator, types.FieldHash, types.FieldExtras,
	)

	var result []types.DBSetAccountExtra
	if err := r.db.Select(&result, query, args...); err != nil {
		return []*accountstypes.MsgSetExtra{}, errs.Internal{Cause: err.Error()}
	}

	migrates := make([]*accountstypes.MsgSetExtra, 0, len(result))
	for _, extra := range result {
		migrates = append(migrates, toSetExtraDomain(extra))
	}

	return migrates, nil
}
