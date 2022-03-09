package pricefeed

import "github.com/forbole/bdjuno/v2/types"

type HistoryModule interface {
	IsHistoryModuleEnabled() bool 
	UpdatePricesHistory([]types.TokenPrice) error
}
