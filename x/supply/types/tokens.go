package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ValidatorUptime contains the uptime information of a single
// validator for a specific height and point in time
type TotalCoins struct {
	Denom  string
	Amount sdk.Int
	Height int64
}

// NewValidatorUptime allows to build a new ValidatorUptime instance
func NewTotalCoins(denom string, amount sdk.Int, height int64) TotalCoins {
	return TotalCoins{
		Denom:  denom,
		Amount: amount,
		Height: height,
	}
}

// Equal tells whether v and w represent the same uptime
func (v TotalCoins) Equal(w TotalCoins) bool {
	return v.Denom == w.Denom &&
		v.Amount == w.Amount &&
		v.Height == w.Height
}
