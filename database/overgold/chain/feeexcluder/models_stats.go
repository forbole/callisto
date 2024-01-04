package feeexcluder

import (
	"time"

	"git.ooo.ua/vipcoin/ovg-chain/x/feeexcluder/types"
	"github.com/lib/pq"

	fe "git.ooo.ua/vipcoin/ovg-chain/x/feeexcluder/types"

	db "github.com/forbole/bdjuno/v4/database/types"
)

// BLOCK DailyStats

type (
	// dailyStats represents a single row inside the 'overgold_feeexcluder_daily_stats' table (used only for SELECT)
	dailyStats struct {
		CountWithFee  int32          `db:"count_with_fee"`
		CountNoFee    int32          `db:"count_no_fee"`
		ID            uint64         `db:"id"`
		MsgID         uint64         `db:"msg_id"`
		AmountWithFee pq.StringArray `db:"amount_with_fee"`
		AmountNoFee   pq.StringArray `db:"amount_no_fee"`
		Fee           pq.StringArray `db:"fee"`
	}

	dailyStatsList []dailyStats
)

// toDomain - mapping func to a domain.
func (s dailyStats) toDomain() (fe.DailyStats, error) {
	amountWithFee, err := db.FromPqStringArrayToCoins(s.AmountWithFee)
	if err != nil {
		return fe.DailyStats{}, err
	}

	amountNoFee, err := db.FromPqStringArrayToCoins(s.AmountNoFee)
	if err != nil {
		return fe.DailyStats{}, err
	}

	fee, err := db.FromPqStringArrayToCoins(s.Fee)
	if err != nil {
		return fe.DailyStats{}, err
	}

	return fe.DailyStats{
		Id:            s.MsgID,
		CountWithFee:  s.CountWithFee,
		CountNoFee:    s.CountNoFee,
		AmountWithFee: amountWithFee,
		AmountNoFee:   amountNoFee,
		Fee:           fee,
	}, nil
}

// toDomain - mapping func to a domain list.
func (list dailyStatsList) toDomain() ([]fe.DailyStats, error) {
	res := make([]fe.DailyStats, 0, len(list))
	for _, ds := range list {
		d, err := ds.toDomain()
		if err != nil {
			return nil, err
		}

		res = append(res, d)
	}

	return res, nil
}

// toDailyStatsDomain - mapping func to a domain model.
func toDailyStatsDomain(s db.FeeExcluderDailyStats) types.DailyStats {
	return types.DailyStats{
		Id:            s.MsgID,
		AmountWithFee: s.AmountWithFee.ToCoins(),
		AmountNoFee:   s.AmountNoFee.ToCoins(),
		Fee:           s.Fee.ToCoins(),
		CountWithFee:  s.CountWithFee,
		CountNoFee:    s.CountNoFee,
	}
}

// toDailyStatsDatabase - mapping func to a database model.
func toDailyStatsDatabase(id uint64, s types.DailyStats) db.FeeExcluderDailyStats {
	return db.FeeExcluderDailyStats{
		ID:            id,
		MsgID:         s.Id,
		CountWithFee:  s.CountWithFee,
		CountNoFee:    s.CountNoFee,
		AmountWithFee: db.NewDbCoins(s.AmountWithFee),
		AmountNoFee:   db.NewDbCoins(s.AmountNoFee),
		Fee:           db.NewDbCoins(s.Fee),
	}
}

// toDailyStatsDomainList - mapping func to a domain list.
func toDailyStatsDomainList(s []db.FeeExcluderDailyStats) []fe.DailyStats {
	res := make([]fe.DailyStats, 0, len(s))
	for _, ds := range s {
		res = append(res, toDailyStatsDomain(ds))
	}

	return res
}

// BLOCK Stats

// toStatsDomain - mapping func to a domain model.
func toStatsDomain(dailyStats *types.DailyStats, s db.FeeExcluderStats) types.Stats {
	return types.Stats{
		Index: s.ID,
		Date:  s.Date.Format(layoutDate),
		Stats: dailyStats,
	}
}

// toStatsDatabase - mapping func to a database model.
func toStatsDatabase(dailyStatsID uint64, s fe.Stats) (db.FeeExcluderStats, error) {
	date, err := time.Parse(layoutDate, s.Date)
	if err != nil {
		return db.FeeExcluderStats{}, err
	}

	return db.FeeExcluderStats{
		ID:           s.Index,
		DailyStatsID: dailyStatsID,
		Date:         date,
	}, nil
}

// toStatsDomainList - mapping func to a domain list.
func toStatsDomainList(dailyStats *types.DailyStats, s []db.FeeExcluderStats) []fe.Stats {
	res := make([]fe.Stats, 0, len(s))
	for _, ds := range s {
		res = append(res, toStatsDomain(dailyStats, ds))
	}

	return res
}
