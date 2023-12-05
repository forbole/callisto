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

// GetAllDeleteByAddresses - method that get data from a db (overgold_allowed_delete_by_addresses).
func (r Repository) GetAllDeleteByAddresses(filter filter.Filter) ([]allowed.MsgDeleteByAddresses, error) {
	query, args := filter.Build(tableDeleteByAddresses)

	var result []types.AllowedDeleteByAddresses
	if err := r.db.Select(&result, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NotFound{What: tableDeleteByAddresses}
		}

		return nil, errs.Internal{Cause: err.Error()}
	}
	if len(result) == 0 {
		return nil, errs.NotFound{What: tableDeleteByAddresses}
	}

	return toDeleteByAddressesDomainList(result), nil
}

// InsertToDeleteByAddresses - insert a new MsgCreateAddresses in a database (overgold_allowed_delete_by_addresses).
func (r Repository) InsertToDeleteByAddresses(hash string, msgs ...*allowed.MsgDeleteByAddresses) error {
	if len(msgs) == 0 || hash == "" {
		return nil
	}

	tx, err := r.db.BeginTxx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return errs.Internal{Cause: err.Error()}
	}

	defer func() {
		_ = tx.Rollback()
	}()

	query := `
		INSERT INTO overgold_allowed_delete_by_addresses (
			tx_hash, creator, address
		) VALUES (
			:tx_hash, :creator, :address
		) RETURNING
			id, tx_hash, creator, address
	`

	for _, m := range msgs {
		if _, err = tx.NamedExec(query, toDeleteByAddressesDatabase(hash, m)); err != nil {
			return errs.Internal{Cause: err.Error()}
		}
	}

	return tx.Commit()
}
