package source

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/ingenuity-build/quicksilver/x/mint/types"
)

type Source interface {
	GetEpochProvisions(height int64) (sdk.Dec, error)
	Params(height int64) (minttypes.Params, error)
}
