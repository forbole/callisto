package types

import (
	"github.com/ovrclk/akash/provider"
	clustertypes "github.com/ovrclk/akash/provider/cluster/types/v1beta2"
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

type ProviderStatus struct {
	ProviderAddress       string
	Active                bool
	LeasesNumber          uint32
	BidengineOrders       uint32
	ManifestDeployments   uint32
	ClusterPublicHostname string
	InventoryStatusRaw    clustertypes.InventoryStatus
	ActiveInventorySum    *Resource
	PendingInventorySum   *Resource
	AvailableInventorySum *Resource
	Height                int64
}

// NewProviderStatus allows to build a new ProviderStatus instance
func NewProviderStatus(
	providerAddress string, active bool, s *provider.Status,
	activeInventorySum *Resource, pendingInventorySum *Resource, availableInventorySum *Resource,
	height int64,
) *ProviderStatus {
	return &ProviderStatus{
		ProviderAddress:       providerAddress,
		Active:                active,
		LeasesNumber:          s.Cluster.Leases,
		BidengineOrders:       s.Bidengine.Orders,
		ManifestDeployments:   s.Manifest.Deployments,
		ClusterPublicHostname: s.ClusterPublicHostname,
		InventoryStatusRaw:    s.Cluster.Inventory,
		ActiveInventorySum:    activeInventorySum,
		PendingInventorySum:   pendingInventorySum,
		AvailableInventorySum: availableInventorySum,
		Height:                height,
	}
}

// NewInactiveProviderStatus allows to build a new ProviderStatus instance of inactive status
func NewInactiveProviderStatus(
	providerAddress string, active bool, s *provider.Status,
	activeInventorySum *Resource, pendingInventorySum *Resource, availableInventorySum *Resource,
	height int64,
) *ProviderStatus {
	return &ProviderStatus{
		ProviderAddress: providerAddress,
		Active:          active,
		Height:          height,
	}
}

type Resource struct {
	CPU              uint64
	Memory           uint64
	StorageEphemeral uint64
}

func NewProviderResouce(
	cpu uint64, memory uint64, storageEphemeral uint64,
) *Resource {
	return &Resource{
		CPU:              cpu,
		Memory:           memory,
		StorageEphemeral: storageEphemeral,
	}
}
