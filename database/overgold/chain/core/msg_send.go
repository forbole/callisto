package core

import (
	"context"
	"database/sql"
	"errors"

	"git.ooo.ua/vipcoin/lib/errs"
	"git.ooo.ua/vipcoin/lib/filter"
	core "git.ooo.ua/vipcoin/ovg-chain/x/core/types"

	db "github.com/forbole/bdjuno/v4/database/types"
)

// GetAllMsgSend - method that get data from a db (overgold_core_send).
func (r Repository) GetAllMsgSend(filter filter.Filter) ([]core.MsgSend, error) {
	query, args := filter.Build(tableSend)

	var result []db.CoreMsgSend
	if err := r.db.Select(&result, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NotFound{What: tableSend}
		}

		return nil, errs.Internal{Cause: err.Error()}
	}
	if len(result) == 0 {
		return nil, errs.NotFound{What: tableSend}
	}

	return toMsgSendDomainList(result), nil
}

// InsertMsgSend - insert a new MsgCreateAddresses in a database (overgold_core_send).
func (r Repository) InsertMsgSend(hash string, msgs ...core.MsgSend) error {
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
		INSERT INTO overgold_core_send (
			tx_hash, creator, amount, denom, address_from, address_to
		) VALUES (
			:tx_hash, :creator, :amount, :denom, :address_from, :address_to
		) RETURNING
			id, tx_hash, creator, amount, denom, address_from, address_to
	`

	for _, msg := range msgs {
		model, err := toMsgSendDatabase(hash, msg)
		if err != nil {
			return err
		}

		if _, err = tx.NamedExec(query, model); err != nil {
			return errs.Internal{Cause: err.Error()}
		}
	}

	return tx.Commit()
}
