package feeexcluder

import (
	"context"
	"database/sql"
	"errors"

	"git.ooo.ua/vipcoin/lib/errs"
	"git.ooo.ua/vipcoin/lib/filter"
	fe "git.ooo.ua/vipcoin/ovg-chain/x/feeexcluder/types"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// GetAllDailyStats - method that get data from a db (overgold_feeexcluder_daily_stats).
func (r Repository) GetAllDailyStats(f filter.Filter) ([]fe.DailyStats, error) {
	q, args := f.Build(tableDailyStats)

	var ds dailyStatsList
	if err := r.db.Select(&ds, q, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NotFound{What: tableDailyStats}
		}

		return nil, errs.Internal{Cause: err.Error()}
	}
	if len(ds) == 0 {
		return nil, errs.NotFound{What: tableDailyStats}
	}

	return ds.toDomain()
}

// InsertToDailyStats - insert new data in a database (overgold_feeexcluder_daily_stats).
func (r Repository) InsertToDailyStats(tx *sqlx.Tx, dailyStats fe.DailyStats) (lastID uint64, err error) {
	if tx == nil {
		tx, err = r.db.BeginTxx(context.Background(), &sql.TxOptions{})
		if err != nil {
			return 0, errs.Internal{Cause: err.Error()}
		}

		defer commit(tx, err)
	}

	q := `
		INSERT INTO overgold_feeexcluder_daily_stats (
			msg_id, amount_with_fee, amount_no_fee, fee, count_with_fee, count_no_fee
		) VALUES (
			$1, $2, $3, $4, $5, $6
		) RETURNING	id
	`

	m := toDailyStatsDatabase(0, dailyStats)
	if err = tx.QueryRowx(q,
		m.MsgID,
		pq.Array(m.AmountWithFee),
		pq.Array(m.AmountNoFee),
		pq.Array(m.Fee),
		m.CountWithFee,
		m.CountNoFee,
	).Scan(&lastID); err != nil {
		return 0, errs.Internal{Cause: err.Error()}
	}

	return lastID, nil
}

// UpdateDailyStats - method that deletes in a database (overgold_feeexcluder_daily_stats).
func (r Repository) UpdateDailyStats(tx *sqlx.Tx, id uint64, ut fe.DailyStats) (err error) {
	if tx == nil {
		tx, err = r.db.BeginTxx(context.Background(), &sql.TxOptions{})
		if err != nil {
			return errs.Internal{Cause: err.Error()}
		}

		defer commit(tx, err)
	}

	q := `UPDATE overgold_feeexcluder_daily_stats SET
                 msg_id = $1,
				 amount_with_fee = $2,
				 amount_no_fee = $3,
            	 fee = $4,
            	 count_with_fee = $5,
                 count_no_fee = $6
			 WHERE id = $7`

	m := toDailyStatsDatabase(id, ut)
	if _, err := tx.Exec(q,
		m.MsgID,
		pq.Array(m.AmountWithFee),
		pq.Array(m.AmountNoFee),
		pq.Array(m.Fee),
		m.CountWithFee,
		m.CountNoFee,
		m.ID,
	); err != nil {
		return errs.Internal{Cause: err.Error()}
	}

	return nil
}

// DeleteDailyStats - method that deletes data in a database (overgold_feeexcluder_daily_stats).
func (r Repository) DeleteDailyStats(tx *sqlx.Tx, id uint64) (err error) {
	if tx == nil {
		tx, err = r.db.BeginTxx(context.Background(), &sql.TxOptions{})
		if err != nil {
			return errs.Internal{Cause: err.Error()}
		}

		defer commit(tx, err)
	}

	q := `DELETE FROM overgold_feeexcluder_daily_stats WHERE id IN ($1)`

	if _, err = tx.Exec(q, id); err != nil {
		return errs.Internal{Cause: err.Error()}
	}

	return nil
}
