package accounts

import (
	accountstypes "git.ooo.ua/vipcoin/chain/x/accounts/types"
	"git.ooo.ua/vipcoin/lib/filter"

	"github.com/forbole/bdjuno/v2/database/types"
)

// SaveState - saves the given state inside the database
func (r Repository) SaveState(msg ...*accountstypes.MsgSetState) error {
	if len(msg) == 0 {
		return nil
	}

	query := `INSERT INTO vipcoin_chain_accounts_set_state 
			(creator, hash, state, reason) 
			VALUES 
			(:creator, :hash, :state, :reason)`

	if _, err := r.db.NamedExec(query, toSetStatesDatabase(msg...)); err != nil {
		return err
	}

	return nil
}

// GetState - get the given state from database
func (r Repository) GetState(accfilter filter.Filter) ([]*accountstypes.MsgSetState, error) {
	query, args := accfilter.Build(
		tableState,
		types.FieldCreator, types.FieldHash,
		types.FieldState, types.FieldReason,
	)

	var result []types.DBSetState
	if err := r.db.Select(&result, query, args...); err != nil {
		return []*accountstypes.MsgSetState{}, err
	}

	states := make([]*accountstypes.MsgSetState, 0, len(result))
	for _, state := range result {
		states = append(states, toSetStateDomain(state))
	}

	return states, nil
}
