package referral

import (
	"context"
	"database/sql"
	"errors"

	"git.ooo.ua/vipcoin/lib/errs"
	"git.ooo.ua/vipcoin/lib/filter"
	"git.ooo.ua/vipcoin/ovg-chain/x/referral/types"

	db "github.com/forbole/bdjuno/v4/database/types"
)

// GetAllMsgSetReferrer - method that get data from a db (overgold_referral_set_referrer).
func (r Repository) GetAllMsgSetReferrer(filter filter.Filter) ([]types.MsgSetReferrer, error) {
	q, args := filter.Build(tableSetReferrer)

	var result []db.DbReferralSetReferrer
	if err := r.db.Select(&result, q, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NotFound{What: tableSetReferrer}
		}

		return nil, errs.Internal{Cause: err.Error()}
	}
	if len(result) == 0 {
		return nil, errs.NotFound{What: tableSetReferrer}
	}

	return toMsgSetReferrerDomainList(result), nil
}

// InsertMsgSetReferrer - insert a new MsgIssue in a database (overgold_referral_set_referrer).
func (r Repository) InsertMsgSetReferrer(hash string, msgs ...types.MsgSetReferrer) error {
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
		INSERT INTO overgold_referral_set_referrer (
			tx_hash, creator, referrer_address, referral_address
		) VALUES (
			$1, $2, $3, $4
		) RETURNING
			id, tx_hash, creator, referrer_address, referral_address
	`

	for _, msg := range msgs {
		m := toMsgSetReferrerDatabase(hash, msg)
		if _, err = tx.Exec(q, m.TxHash, m.Creator, m.ReferrerAddress, m.ReferralAddress); err != nil {
			return errs.Internal{Cause: err.Error()}
		}
	}

	return tx.Commit()
}
