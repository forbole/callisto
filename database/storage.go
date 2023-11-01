package database

import (
	"encoding/json"
	"fmt"

	storagetypes "github.com/jackalLabs/canine-chain/v3/x/storage/types"
	"github.com/lib/pq"

	"github.com/forbole/bdjuno/v4/types"
)

// SaveStorageParams allows to store the given params inside the database
func (db *Db) SaveStorageParams(params *types.StorageParams) error {
	paramsBz, err := json.Marshal(&params.Params)
	if err != nil {
		return fmt.Errorf("error while marshaling storage params: %s", err)
	}

	stmt := `
INSERT INTO storage_params (params, height) 
VALUES ($1, $2)
ON CONFLICT (one_row_id) DO UPDATE 
    SET params = excluded.params,
        height = excluded.height
WHERE storage_params.height <= excluded.height`

	_, err = db.SQL.Exec(stmt, string(paramsBz), params.Height)
	if err != nil {
		return fmt.Errorf("error while storing storage params: %s", err)
	}

	return nil
}

// SaveStorageProviders allows the bulk saving of a list of storage providers.
func (db *Db) SaveStorageProviders(providers []storagetypes.Providers, height int64) error {
	if len(providers) == 0 {
		return nil
	}

	storageProvidersQuery := `
INSERT INTO storage_providers (address, ip, total_space, burned_contracts, creator, keybase_identity, auth_claimers, height) 
VALUES `
	var storageProviders []interface{}
	var accounts []types.Account

	for i, provider := range providers {
		vi := i * 8 // Starting position for storage providers
		accounts = append(accounts, types.NewAccount(provider.Address))

		storageProvidersQuery += fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d),", vi+1, vi+2, vi+3, vi+4, vi+5, vi+6, vi+7, vi+8)
		storageProviders = append(storageProviders,
			provider.Address,
			provider.Ip,
			provider.Totalspace,
			provider.BurnedContracts,
			provider.Creator,
			provider.KeybaseIdentity,
			pq.StringArray(provider.AuthClaimers),
			height,
		)
	}

	// Store the accounts
	err := db.SaveAccounts(accounts)
	if err != nil {
		return fmt.Errorf("error while storing providers accounts: %s", err)
	}

	storageProvidersQuery = storageProvidersQuery[:len(storageProvidersQuery)-1] // Remove the trailing ","
	storageProvidersQuery += `
ON CONFLICT (address) DO UPDATE 
	SET ip = excluded.ip,
		total_space = excluded.total_space,
		burned_contracts = excluded.burned_contracts,
		creator = excluded.creator,
		keybase_identity = excluded.keybase_identity,
		auth_claimers = excluded.auth_claimers,
		height = excluded.height
WHERE storage_providers.height <= excluded.height`
	_, err = db.SQL.Exec(storageProvidersQuery, storageProviders...)
	if err != nil {
		return fmt.Errorf("error while storing storage providers infos: %s", err)
	}

	return nil
}
