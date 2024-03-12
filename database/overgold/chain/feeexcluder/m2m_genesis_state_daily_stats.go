package feeexcluder

import (
	"database/sql"
	"errors"

	"git.ooo.ua/vipcoin/lib/errs"
	"git.ooo.ua/vipcoin/lib/filter"
	"github.com/jmoiron/sqlx"

	"github.com/forbole/bdjuno/v4/database/overgold/chain"
	"github.com/forbole/bdjuno/v4/database/types"
)

// GetAllM2MGenesisStateDailyStats - method that get data from a db (overgold_feeexcluder_m2m_genesis_state_daily_stats).
func (r Repository) GetAllM2MGenesisStateDailyStats(filter filter.Filter) ([]types.FeeExcluderM2MGenesisStateDailyStats, error) {
	q, args := filter.Build(tableM2MGenesisStateDailyStats)

	var result []types.FeeExcluderM2MGenesisStateDailyStats
	if err := r.db.Select(&result, q, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NotFound{What: tableM2MGenesisStateDailyStats}
		}

		return nil, errs.Internal{Cause: err.Error()}
	}
	if len(result) == 0 {
		return nil, errs.NotFound{What: tableM2MGenesisStateDailyStats}
	}

	return result, nil
}

// InsertToM2MGenesisStateDailyStats - insert new data in a database (overgold_feeexcluder_m2m_genesis_state_daily_stats).
func (r Repository) InsertToM2MGenesisStateDailyStats(_ *sqlx.Tx, ids ...types.FeeExcluderM2MGenesisStateDailyStats) (err error) {
	if len(ids) == 0 {
		return nil
	}
	q := `
		INSERT INTO overgold_feeexcluder_m2m_genesis_state_daily_stats (
			genesis_state_id, daily_stats_id
		) VALUES (
			$1, $2
		) RETURNING
			genesis_state_id, daily_stats_id
	`

	for _, m := range ids {
		if _, err = r.db.Exec(q, m.GenesisStateID, m.DailyStatsID); err != nil {
			if chain.IsAlreadyExists(err) {
				continue
			}
			return errs.Internal{Cause: err.Error()}
		}
	}

	return nil
}

// DeleteM2MGenesisStateDailyStatsByGenesisState - method that deletes data in a database (overgold_feeexcluder_m2m_genesis_state_daily_stats).
func (r Repository) DeleteM2MGenesisStateDailyStatsByGenesisState(tx *sqlx.Tx, id uint64) (err error) {
	q := `DELETE FROM overgold_feeexcluder_m2m_genesis_state_daily_stats WHERE genesis_state_id IN ($1)`

	if _, err = r.db.Exec(q, id); err != nil {
		return errs.Internal{Cause: err.Error()}
	}

	return nil
}
