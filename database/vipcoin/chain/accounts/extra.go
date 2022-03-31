package accounts

import (
	accountstypes "git.ooo.ua/vipcoin/chain/x/accounts/types"
	"git.ooo.ua/vipcoin/lib/filter"

	"github.com/forbole/bdjuno/v2/database/types"
)

// SaveExtra - saves the given extra inside the database
func (r Repository) SaveExtra(msg ...*accountstypes.MsgSetExtra) error {
	if len(msg) == 0 {
		return nil
	}

	query := `INSERT INTO vipcoin_chain_accounts_set_extra 
			(creator, hash, extras) 
		VALUES 
			(:creator, :hash, :extras)`

	if _, err := r.db.NamedExec(query, toSetExtrasDatabase(msg...)); err != nil {
		return err
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
		return []*accountstypes.MsgSetExtra{}, err
	}

	migrates := make([]*accountstypes.MsgSetExtra, 0, len(result))
	for _, extra := range result {
		migrates = append(migrates, toSetExtraDomain(extra))
	}

	return migrates, nil
}
