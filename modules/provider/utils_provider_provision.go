package provider

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"

	provider "github.com/akash-network/provider"
	"github.com/forbole/bdjuno/v4/types"
	"github.com/rs/zerolog/log"
)

var (
	statusEndpoint string = "/status"
)

func (m *Module) updateProviderInventoryStatus(address string, hostURI string, height int64) {
	// Get the inventory status of a provider
	statusURL := hostURI + statusEndpoint
	status, err := m.getProviderInventoryStatus(statusURL)
	if err != nil {
		err := m.db.SetProviderStatus(address, false, height)
		if err != nil {
			log.Error().Str("module", "provider").
				Msgf("error while setting provider status inactive: %s", err)
		}
		return
	}

	if status.Cluster != nil {
		// Calculate inventory sum of each state
		active, pending, available := m.calculateInventorySum(status)

		err = m.db.SaveProviderInventoryStatus(types.NewProviderInventoryStatus(
			address, true, status, active, pending, available, height,
		))
		if err != nil {
			log.Error().Str("module", "provider").
				Msgf("error while storing provider inventory status: %s", err)
		}
	}

}

// getProviderInventoryStatus allows to get provider inventory status by rest client with insecure TLS config
func (m *Module) getProviderInventoryStatus(statusURL string) (*provider.Status, error) {
	transport := &http.Transport{
		//nolint:gosec
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	client := &http.Client{Transport: transport}
	res, err := client.Get(statusURL)
	if err != nil {
		return nil, fmt.Errorf("error while getting inventory status with client: %s", err)
	}

	defer res.Body.Close()

	var providerStatus = new(provider.Status)
	err = json.NewDecoder(res.Body).Decode(providerStatus)
	if err != nil {
		return nil, fmt.Errorf("error while reading json response: %s", err)
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
