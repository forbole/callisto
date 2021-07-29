package types

// StakingParamsRow represents a single row inside the staking_params table
type StakingParamsRow struct {
	BondName          string `db:"bond_denom"`
	MinCommissionRate string `db:"min_commission_rate"`
	UnbondingTime     uint64 `db:"unbonding_time"`
	Height            int64  `db:"height"`
	MaxEntries        uint32 `db:"max_entries"`
	HistoricalEntries uint32 `db:"historical_entries"`
	MaxValidators     uint32 `db:"max_validators"`
	OneRowID          bool   `db:"one_row_id"`
}

// NewStakingParamsRow allows to build a new StakingParamsRow object
func NewStakingParamsRow(
	bondName string, unbondingTime uint64, maxEntries uint32, historicalEntries uint32, maxValidators uint32, minCommissionRate string, height int64,
) StakingParamsRow {
	return StakingParamsRow{
		OneRowID:          true,
		BondName:          bondName,
		UnbondingTime:     unbondingTime,
		MaxEntries:        maxEntries,
		HistoricalEntries: historicalEntries,
		MaxValidators:     maxValidators,
		MinCommissionRate: minCommissionRate,
		Height:            height,
	}
}

// Equal tells whether r and s contain the same data
func (r StakingParamsRow) Equal(s StakingParamsRow) bool {
	return r.BondName == s.BondName &&
		r.UnbondingTime == s.UnbondingTime &&
		r.MaxEntries == s.MaxEntries &&
		r.HistoricalEntries == s.HistoricalEntries &&
		r.MaxValidators == s.MaxValidators &&
		r.MinCommissionRate == s.MinCommissionRate &&
		r.Height == s.Height
}
