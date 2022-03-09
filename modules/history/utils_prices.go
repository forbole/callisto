package history

import (
	"github.com/forbole/bdjuno/v2/types"
)

// UpdatePricesHistory stores the given prices inside the price history table
func (m *Module) UpdatePricesHistory(prices []types.TokenPrice) error {
	return m.db.SaveTokenPricesHistory(prices)
}
