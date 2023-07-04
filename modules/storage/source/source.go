package source

import (
	storagetypes "github.com/jackalLabs/canine-chain/v3/x/storage/types"
)

type Source interface {
	Params(height int64) (storagetypes.Params, error)
	Providers(height int64) ([]storagetypes.Providers, error)
}
