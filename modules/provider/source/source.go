package source

import (
	providertypes "github.com/akash-network/akash-api/go/node/provider/v1beta3"
)

type Source interface {
	GetProvider(height int64, ownerAddress string) (providertypes.Provider, error)
	GetProviders(height int64) ([]providertypes.Provider, error)
}
