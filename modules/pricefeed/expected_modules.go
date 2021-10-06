package pricefeed

import "github.com/forbole/bdjuno/v2/types"

type HistoryModule interface {
	UpdatePricesHistory([]types.TokenPrice) error
}
