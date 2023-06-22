package wallets

import (
	"context"
	"database/sql"

	walletstypes "git.ooo.ua/vipcoin/chain/x/wallets/types"
	"git.ooo.ua/vipcoin/lib/errs"
	"git.ooo.ua/vipcoin/lib/filter"

	"github.com/forbole/bdjuno/v2/database/types"
)

// SaveStates - inserts into the "overgold_chain_wallets_set_wallet_state" table
func (r Repository) SaveStates(messages *walletstypes.MsgSetWalletState, transactionHash string) error {
	query := `INSERT INTO overgold_chain_wallets_set_wallet_state 
			(transaction_hash, creator, address, state) 
			VALUES 
			(:transaction_hash, :creator, :address, :state)
			ON CONFLICT (id) DO NOTHING`

	if _, err := r.db.NamedExec(query, toSetWalletStateDatabase(messages, transactionHash)); err != nil {
		return errs.Internal{Cause: err.Error()}
	}

	return nil
}

// GetStates - get states from the "overgold_chain_wallets_set_wallet_state" table
func (r Repository) GetStates(filter filter.Filter) ([]*walletstypes.MsgSetWalletState, error) {
	query, args := filter.Build("overgold_chain_wallets_set_wallet_state",
		`creator, address, state`)

	var result []types.DBSetWalletState
	if err := r.db.Select(&result, query, args...); err != nil {
		return []*walletstypes.MsgSetWalletState{}, errs.Internal{Cause: err.Error()}
	}

	states := make([]*walletstypes.MsgSetWalletState, 0, len(result))
	for _, kind := range result {
		states = append(states, toSetWalletStateDomain(kind))
	}

	return states, nil
}

// UpdateStates - method that updates states in the "overgold_chain_wallets_set_wallet_state" table
func (r Repository) UpdateStates(messages ...*walletstypes.MsgSetWalletState) error {
	if len(messages) == 0 {
		return nil
	}

	tx, err := r.db.BeginTxx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	defer tx.Rollback()

	query := `UPDATE overgold_chain_wallets_set_wallet_state SET
				 creator = :creator, address = :address, state = :state
				 WHERE address = :address`

	for _, m := range messages {
		stateDB := toSetWalletStateDatabase(m, "")
		if err != nil {
			return errs.Internal{Cause: err.Error()}
		}

		if _, err := tx.NamedExec(query, stateDB); err != nil {
			return errs.Internal{Cause: err.Error()}
		}
	}

	return tx.Commit()
}
