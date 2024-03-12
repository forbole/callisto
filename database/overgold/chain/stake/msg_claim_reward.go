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

// InsertMsgClaimReward - insert a new ClaimReward in a database (overgold_stake_claim_reward).
func (r Repository) InsertMsgClaimReward(hash string, msgs ...stake.MsgClaimReward) error {
	if len(msgs) == 0 || hash == "" {
		return nil
	}

	q := `
		INSERT INTO overgold_stake_claim_reward (
			tx_hash, creator, amount
		) VALUES (
			$1, $2, $3
		) RETURNING 
			id, tx_hash, creator, amount
	`

	for _, msg := range msgs {
		m := toMsgClaimRewardDatabase(hash, msg)

		if _, err := r.db.Exec(q, m.TxHash, m.Creator, m.Amount); err != nil {
			if chain.IsAlreadyExists(err) {
				continue
			}
			return errs.Internal{Cause: err.Error()}
		}
	}

	return nil
}
