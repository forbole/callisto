package feeexcluder

import (
	"context"
	"database/sql"
	"errors"

	"git.ooo.ua/vipcoin/lib/errs"
	"git.ooo.ua/vipcoin/lib/filter"
	fe "git.ooo.ua/vipcoin/ovg-chain/x/feeexcluder/types"
	"github.com/jmoiron/sqlx"

	"github.com/forbole/bdjuno/v4/database/overgold/chain"
	"github.com/forbole/bdjuno/v4/database/types"
)

// GetAllStats - method that get data from a db (overgold_feeexcluder_stats).
// TODO: use JOIN and other db model
func (r Repository) GetAllStats(f filter.Filter) ([]fe.Stats, error) {
	q, args := f.Build(tableStats)

	// 1) get stats
	var stats []types.FeeExcluderStats
	if err := r.db.Select(&stats, q, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NotFound{What: tableStats}
		}

		return nil, errs.Internal{Cause: err.Error()}
	}
	if len(stats) == 0 {
		return nil, errs.NotFound{What: tableStats}
	}

	// 2) get daily stats
	result := make([]fe.Stats, 0, len(stats))
	for _, s := range stats {
		dailyStats, err := r.GetAllDailyStats(filter.NewFilter().SetArgument(types.FieldID, s.DailyStatsID))
		if err != nil {
			return nil, err
		}
		if len(dailyStats) == 0 {
			return nil, errs.NotFound{What: tableDailyStats}
		}

		result = append(result, toStatsDomain(&dailyStats[0], s))
	}

	return result, nil
}

// InsertToStats - insert new data in a database (overgold_feeexcluder_stats).
func (r Repository) InsertToStats(_ *sqlx.Tx, stats fe.Stats) (lastID string, err error) {
	// 1) add daily stats and get unique ids
	dailyStatsID, err := r.InsertToDailyStats(nil, *stats.Stats)
	if err != nil {
		return "", err
	}

	// 2) add stats
	q := `
		INSERT INTO overgold_feeexcluder_stats (
			id, date, daily_stats_id
		) VALUES (
			$1, $2, $3
		) RETURNING	id
	`

	m, err := toStatsDatabase(dailyStatsID, stats)
	if err != nil {
		return "", errs.Internal{Cause: err.Error()}
	}

	if err = r.db.QueryRowx(q, m.ID, m.Date, m.DailyStatsID).Scan(&lastID); err != nil {
		if chain.IsAlreadyExists(err) {
			return "", nil
		}
		return "", errs.Internal{Cause: err.Error()}
	}

	return lastID, nil
}

// UpdateStats - method that updates in a database (overgold_feeexcluder_stats).
func (r Repository) UpdateStats(_ *sqlx.Tx, stats fe.Stats) (err error) {
	// 1) update stats and get unique id for daily stats
	// 1.a) get daily stats id via stats index
	s, err := r.getStatsWithUniqueID(filter.NewFilter().SetArgument(types.FieldID, stats.Index))
	if err != nil {
		return err
	}

	// 1.b) update stats
	q := `UPDATE overgold_feeexcluder_stats SET
				 date = $1,
				 daily_stats_id = $2
			 WHERE id = $3`

	m, err := toStatsDatabase(s.DailyStatsID, stats)
	if err != nil {
		return errs.Internal{Cause: err.Error()}
	}

	if _, err = r.db.Exec(q, m.Date, m.DailyStatsID, m.ID); err != nil {
		return err
	}

	// 2) update daily stats
	if err = r.UpdateDailyStats(nil, s.DailyStatsID, *stats.Stats); err != nil {
		return err
	}

	return nil
}

// DeleteStats - method that deletes data in a database (overgold_feeexcluder_stats).
func (r Repository) DeleteStats(_ *sqlx.Tx, id string) (err error) {
	// 1) delete stats and get unique id via stats index
	// 1.a) get stats id via stats index
	s, err := r.getStatsWithUniqueID(filter.NewFilter().SetArgument(types.FieldID, id))
	if err != nil {
		return err
	}

	// 1.b) delete daily stats
	q := `DELETE FROM overgold_feeexcluder_stats WHERE id IN ($1)`

	if _, err = r.db.Exec(q, id); err != nil {
		return errs.Internal{Cause: err.Error()}
	}

	// 2) delete daily stats
	if err = r.DeleteDailyStats(nil, s.DailyStatsID); err != nil {
		return err
	}

	return nil
}

// getStatsWithUniqueID - method that get data from a db (overgold_feeexcluder_stats).
func (r Repository) getStatsWithUniqueID(req filter.Filter) (types.FeeExcluderStats, error) {
	query, args := req.SetLimit(1).Build(tableStats)

	var result types.FeeExcluderStats
	if err := r.db.GetContext(context.Background(), &result, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return types.FeeExcluderStats{}, errs.NotFound{What: tableStats}
		}

		return types.FeeExcluderStats{}, errs.Internal{Cause: err.Error()}
	}

	return result, nil
}

// getAllStatsWithUniqueID - method that get data from a db (overgold_feeexcluder_stats).
func (r Repository) getAllStatsWithUniqueID(f filter.Filter) ([]types.FeeExcluderStats, error) {
	q, args := f.Build(tableStats)

	var stats []types.FeeExcluderStats
	if err := r.db.Select(&stats, q, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NotFound{What: tableStats}
		}

		return nil, errs.Internal{Cause: err.Error()}
	}
	if len(stats) == 0 {
		return nil, errs.NotFound{What: tableStats}
	}

	return stats, nil
}
