package source

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/osmosis-labs/osmosis/v14/x/mint/types"
)

type Source interface {
	Params(height int64) (minttypes.Params, error)
	EpochProvisions(height int64) (sdk.Dec, error)
}
