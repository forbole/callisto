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

// GetAllMsgSell - method that get data from a db (overgold_stake_sell).
func (r Repository) GetAllMsgSell(filter filter.Filter) ([]stake.MsgSellRequest, error) {
	query, args := filter.Build(tableSell)

	var result []db.StakeMsgSell
	if err := r.db.Select(&result, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NotFound{What: tableSell}
		}

		return nil, errs.Internal{Cause: err.Error()}
	}
	if len(result) == 0 {
		return nil, errs.NotFound{What: tableSell}
	}

	return toMsgSellDomainList(result), nil
}

// InsertMsgSell - insert a new MsgSellRequest in a database (overgold_stake_sell).
func (r Repository) InsertMsgSell(hash string, msgs ...stake.MsgSellRequest) error {
	if len(msgs) == 0 || hash == "" {
		return nil
	}

	q := `
		INSERT INTO overgold_stake_sell (
			tx_hash, creator, amount
		) VALUES (
			$1, $2, $3
		) RETURNING
			id, tx_hash, creator, amount
	`

	for _, msg := range msgs {
		m, err := toMsgSellDatabase(hash, msg)
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
