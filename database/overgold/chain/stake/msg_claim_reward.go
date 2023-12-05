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

// GetAllMsgClaimReward - method that get data from a db (overgold_stake_claim_reward).
func (r Repository) GetAllMsgClaimReward(filter filter.Filter) ([]stake.MsgClaimReward, error) {
	query, args := filter.Build(tableClaimReward)

	var result []db.StakeMsgClaim
	if err := r.db.Select(&result, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NotFound{What: tableClaimReward}
		}

		return nil, errs.Internal{Cause: err.Error()}
	}
	if len(result) == 0 {
		return nil, errs.NotFound{What: tableClaimReward}
	}

	return toMsgClaimRewardDomainList(result), nil
}

// InsertClaimReward - insert a new ClaimReward in a database (overgold_stake_claim_reward).
func (r Repository) InsertMsgClaimReward(hash string, msgs ...stake.MsgClaimReward) error {
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
		INSERT INTO overgold_stake_claim_reward (
			tx_hash, creator, amount
		) VALUES (
			:tx_hash, :creator, :amount
		) RETURNING
			id, tx_hash, creator, amount
	`

	for _, msg := range msgs {
		model := toMsgClaimRewardDatabase(hash, msg)

		if _, err = tx.NamedExec(query, model); err != nil {
			return errs.Internal{Cause: err.Error()}
		}
	}

	return tx.Commit()
}
