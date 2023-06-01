package source

import (
	stakersquerytypes "github.com/KYVENetwork/chain/x/query/types"
	stakerstypes "github.com/KYVENetwork/chain/x/stakers/types"
)

type Source interface {
	Params(height int64) (stakerstypes.Params, error)
	Stakers(height int64) ([]stakersquerytypes.FullStaker, error)
}
