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

// GetAllMsgWithdraw - method that get data from a db (overgold_core_withdraw).
func (r Repository) GetAllMsgWithdraw(filter filter.Filter) ([]core.MsgWithdraw, error) {
	q, args := filter.Build(tableWithdraw)

	var result []db.CoreMsgWithdraw
	if err := r.db.Select(&result, q, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NotFound{What: tableWithdraw}
		}

		return nil, errs.Internal{Cause: err.Error()}
	}
	if len(result) == 0 {
		return nil, errs.NotFound{What: tableWithdraw}
	}

	return toMsgWithdrawDomainList(result), nil
}

// InsertMsgWithdraw - insert a new MsgWithdraw in a database (overgold_core_withdraw).
func (r Repository) InsertMsgWithdraw(hash string, msgs ...core.MsgWithdraw) error {
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

	q := `
		INSERT INTO overgold_core_withdraw (
			tx_hash, creator, amount, denom, address
		) VALUES (
			$1, $2, $3, $4, $5
		) RETURNING
			id, tx_hash, creator, amount, denom, address
	`

	for _, msg := range msgs {
		m, err := toMsgWithdrawDatabase(hash, msg)
		if err != nil {
			return err
		}

		if _, err = tx.Exec(q, m.TxHash, m.Creator, m.Amount, m.Denom, m.Address); err != nil {
			return errs.Internal{Cause: err.Error()}
		}
	}

	return tx.Commit()
}
