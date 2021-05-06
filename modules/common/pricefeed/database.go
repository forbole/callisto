package pricefeed

import (
	"github.com/forbole/bdjuno/types"
)

// DB represents a generic database that allows to perform token prices related operations
type DB interface {
	SaveToken(token types.Token) error
	GetTokenUnits() ([]string, error)
	SaveTokensPrices(prices []types.TokenPrice) error
}
