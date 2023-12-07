package bank

import (
	"context"
	"database/sql"
	"errors"

	"git.ooo.ua/vipcoin/lib/errs"
	"git.ooo.ua/vipcoin/lib/filter"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/lib/pq"
)

// GetAllMsgSend - method that get data from a db (msg_send).
func (r Repository) GetAllMsgSend(filter filter.Filter) ([]bank.MsgSend, error) {
	query, args := filter.Build(tableMsgSend)

	var result []msgSend
	if err := r.db.Select(&result, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NotFound{What: tableMsgSend}
		}

		return nil, errs.Internal{Cause: err.Error()}
	}
	if len(result) == 0 {
		return nil, errs.NotFound{What: tableMsgSend}
	}

	return toGetMsgSendDomainList(result)
}

// InsertMsgSend - insert a new MsgCreateAddresses in a database (msg_send).
func (r Repository) InsertMsgSend(hash string, msgs ...bank.MsgSend) error {
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

	q := `INSERT INTO msg_send (
	    	tx_hash, from_address, to_address, amount
	    ) VALUES (
	    	$1, $2, $3, $4
	    	) RETURNING
			id, tx_hash, from_address, to_address, amount
	`

	// NOTE: use tx.Exec for custom type pq.Array(DbCoins)
	for _, msg := range msgs {
		m := toMsgSendDatabase(hash, msg)
		if _, err = tx.Exec(q, m.TxHash, m.FromAddress, m.ToAddress, pq.Array(m.Amount)); err != nil {
			return errs.Internal{Cause: err.Error()}
		}
	}

	return tx.Commit()
}