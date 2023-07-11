package source

import (
	storagetypes "github.com/MonikaCat/canine-chain/v2/x/storage/types"
)

type Source interface {
	Params(height int64) (storagetypes.Params, error)
	Providers(height int64) ([]storagetypes.Providers, error)
	Strays(height int64) ([]storagetypes.Strays, error)
}
