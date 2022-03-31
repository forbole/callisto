package wallets

import (
	"context"
	"database/sql"

	walletstypes "git.ooo.ua/vipcoin/chain/x/wallets/types"
	"git.ooo.ua/vipcoin/lib/filter"

	"github.com/forbole/bdjuno/v2/database/types"
)

// SaveStates - inserts into the "vipcoin_chain_wallets_set_wallet_state" table
func (r Repository) SaveStates(messages ...*walletstypes.MsgSetWalletState) error {
	if len(messages) == 0 {
		return nil
	}

	tx, err := r.db.BeginTxx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	defer func() {
		_ = tx.Rollback()
	}()

	query := `INSERT INTO vipcoin_chain_wallets_set_wallet_state 
			(creator, address, state) 
			VALUES 
			(:creator, :address, :state)`

	for _, m := range messages {
		if _, err := tx.NamedExec(query, toSetWalletStateDatabase(m)); err != nil {
			return err
		}
	}

	return tx.Commit()
}

// GetStates - get states from the "vipcoin_chain_wallets_set_wallet_state" table
func (r Repository) GetStates(filter filter.Filter) ([]*walletstypes.MsgSetWalletState, error) {
	query, args := filter.Build("vipcoin_chain_wallets_set_wallet_state",
		`creator, address, state`)

	var result []types.DBSetWalletState
	if err := r.db.Select(&result, query, args...); err != nil {
		return []*walletstypes.MsgSetWalletState{}, err
	}

	states := make([]*walletstypes.MsgSetWalletState, 0, len(result))
	for _, kind := range result {
		states = append(states, toSetWalletStateDomain(kind))
	}

	return states, nil
}

// UpdateStates - method that updates states in the "vipcoin_chain_wallets_set_wallet_state" table
func (r Repository) UpdateStates(messages ...*walletstypes.MsgSetWalletState) error {
	if len(messages) == 0 {
		return nil
	}

	tx, err := r.db.BeginTxx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	defer tx.Rollback()

	query := `UPDATE vipcoin_chain_wallets_set_wallet_state SET
				 creator = :creator, address = :address, state = :state
				 WHERE address = :address`

	for _, m := range messages {
		stateDB := toSetWalletStateDatabase(m)
		if err != nil {
			return err
		}

		if _, err := tx.NamedExec(query, stateDB); err != nil {
			return err
		}
	}

	return tx.Commit()
}
