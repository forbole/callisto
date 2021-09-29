package pricefeed

import "github.com/forbole/bdjuno/types"

type HistoryModule interface {
	UpdatePricesHistory([]types.TokenPrice) error
}
