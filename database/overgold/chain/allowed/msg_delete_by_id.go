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

// GetAllDeleteByID - method that get data from a db (overgold_allowed_delete_by_id).
func (r Repository) GetAllDeleteByID(filter filter.Filter) ([]allowed.MsgDeleteByID, error) {
	query, args := filter.Build(tableDeleteByID)

	var result []types.AllowedDeleteByID
	if err := r.db.Select(&result, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NotFound{What: tableDeleteByID}
		}

		return nil, errs.Internal{Cause: err.Error()}
	}
	if len(result) == 0 {
		return nil, errs.NotFound{What: tableDeleteByID}
	}

	return toDeleteByIDDomainList(result), nil
}

// InsertToDeleteByID - insert a new MsgDeleteByID in a database (overgold_allowed_delete_by_id).
func (r Repository) InsertToDeleteByID(hash string, msgs ...*allowed.MsgDeleteByID) error {
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
		INSERT INTO overgold_allowed_delete_by_id (
			id, tx_hash, creator
		) VALUES (
			:id, :tx_hash, :creator
		) RETURNING
			id, tx_hash, creator
	`

	for _, m := range msgs {
		if _, err = tx.NamedExec(query, toDeleteByIDDatabase(hash, m)); err != nil {
			return errs.Internal{Cause: err.Error()}
		}
	}

	return tx.Commit()
}
