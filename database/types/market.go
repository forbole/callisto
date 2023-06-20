package types

// MarketParamsRow represents a single row inside the market_params table
type MarketParamsRow struct {
	OneRowID bool   `db:"one_row_id"`
	Params   string `db:"params"`
	Height   int64  `db:"height"`
}

// NewMarketParamsRow builds a new MarketParamsRow instance
func NewMarketParamsRow(
	params string, height int64,
) MarketParamsRow {
	return MarketParamsRow{
		OneRowID: true,
		Params:   params,
		Height:   height,
	}
}

// Equal reports whether m and n represent the same table rows.
func (m MarketParamsRow) Equal(n MarketParamsRow) bool {
	return m.Params == n.Params &&
		m.Height == n.Height
}
