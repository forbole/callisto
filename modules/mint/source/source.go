package source

import (
	minttypes "github.com/MonikaCat/canine-chain/v2/x/jklmint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Source interface {
	GetInflation(height int64) (sdk.Dec, error)
	Params(height int64) (minttypes.Params, error)
}
