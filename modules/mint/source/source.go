package source

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
)

type Source interface {
	GetInflation(height int64) (sdk.Dec, error)
	Params(height int64) (minttypes.Params, error)
}
