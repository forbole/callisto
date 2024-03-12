package core

import (
	"database/sql"
	"errors"

	"git.ooo.ua/vipcoin/lib/errs"
	"git.ooo.ua/vipcoin/lib/filter"
	core "git.ooo.ua/vipcoin/ovg-chain/x/core/types"

	"github.com/forbole/bdjuno/v4/database/overgold/chain"
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

// InsertMsgSend - insert a new MsgSend in a database (overgold_core_send).
func (r Repository) InsertMsgSend(hash string, msgs ...core.MsgSend) error {
	if len(msgs) == 0 || hash == "" {
		return nil
	}

	q := `
		INSERT INTO overgold_core_send (
			tx_hash, creator, amount, denom, address_from, address_to
		) VALUES (
			$1, $2, $3, $4, $5, $6
		) RETURNING
			id, tx_hash, creator, amount, denom, address_from, address_to
	`

	for _, msg := range msgs {
		m, err := toMsgSendDatabase(hash, msg)
		if err != nil {
			return err
		}

		if _, err := r.db.Exec(q, m.TxHash, m.Creator, m.Amount, m.Denom, m.AddressFrom, m.AddressTo); err != nil {
			if chain.IsAlreadyExists(err) {
				continue
			}
			return errs.Internal{Cause: err.Error()}
		}
	}

	return nil
}
