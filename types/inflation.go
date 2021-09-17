package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type EmoneyInflation struct {
	Issuer string
	Denom  string
	Rate   sdk.Dec
	Height int64
}

func NewEmoneyInflation(issuer string, denom string, rate sdk.Dec, height int64) EmoneyInflation {
	return EmoneyInflation{
		Issuer: issuer,
		Denom:  denom,
		Rate:   rate,
		Height: height,
	}
}
