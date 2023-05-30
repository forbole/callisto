package source

import (
	poolquerytypes "github.com/KYVENetwork/chain/x/query/types"
)

type Source interface {
	Pools(height int64) ([]poolquerytypes.PoolResponse, error)
}
