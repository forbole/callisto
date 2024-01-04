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

// GetAllMsgTransferFromUser - method that get data from a db (overgold_stake_transfer_from_user).
func (r Repository) GetAllMsgTransferFromUser(filter filter.Filter) ([]stake.MsgTransferFromUser, error) {
	query, args := filter.Build(tableTransferFromUser)

	var result []db.StakeMsgTransferFromUser
	if err := r.db.Select(&result, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NotFound{What: tableTransferFromUser}
		}

		return nil, errs.Internal{Cause: err.Error()}
	}
	if len(result) == 0 {
		return nil, errs.NotFound{What: tableTransferFromUser}
	}

	return toMsgTransferFromUserDomainList(result), nil
}

// InsertMsgTransferFromUser - insert a new MsgTransferFromUser in a database (overgold_stake_transfer_from_user).
func (r Repository) InsertMsgTransferFromUser(hash string, msgs ...stake.MsgTransferFromUser) error {
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
		INSERT INTO overgold_stake_transfer_from_user ( tx_hash, creator, amount, address) 
		VALUES ( $1, $2, $3, $4 )
	`

	for _, msg := range msgs {
		m, err := toMsgTransferFromUserDatabase(hash, msg)
		if err != nil {
			return err
		}

		if _, err = tx.Exec(q, m.TxHash, m.Creator, m.Amount, m.Address); err != nil {
			return errs.Internal{Cause: err.Error()}
		}
	}

	return tx.Commit()
}
