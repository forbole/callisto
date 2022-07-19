package source

import (
	minttypes "github.com/MonOsmosis/osmosis/v10/x/mint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Source interface {
	Params(height int64) (minttypes.Params, error)
	EpochProvisions(height int64) (sdk.Dec, error)
}
