package types

// InflationRow represents a single row inside the inflation table
type InflationRow struct {
	OneRowID bool    `db:"one_row_id"`
	Value    float64 `db:"value"`
	Height   int64   `db:"height"`
}

// NewInflationRow builds a new InflationRows instance
func NewInflationRow(value float64, height int64) InflationRow {
	return InflationRow{
		OneRowID: true,
		Value:    value,
		Height:   height,
	}
}

// Equal reports whether i and j represent the same table rows.
func (i InflationRow) Equal(j InflationRow) bool {
	return i.Value == j.Value &&
		i.Height == j.Height
}

// --------------------------------------------------------------------------------------------------------------------

// MintParamsRow represents a single row inside the mint_params table
type MintParamsRow struct {
	OneRowID            bool   `db:"one_row_id"`
	MintDenom           string `db:"mint_denom"`
	InflationRateChange string `db:"inflation_rate_change"`
	InflationMin        string `db:"inflation_min"`
	InflationMax        string `db:"inflation_max"`
	GoalBonded          string `db:"goal_bonded"`
	BlocksPerYear       uint64 `db:"blocks_per_year"`
	Height              int64  `db:"height"`
}

// NewMintParamsRow builds a new MintParamsRow instance
func NewMintParamsRow(
	mintDenom, inflationRateChange, inflationMin, inflationMax, goalBonded string, blockPerYear uint64, height int64,
) MintParamsRow {
	return MintParamsRow{
		OneRowID:            true,
		MintDenom:           mintDenom,
		InflationRateChange: inflationRateChange,
		InflationMin:        inflationMin,
		InflationMax:        inflationMax,
		GoalBonded:          goalBonded,
		BlocksPerYear:       blockPerYear,
		Height:              height,
	}
}

// Equal reports whether m and n represent the same table rows.
func (m MintParamsRow) Equal(n MintParamsRow) bool {
	return m.MintDenom == n.MintDenom &&
		m.InflationRateChange == n.InflationRateChange &&
		m.InflationMin == n.InflationMin &&
		m.InflationMax == n.InflationMax &&
		m.GoalBonded == n.GoalBonded &&
		m.BlocksPerYear == n.BlocksPerYear &&
		m.Height == n.Height
}
