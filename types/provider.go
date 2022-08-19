package types

import (
	akashtypes "github.com/ovrclk/akash/types/v1beta2"
	providertypes "github.com/ovrclk/akash/x/provider/types/v1beta2"
)

type Provider struct {
	OwnerAddress string
	HostURI      string
	Attributes   []akashtypes.Attribute
	Info         providertypes.ProviderInfo
	Height       int64
}

// NewProvider allows to build a new Provider instance
func NewProvider(p providertypes.Provider, height int64) *Provider {
	return &Provider{
		OwnerAddress: p.Owner,
		HostURI:      p.HostURI,
		Attributes:   p.Attributes,
		Info:         p.Info,
		Height:       height,
	}
}
