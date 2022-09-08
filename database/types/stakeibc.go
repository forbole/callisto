package types

// StakeIBCParamsRow represents a single row inside the stakeibc_params table
type StakeIBCParamsRow struct {
	OneRowID bool   `db:"one_row_id"`
	Params   string `db:"params"`
	Height   int64  `db:"height"`
}

// NewStakeIBCParamsRow builds a new StakeIBCParamsRow instance
func NewStakeIBCParamsRow(
	params string, height int64,
) StakeIBCParamsRow {
	return StakeIBCParamsRow{
		OneRowID: true,
		Params:   params,
		Height:   height,
	}
}

// Equal reports whether m and n represent the same table rows.
func (m StakeIBCParamsRow) Equal(n StakeIBCParamsRow) bool {
	return m.Params == n.Params &&
		m.Height == n.Height
}
