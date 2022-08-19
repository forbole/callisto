package provider

import (
	"github.com/forbole/bdjuno/v3/types"
	"github.com/ovrclk/akash/provider"
)

func (m *Module) getProviderInventory(address string, height int64, invCh chan *types.ProviderStatus) {
	// Get the provision status of a provider
	status, err := m.source.GetProviderInventoryStatus(address)
	if err != nil {
		invCh <- types.NewInactiveProviderStatus(
			address, false, nil, nil, nil, nil, height,
		)
		return
	}

	// Calculate inventory sum of each state
	active, pending, available := m.calculateInventorySum(status)

	invCh <- types.NewProviderStatus(
		address, true, status, active, pending, available, height,
	)
}

// calculateInventorySum calculates the sum of inventory in different statuses
func (m *Module) calculateInventorySum(status *provider.Status) (
	activeInventorySum *types.Resource, pendingInventorySum *types.Resource, availableInventorySum *types.Resource,
) {
	var cpu uint64 = 0
	var memory uint64 = 0
	var storage uint64 = 0

	// Sum up active inventory
	for _, active := range status.Cluster.Inventory.Active {
		cpu += active.CPU
		memory += active.Memory
		storage += active.StorageEphemeral
	}
	activeInventorySum = types.NewProviderResouce(cpu, memory, storage)

	// Sum up pending inventory
	cpu = 0
	memory = 0
	storage = 0
	for _, pending := range status.Cluster.Inventory.Pending {
		cpu += pending.CPU
		memory += pending.Memory
		storage += pending.StorageEphemeral
	}
	pendingInventorySum = types.NewProviderResouce(cpu, memory, storage)

	// Sum up available inventory
	cpu = 0
	memory = 0
	storage = 0
	for _, available := range status.Cluster.Inventory.Available.Nodes {
		cpu += available.CPU
		memory += available.Memory
		storage += available.StorageEphemeral
	}
	availableInventorySum = types.NewProviderResouce(cpu, memory, storage)

	return
}
