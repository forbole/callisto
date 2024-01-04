package feeexcluder

import (
	"context"
	"database/sql"
	"errors"

	"git.ooo.ua/vipcoin/lib/errs"
	"git.ooo.ua/vipcoin/lib/filter"
	"github.com/jmoiron/sqlx"

	"github.com/forbole/bdjuno/v4/database/types"
)

// GetAllM2MGenesisStateAddress - method that get data from a db (overgold_feeexcluder_m2m_genesis_state_address).
func (r Repository) GetAllM2MGenesisStateAddress(filter filter.Filter) ([]types.FeeExcluderM2MGenesisStateAddress, error) {
	q, args := filter.Build(tableM2MGenesisStateAddress)

	var result []types.FeeExcluderM2MGenesisStateAddress
	if err := r.db.Select(&result, q, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NotFound{What: tableM2MGenesisStateAddress}
		}

		return nil, errs.Internal{Cause: err.Error()}
	}
	if len(result) == 0 {
		return nil, errs.NotFound{What: tableM2MGenesisStateAddress}
	}

	return result, nil
}

// InsertToM2MGenesisStateAddress - insert new data in a database (overgold_feeexcluder_m2m_genesis_state_address).
func (r Repository) InsertToM2MGenesisStateAddress(tx *sqlx.Tx, ids ...types.FeeExcluderM2MGenesisStateAddress) (err error) {
	if len(ids) == 0 {
		return nil
	}

	if tx == nil {
		tx, err = r.db.BeginTxx(context.Background(), &sql.TxOptions{})
		if err != nil {
			return errs.Internal{Cause: err.Error()}
		}

		defer commit(tx, err)
	}

	q := `
		INSERT INTO overgold_feeexcluder_m2m_genesis_state_address (
			genesis_state_id, address_id
		) VALUES (
			$1, $2
		) RETURNING
			genesis_state_id, address_id
	`

	for _, m := range ids {
		if _, err = tx.Exec(q, m.GenesisStateID, m.AddressID); err != nil {
			return errs.Internal{Cause: err.Error()}
		}
	}

	return nil
}

// DeleteM2MGenesisStateAddressByGenesisState - method that deletes data in a database (overgold_feeexcluder_m2m_genesis_state_address).
func (r Repository) DeleteM2MGenesisStateAddressByGenesisState(tx *sqlx.Tx, id uint64) (err error) {
	if tx == nil {
		tx, err = r.db.BeginTxx(context.Background(), &sql.TxOptions{})
		if err != nil {
			return errs.Internal{Cause: err.Error()}
		}

		defer commit(tx, err)
	}

	q := `DELETE FROM overgold_feeexcluder_m2m_genesis_state_address WHERE genesis_state_id IN ($1)`

	if _, err = tx.Exec(q, id); err != nil {
		return errs.Internal{Cause: err.Error()}
	}

	return nil
}
