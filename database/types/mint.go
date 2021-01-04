package types

// InflationRow represents a single row inside the inflation table
type InflationRow struct {
	Value  float64 `db:"value"`
	Height int64   `db:"height"`
}

// NewInflationRow builds a new InflationRows instance
func NewInflationRow(value float64, height int64) InflationRow {
	return InflationRow{
		Value:  value,
		Height: height,
	}
}

// Equal reports whether i and j represent the same table rows.
func (i InflationRow) Equal(j InflationRow) bool {
	return i.Value == j.Value &&
		i.Height == j.Height
}
