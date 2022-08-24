package provider

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/forbole/bdjuno/v3/types"
	"github.com/ovrclk/akash/provider"
	"github.com/rs/zerolog/log"
)

func (m *Module) updateProviderInventoryStatus(address string, height int64) {
	// Get the inventory status of a provider
	status, err := m.getProviderInventoryStatus(address)
	if err != nil {
		err := m.db.SetProviderStatus(address, false, height)
		if err != nil {
			log.Error().Str("module", "provider").
				Msgf("error while setting provider status inactive %s", err)
		}
		return
	}

	// Calculate inventory sum of each state
	active, pending, available := m.calculateInventorySum(status)

	err = m.db.SaveProviderInventoryStatus(types.NewProviderInventoryStatus(
		address, true, status, active, pending, available, height,
	))
	if err != nil {
		log.Error().Str("module", "provider").
			Msgf("error while storing provider inventory status %s", err)
	}
}

func (m *Module) getProviderInventoryStatus(hostURI string) (*provider.Status, error) {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: transport}
	res, err := client.Get(hostURI)
	if err != nil {
		fmt.Println(err)
	}

	defer res.Body.Close()

	var providerStatus = new(provider.Status)
	err = json.NewDecoder(res.Body).Decode(providerStatus)
	if err != nil {
		return nil, fmt.Errorf("error while reading json response : %s", err)
	}

	return providerStatus, nil
}

// calculateInventorySum calculates the sum of inventory in different statuses
func (m *Module) calculateInventorySum(status *provider.Status) (
	*types.Resource, *types.Resource, *types.Resource,
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
	activeInventorySum := types.NewProviderResouce(cpu, memory, storage)

	// Sum up pending inventory
	cpu = 0
	memory = 0
	storage = 0
	for _, pending := range status.Cluster.Inventory.Pending {
		cpu += pending.CPU
		memory += pending.Memory
		storage += pending.StorageEphemeral
	}
	pendingInventorySum := types.NewProviderResouce(cpu, memory, storage)

	// Sum up available inventory
	cpu = 0
	memory = 0
	storage = 0
	for _, available := range status.Cluster.Inventory.Available.Nodes {
		cpu += available.CPU
		memory += available.Memory
		storage += available.StorageEphemeral
	}
	availableInventorySum := types.NewProviderResouce(cpu, memory, storage)

	return activeInventorySum, pendingInventorySum, availableInventorySum
}
