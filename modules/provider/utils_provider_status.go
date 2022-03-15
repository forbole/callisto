package provider

import (
	providertypes "github.com/ovrclk/akash/x/provider/types/v1beta2"
)

func (m *Module) getProviderStatuses(height int64) ([]providertypes.Provider, error) {
	providers, err := m.source.GetProviders(height)
	if err != nil {
		return nil, err
	}

	return providers, nil
}
