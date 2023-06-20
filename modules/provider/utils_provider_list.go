package provider

import "github.com/forbole/bdjuno/v4/types"

func (m *Module) getProviderList(height int64) ([]*types.Provider, error) {
	providers, err := m.source.GetProviders(height)
	if err != nil {
		return nil, err
	}

	infos := make([]*types.Provider, len(providers))
	for index, info := range providers {
		infos[index] = types.NewProvider(info, height)
	}

	return infos, nil
}
