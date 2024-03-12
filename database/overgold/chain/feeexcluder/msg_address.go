package feeexcluder

import (
	"database/sql"
	"errors"

	"git.ooo.ua/vipcoin/lib/errs"
	"git.ooo.ua/vipcoin/lib/filter"
	fe "git.ooo.ua/vipcoin/ovg-chain/x/feeexcluder/types"
	"github.com/jmoiron/sqlx"

	"github.com/forbole/bdjuno/v4/database/overgold/chain"
	"github.com/forbole/bdjuno/v4/database/types"
)

// GetAllAddress - method that get data from a db (overgold_feeexcluder_address).
func (r Repository) GetAllAddress(filter filter.Filter) ([]fe.Address, error) {
	q, args := filter.Build(tableAddress)

	var result []types.FeeExcluderAddress
	if err := r.db.Select(&result, q, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NotFound{What: tableAddress}
		}

		return nil, errs.Internal{Cause: err.Error()}
	}
	if len(result) == 0 {
		return nil, errs.NotFound{What: tableAddress}
	}

	return toAddressDomainList(result), nil
}

// InsertToAddress - insert new data in a database (overgold_feeexcluder_address).
func (r Repository) InsertToAddress(_ *sqlx.Tx, address fe.Address) (lastID uint64, err error) {
	q := `
		INSERT INTO overgold_feeexcluder_address (
			msg_id, creator, address
		) VALUES (
			$1, $2, $3
		) RETURNING	id
	`

	m := toAddressDatabase(0, address)
	if err = r.db.QueryRowx(q, m.MsgID, m.Creator, m.Address).Scan(&lastID); err != nil {
		if chain.IsAlreadyExists(err) {
			return 0, nil
		}
		return 0, errs.Internal{Cause: err.Error()}
	}

	return lastID, nil
}

// UpdateAddress - method that updates in a database (overgold_feeexcluder_address).
func (r Repository) UpdateAddress(_ *sqlx.Tx, id uint64, address fe.Address) (err error) {
	q := `UPDATE overgold_feeexcluder_address SET
				 msg_id = $1,
				 creator = $2,
				 address = $3
			 WHERE id = $4`

	m := toAddressDatabase(id, address)
	if _, err = r.db.Exec(q, m.MsgID, m.Creator, m.Address, m.ID); err != nil {
		return err
	}

	return nil
}

// DeleteAddress - method that deletes data in a database (overgold_feeexcluder_address).
func (r Repository) DeleteAddress(_ *sqlx.Tx, id uint64) (err error) {
	q := `DELETE FROM overgold_feeexcluder_address WHERE id IN ($1)`

	if _, err = r.db.Exec(q, id); err != nil {
		return errs.Internal{Cause: err.Error()}
	}

	return nil
}
