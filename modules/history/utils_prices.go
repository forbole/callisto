package history

import (
	"github.com/forbole/bdjuno/v2/types"
)

// IsHistoryModuleEnabled checks if history module is enabled inside config.yaml file
func (m *Module) IsHistoryModuleEnabled() bool {
	if m.cfg.IsModuleEnabled(moduleName) {
		return true
	}
	return false
}

// UpdatePricesHistory stores the given prices inside the price history table
func (m *Module) UpdatePricesHistory(prices []types.TokenPrice) error {
	return m.db.SaveTokenPricesHistory(prices)
}
