package history

import (
	"github.com/forbole/bdjuno/types"
)

// UpdatePricesHistory stores the given prices inside the price history table
func (m *Module) UpdatePricesHistory(prices []types.TokenPrice) error {
	if !m.cfg.IsModuleEnabled(moduleName) {
		return nil
	}

	return m.db.SaveTokenPricesHistory(prices)
}
