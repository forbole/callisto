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

// GetAllUpdateAddresses - method that get data from a db (overgold_allowed_update_addresses).
func (r Repository) GetAllUpdateAddresses(filter filter.Filter) ([]allowed.MsgUpdateAddresses, error) {
	query, args := filter.Build(tableUpdateAddresses)

	var result []types.AllowedUpdateAddresses
	if err := r.db.Select(&result, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NotFound{What: "update_addresses"}
		}

		return nil, errs.Internal{Cause: err.Error()}
	}

	return toUpdateAddressesDomainList(result), nil
}

// InsertToUpdateAddresses - insert a new MsgUpdateAddresses in a database (overgold_allowed_update_addresses).
func (r Repository) InsertToUpdateAddresses(hash string, msgs ...*allowed.MsgUpdateAddresses) error {
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
		INSERT INTO overgold_allowed_update_addresses (
			id, tx_hash, creator, address
		) VALUES (
			:id, :tx_hash, :creator, :address
		) RETURNING
			id, tx_hash, creator, address
	`

	for _, m := range msgs {
		if _, err = tx.NamedExec(query, toUpdateAddressesDatabase(hash, m)); err != nil {
			return errs.Internal{Cause: err.Error()}
		}
	}

	return tx.Commit()
}
