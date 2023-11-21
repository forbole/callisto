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

// GetAllCreateAddresses - method that get data from a db (overgold_allowed_create_addresses).
func (r Repository) GetAllCreateAddresses(filter filter.Filter) ([]allowed.MsgCreateAddresses, error) {
	query, args := filter.Build(tableCreateAddresses)

	var result []types.AllowedCreateAddresses
	if err := r.db.Select(&result, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NotFound{What: "create_addresses"}
		}

		return nil, errs.Internal{Cause: err.Error()}
	}

	return toCreateAddressesDomainList(result), nil
}

// InsertToCreateAddresses - insert a new MsgCreateAddresses in a database (overgold_allowed_create_addresses).
func (r Repository) InsertToCreateAddresses(hash string, msgs ...*allowed.MsgCreateAddresses) error {
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
		INSERT INTO overgold_allowed_create_addresses (
			tx_hash, creator, address
		) VALUES (
			:tx_hash, :creator, :address
		) RETURNING
			id, tx_hash, creator, address
	`

	for _, m := range msgs {
		if _, err = tx.NamedExec(query, toCreateAddressesDatabase(hash, m)); err != nil {
			return errs.Internal{Cause: err.Error()}
		}
	}

	return tx.Commit()
}
