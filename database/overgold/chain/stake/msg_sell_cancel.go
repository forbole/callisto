package stake

import (
	"context"
	"database/sql"
	"errors"

	"git.ooo.ua/vipcoin/lib/errs"
	"git.ooo.ua/vipcoin/lib/filter"
	stake "git.ooo.ua/vipcoin/ovg-chain/x/stake/types"

	db "github.com/forbole/bdjuno/v4/database/types"
)

// GetAllMsgSellCancel - method that get data from a db (overgold_stake_sell_cancel).
func (r Repository) GetAllMsgSellCancel(filter filter.Filter) ([]stake.MsgMsgCancelSell, error) {
	query, args := filter.Build(tableSellCancel)

	var result []db.StakeMsgSellCancel
	if err := r.db.Select(&result, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NotFound{What: tableSellCancel}
		}

		return nil, errs.Internal{Cause: err.Error()}
	}
	if len(result) == 0 {
		return nil, errs.NotFound{What: tableSellCancel}
	}

	return toMsgSellCancelDomainList(result), nil
}

// InsertMsgSellCancel - insert a new MsgMsgCancelSell in a database (overgold_stake_sell_cancel).
func (r Repository) InsertMsgSellCancel(hash string, msgs ...stake.MsgMsgCancelSell) error {
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
		INSERT INTO overgold_stake_sell_cancel (
			tx_hash, creator, amount
		) VALUES (
			:tx_hash, :creator, :amount
		) RETURNING
			id, tx_hash, creator, amount
	`

	for _, msg := range msgs {
		model, err := toMsgSellCancelDatabase(hash, msg)
		if err != nil {
			return err
		}

		if _, err = tx.NamedExec(query, model); err != nil {
			return errs.Internal{Cause: err.Error()}
		}
	}

	return tx.Commit()
}
