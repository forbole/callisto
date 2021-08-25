package types

// IscnParamsRow represents a single row inside the iscn_params table
type IscnParamsRow struct {
	OneRowID bool   `db:"one_row_id"`
	Params   string `db:"params"`
	Height   int64  `db:"height"`
}

// NewIscnParamsRow builds a new IscnParamsRow instance
func NewIscnParamsRow(
	params string, height int64,
) IscnParamsRow {
	return IscnParamsRow{
		OneRowID: true,
		Params:   params,
		Height:   height,
	}
}

// Equal reports whether m and n represent the same table rows.
func (m IscnParamsRow) Equal(n IscnParamsRow) bool {
	return m.Params == n.Params &&
		m.Height == n.Height
}
