package source

import (
	providertypes "github.com/ovrclk/akash/x/provider/types/v1beta2"
)

type Source interface {
	GetProvider(height int64, ownerAddress string) (providertypes.Provider, error)
	GetProviders(height int64) ([]providertypes.Provider, error)
}
