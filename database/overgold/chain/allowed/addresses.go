package allowed

import (
	"context"
	"database/sql"
	"errors"

	"git.ooo.ua/vipcoin/lib/errs"
	"git.ooo.ua/vipcoin/lib/filter"
	allowed "git.ooo.ua/vipcoin/ovg-chain/x/allowed/types"

	"github.com/forbole/bdjuno/v4/database/types"
)

// GetAllAddresses - method that get data from a db (overgold_allowed_addresses).
func (r Repository) GetAllAddresses(filter filter.Filter) ([]allowed.Addresses, error) {
	q, args := filter.Build(tableAddresses)

	var result []types.AllowedAddresses
	if err := r.db.Select(&result, q, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NotFound{What: tableAddresses}
		}

		return nil, errs.Internal{Cause: err.Error()}
	}
	if len(result) == 0 {
		return nil, errs.NotFound{What: tableAddresses}
	}

	return toAddressesDomainList(result), nil
}

// InsertToAddresses - insert a new Addresses in a database (overgold_allowed_addresses).
func (r Repository) InsertToAddresses(addresses ...allowed.Addresses) error {
	if len(addresses) == 0 {
		return nil
	}

	tx, err := r.db.BeginTxx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return errs.Internal{Cause: err.Error()}
	}

	defer func() {
		_ = tx.Rollback()
	}()

	q := `
		INSERT INTO overgold_allowed_addresses (
			creator, address
		) VALUES (
			$1, $2
		) RETURNING
			creator, address
	`

	for _, a := range addresses {
		m := toAddressesDatabase(a)
		if _, err = tx.Exec(q, m.Creator, m.Address); err != nil {
			return errs.Internal{Cause: err.Error()}
		}
	}

	return tx.Commit()
}

// UpdateAddresses - method that updates in a database (overgold_allowed_addresses).
func (r Repository) UpdateAddresses(assets ...allowed.Addresses) error {
	if len(assets) == 0 {
		return nil
	}

	tx, err := r.db.BeginTxx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	defer func() {
		_ = tx.Rollback()
	}()

	q := `UPDATE overgold_allowed_addresses SET
				 creator = $1,
				 address = $2
			 WHERE id = $3`

	for _, asset := range assets {
		m := toAddressesDatabase(asset)
		if _, err = tx.Exec(q, m.Creator, m.Address, m.ID); err != nil {
			return err
		}
	}

	return tx.Commit()
}

// DeleteAddressesByAddress - method that deletes data in a database (overgold_allowed_addresses).
func (r Repository) DeleteAddressesByAddress(addresses ...string) error {
	if len(addresses) == 0 {
		return nil
	}

	tx, err := r.db.BeginTxx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return errs.Internal{Cause: err.Error()}
	}

	defer tx.Rollback()

	q := `DELETE FROM overgold_allowed_addresses WHERE address IN ($1)`

	if _, err = tx.Exec(q, deleteAddressesDB{Address: addresses}); err != nil {
		return errs.Internal{Cause: err.Error()}
	}

	return tx.Commit()
}

// DeleteAddressesByID - method that deletes data in a database (overgold_allowed_addresses).
func (r Repository) DeleteAddressesByID(ids ...uint64) error {
	if len(ids) == 0 {
		return nil
	}

	tx, err := r.db.BeginTxx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return errs.Internal{Cause: err.Error()}
	}

	defer tx.Rollback()

	q := `DELETE FROM overgold_allowed_addresses WHERE id = $1`

	for _, id := range ids {
		if _, err = tx.Exec(q, deleteAddressesDB{ID: id}); err != nil {
			return errs.Internal{Cause: err.Error()}
		}
	}

	return tx.Commit()
}
