package stake

import (
	"database/sql"
	"errors"

	"git.ooo.ua/vipcoin/lib/errs"
	"git.ooo.ua/vipcoin/lib/filter"
	stake "git.ooo.ua/vipcoin/ovg-chain/x/stake/types"

	"github.com/forbole/bdjuno/v4/database/overgold/chain"
	db "github.com/forbole/bdjuno/v4/database/types"
)

// GetAllMsgSellCancel - method that get data from a db (overgold_stake_sell_cancel).
func (r Repository) GetAllMsgSellCancel(filter filter.Filter) ([]stake.MsgMsgCancelSell, error) {
	q, args := filter.Build(tableSellCancel)

	var result []db.StakeMsgSellCancel
	if err := r.db.Select(&result, q, args...); err != nil {
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

	q := `
		INSERT INTO overgold_stake_sell_cancel (
			tx_hash, creator, amount
		) VALUES (
			$1, $2, $3
		) RETURNING
			id, tx_hash, creator, amount
	`

	for _, msg := range msgs {
		m, err := toMsgSellCancelDatabase(hash, msg)
		if err != nil {
			return err
		}

		if _, err := r.db.Exec(q, m.TxHash, m.Creator, m.Amount); err != nil {
			if chain.IsAlreadyExists(err) {
				continue
			}
			return errs.Internal{Cause: err.Error()}
		}
	}

	return nil
}
