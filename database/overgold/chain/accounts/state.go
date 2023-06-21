package accounts

import (
	accountstypes "git.ooo.ua/vipcoin/chain/x/accounts/types"
	"git.ooo.ua/vipcoin/lib/errs"
	"git.ooo.ua/vipcoin/lib/filter"

	"github.com/forbole/bdjuno/v2/database/types"
)

// SaveState - saves the given state inside the database
func (r Repository) SaveState(msg *accountstypes.MsgSetState, transactionHash string) error {
	query := `INSERT INTO overgold_chain_accounts_set_state 
			(transaction_hash, creator, hash, state, reason) 
			VALUES 
			(:transaction_hash, :creator, :hash, :state, :reason)
			ON CONFLICT (id) DO NOTHING`

	if _, err := r.db.NamedExec(query, toSetStateDatabase(msg, transactionHash)); err != nil {
		return errs.Internal{Cause: err.Error()}
	}

	return nil
}

// GetState - get the given state from database
func (r Repository) GetState(accountFilter filter.Filter) ([]*accountstypes.MsgSetState, error) {
	query, args := accountFilter.Build(
		tableState,
		types.FieldCreator, types.FieldHash,
		types.FieldState, types.FieldReason,
	)

	var result []types.DBSetState
	if err := r.db.Select(&result, query, args...); err != nil {
		return []*accountstypes.MsgSetState{}, errs.Internal{Cause: err.Error()}
	}

	states := make([]*accountstypes.MsgSetState, 0, len(result))
	for _, state := range result {
		states = append(states, toSetStateDomain(state))
	}

	return states, nil
}
