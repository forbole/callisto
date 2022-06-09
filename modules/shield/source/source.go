package source

import shieldtypes "github.com/certikfoundation/shentu/v2/x/shield/types"

type Source interface {
	GetPoolParams(height int64) (shieldtypes.PoolParams, error)
	GetPools(height int64) ([]shieldtypes.Pool, error)
	GetShieldProviders(height int64) ([]shieldtypes.Provider, error)
	GetShieldStatus(height int64) (*shieldtypes.QueryShieldStatusResponse, error)
}
