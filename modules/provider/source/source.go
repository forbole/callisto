package source

import (
	providertypes "github.com/ovrclk/akash/x/provider/types/v1beta2"
)

type Source interface {
	GetProviders(height int64) ([]providertypes.Provider, error)
}
