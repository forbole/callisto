package wallets

import (
	walletstypes "git.ooo.ua/vipcoin/chain/x/wallets/types"
	"git.ooo.ua/vipcoin/lib/filter"

	"github.com/forbole/bdjuno/v2/database/types"
)

// SaveExtras - inserts into the "vipcoin_chain_wallets_set_extra" table
func (r Repository) SaveExtras(messages ...*walletstypes.MsgSetExtra) error {
	if len(messages) == 0 {
		return nil
	}

	query := `INSERT INTO vipcoin_chain_wallets_set_extra 
			(creator, address, extras) 
			VALUES 
			(:creator, :address, :extras)`

	if _, err := r.db.NamedExec(query, toSetExtrasDatabase(messages...)); err != nil {
		return err
	}

	return nil
}

// GetExtras - get extras from the "vipcoin_chain_wallets_set_extra" table
func (r Repository) GetExtras(filter filter.Filter) ([]*walletstypes.MsgSetExtra, error) {
	query, args := filter.Build("vipcoin_chain_wallets_set_extra",
		types.FieldCreator, types.FieldAddress, types.FieldExtras)

	var extrasDB []types.DBSetExtra
	if err := r.db.Select(&extrasDB, query, args...); err != nil {
		return []*walletstypes.MsgSetExtra{}, err
	}

	extras := make([]*walletstypes.MsgSetExtra, 0, len(extrasDB))
	for _, extraDB := range extrasDB {
		extras = append(extras, toSetExtraDomain(extraDB))
	}

	return extras, nil
}
