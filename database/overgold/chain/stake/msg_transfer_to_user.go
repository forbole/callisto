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

// GetAllMsgTransferToUser - method that get data from a db (overgold_stake_transfer_to_user).
func (r Repository) GetAllMsgTransferToUser(filter filter.Filter) ([]stake.MsgTransferToUser, error) {
	query, args := filter.Build(tableTransferToUser)

	var result []db.StakeMsgTransferToUser
	if err := r.db.Select(&result, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NotFound{What: tableTransferToUser}
		}

		return nil, errs.Internal{Cause: err.Error()}
	}
	if len(result) == 0 {
		return nil, errs.NotFound{What: tableTransferToUser}
	}

	return toMsgTransferToUserDomainList(result), nil
}

// InsertMsgTransferToUser - insert a new MsgTransferToUser in a database (overgold_stake_transfer_to_user).
func (r Repository) InsertMsgTransferToUser(hash string, msgs ...stake.MsgTransferToUser) error {
	if len(msgs) == 0 || hash == "" {
		return nil
	}

	q := `
		INSERT INTO overgold_stake_transfer_to_user ( tx_hash, creator, amount, address ) 
		VALUES ( $1, $2, $3, $4 )
	`

	for _, msg := range msgs {
		m, err := toMsgTransferToUserDatabase(hash, msg)
		if err != nil {
			return err
		}

		if _, err := r.db.Exec(q, m.TxHash, m.Creator, m.Amount, m.Address); err != nil {
			if chain.IsAlreadyExists(err) {
				continue
			}
			return errs.Internal{Cause: err.Error()}
		}
	}

	return nil
}
