package types

// InterchainStakingParamsRow represents a single row inside the interchain_staking_params table
type InterchainStakingParamsRow struct {
	OneRowID bool   `db:"one_row_id"`
	Params   string `db:"params"`
	Height   int64  `db:"height"`
}

// NewInterchainStakingParamsRow builds a new InterchainStakingParamsRow instance
func NewInterchainStakingParamsRow(
	params string, height int64,
) InterchainStakingParamsRow {
	return InterchainStakingParamsRow{
		OneRowID: true,
		Params:   params,
		Height:   height,
	}
}

// Equal reports whether m and n represent the same table rows.
func (m InterchainStakingParamsRow) Equal(n InterchainStakingParamsRow) bool {
	return m.Params == n.Params &&
		m.Height == n.Height
}
