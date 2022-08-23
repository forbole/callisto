package database

import (
	"encoding/json"
	"fmt"

	dbtypes "github.com/forbole/bdjuno/v3/database/types"
	dbutils "github.com/forbole/bdjuno/v3/database/utils"
	"github.com/forbole/bdjuno/v3/types"
)

// SaveProviders allows to store the providers inside the database
func (db *Db) SaveProviders(providers []*types.Provider, height int64) error {
	if len(providers) == 0 {
		return nil
	}

	paramsNumber := 5
	slices := dbutils.SplitProviders(providers, paramsNumber)

	for _, providers := range slices {
		if len(providers) == 0 {
			continue
		}

		// Store providers
		err := db.saveProviders(paramsNumber, providers, height)
		if err != nil {
			return fmt.Errorf("error while storing providers: %s", err)
		}
	}

	return nil
}

// saveProviders allows to store providers inside the database
func (db *Db) saveProviders(paramsNumber int, providers []*types.Provider, height int64) error {

	stmt := `INSERT INTO akash_provider (owner_address, host_uri, attributes, info, height) VALUES `
	var params []interface{}

	accounts := make([]types.Account, len(providers))
	for i, provider := range providers {

		bi := i * paramsNumber
		stmt += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d),", bi+1, bi+2, bi+3, bi+4, bi+5)

		attributesBz, err := json.Marshal(&provider.Attributes)
		if err != nil {
			return fmt.Errorf("error while marshaling provider attributes: %s", err)
		}

		info := dbtypes.NewDbInfo(provider.Info)
		infoValue, err := info.Value()
		if err != nil {
			return fmt.Errorf("error while converting provider info to DbProviderInfo: %s", err)
		}

		params = append(params, provider.OwnerAddress, provider.HostURI, string(attributesBz), infoValue, height)
		accounts[i] = types.NewAccount(provider.OwnerAddress)
	}

	stmt = stmt[:len(stmt)-1]
	stmt += `
ON CONFLICT ON CONSTRAINT unique_provider DO UPDATE
	SET host_uri = excluded.host_uri, 
		attributes = excluded.attributes,
		info = excluded.info,
	    height = excluded.height 
WHERE akash_provider.height <= excluded.height`

	err := db.SaveAccounts(accounts)
	if err != nil {
		return fmt.Errorf("error while storing accounts: %s", err)
	}

	_, err = db.Sql.Exec(stmt, params...)
	if err != nil {
		return fmt.Errorf("error while storing providers: %s", err)
	}

	return nil
}

// DeleteProvider allows to remove provider record from the database
func (db *Db) DeleteProvider(ownerAddress string) error {
	stmt := `DELETE FROM providers WHERE owner = $1`
	_, err := db.Sql.Exec(stmt, ownerAddress)
	if err != nil {
		return fmt.Errorf("error while deleting provider: %s", err)
	}
	return nil
}

// GetAkashProviders returns the akash provider addresses stored inside the database
func (db *Db) GetAkashProviders() ([]string, error) {
	stmt := `SELECT owner_address FROM akash_provider`

	var addresses []string
	if err := db.Sqlx.Select(&addresses, stmt); err != nil {
		return nil, err
	}

	if len(addresses) == 0 {
		return nil, fmt.Errorf("no provider was saved")
	}

	return addresses, nil
}

// SaveProviderInventoryStatus allows to store provider inventory status inside the database
func (db *Db) SaveProviderInventoryStatus(status *types.ProviderInventoryStatus) error {
	stmt := `INSERT INTO akash_provider_inventory 
	( 
		provider_address, active, lease_count, 
		bidengine_order_count, manifest_deployment_count, 
		cluster_public_hostname, inventory_status_raw, 
		active_inventory_sum, pending_inventory_sum, available_inventory_sum, 
		height 
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) 
ON CONFLICT (provider_address) DO UPDATE 
	SET active = excluded.active, 
		lease_count = excluded.lease_count, 
		bidengine_order_count = excluded.bidengine_order_count, 
		manifest_deployment_count = excluded.manifest_deployment_count, 
		cluster_public_hostname = excluded.cluster_public_hostname, 
		inventory_status_raw = excluded.inventory_status_raw, 
		active_inventory_sum = excluded.active_inventory_sum, 
		pending_inventory_sum = excluded.pending_inventory_sum, 
		available_inventory_sum = excluded.available_inventory_sum, 
		height = excluded.height 
	WHERE akash_provider_inventory.height <= excluded.height`

	// Marshal inventory raw status
	inventoryStatusBz, err := json.Marshal(&status.InventoryStatusRaw)
	if err != nil {
		return fmt.Errorf("error while marshaling provider inventory raw status: %s", err)
	}

	// Get value for active inventory sum
	active := dbtypes.NewDbAkashResource(status.ActiveInventorySum)
	activeValue, _ := active.Value()

	// Get value for active inventory sum
	pending := dbtypes.NewDbAkashResource(status.PendingInventorySum)
	pendingValue, _ := pending.Value()

	// Get value for active inventory sum
	available := dbtypes.NewDbAkashResource(status.AvailableInventorySum)
	availableValue, _ := available.Value()

	_, err = db.Sql.Exec(stmt,
		status.ProviderAddress, status.Active, status.LeaseCount,
		status.BidengineOrderCount, status.ManifestDeploymentCount,
		status.ClusterPublicHostname, string(inventoryStatusBz),
		activeValue, pendingValue, availableValue,
		status.Height,
	)
	if err != nil {
		return fmt.Errorf("error while storing inventory status of provider %s: %s", status.ProviderAddress, err)
	}

	return nil
}

// SetProviderStatus allows to update the provider status inside the database
func (db *Db) SetProviderStatus(providerAddress string, active bool, height int64) error {
	stmt := `UPDATE akash_provider_inventory SET active = $1 WHERE provider_address = $2`
	_, err := db.Sql.Exec(stmt, active, providerAddress)
	if err != nil {
		return fmt.Errorf("error while updating active status of provider %s: %s", providerAddress, err)
	}
	return nil
}
