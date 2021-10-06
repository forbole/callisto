package types

// OracleParamsRow represents a single row inside the oracle_params table
type OracleParamsRow struct {
	OneRowID bool   `db:"one_row_id"`
	Params   string `db:"params"`
	Height   int64  `db:"height"`
}

// NewOracleParamsRow builds a new OracleParamsRow instance
func NewOracleParamsRow(
	params string, height int64,
) OracleParamsRow {
	return OracleParamsRow{
		OneRowID: true,
		Params:   params,
		Height:   height,
	}
}

// Equal reports whether m and n represent the same table rows.
func (m OracleParamsRow) Equal(n OracleParamsRow) bool {
	return m.Params == n.Params &&
		m.Height == n.Height
}
