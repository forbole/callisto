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

// GetAllMsgIssue - method that get data from a db (overgold_core_issue).
func (r Repository) GetAllMsgIssue(filter filter.Filter) ([]core.MsgIssue, error) {
	query, args := filter.Build(tableIssue)

	var result []db.CoreMsgIssue
	if err := r.db.Select(&result, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NotFound{What: tableIssue}
		}

		return nil, errs.Internal{Cause: err.Error()}
	}
	if len(result) == 0 {
		return nil, errs.NotFound{What: tableIssue}
	}

	return toMsgIssueDomainList(result), nil
}

// InsertMsgIssue - insert a new MsgIssue in a database (overgold_core_issue).
func (r Repository) InsertMsgIssue(hash string, msgs ...core.MsgIssue) error {
	if len(msgs) == 0 || hash == "" {
		return nil
	}

	q := `
		INSERT INTO overgold_core_issue (
			tx_hash, creator, amount, denom, address
		) VALUES (
			$1, $2, $3, $4, $5
		) RETURNING
			id, tx_hash, creator, amount, denom, address
	`

	for _, msg := range msgs {
		m, err := toMsgIssueDatabase(hash, msg)
		if err != nil {
			return err
		}

		if _, err := r.db.Exec(q, m.TxHash, m.Creator, m.Amount, m.Denom, m.Address); err != nil {
			if chain.IsAlreadyExists(err) {
				continue
			}
			return errs.Internal{Cause: err.Error()}
		}
	}

	return nil
}
