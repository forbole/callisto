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

// GetAllMsgDistributeRewards - method that get data from a db (overgold_stake_distribute_rewards).
func (r Repository) GetAllMsgDistributeRewards(filter filter.Filter) ([]stake.MsgDistributeRewards, error) {
	query, args := filter.Build(tableDistributeRewards)

	var result []db.StakeMsgDistribute
	if err := r.db.Select(&result, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NotFound{What: tableDistributeRewards}
		}

		return nil, errs.Internal{Cause: err.Error()}
	}
	if len(result) == 0 {
		return nil, errs.NotFound{What: tableDistributeRewards}
	}

	return toMsgDistributeDomainList(result), nil
}

// InsertMsgDistributeRewards - insert a new MsgDistributeRewards in a database (overgold_stake_distribute_rewards).
func (r Repository) InsertMsgDistributeRewards(hash string, msgs ...stake.MsgDistributeRewards) error {
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
		INSERT INTO overgold_stake_distribute_rewards (
			tx_hash, creator
		) VALUES (
			$1, $2
		) RETURNING
			id, tx_hash, creator
	`

	for _, msg := range msgs {
		m := toMsgDistributeDatabase(hash, msg)

		if _, err = tx.Exec(query, m.TxHash, m.Creator); err != nil {
			return errs.Internal{Cause: err.Error()}
		}
	}

	return tx.Commit()
}