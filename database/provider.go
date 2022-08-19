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
	}

	stmt = stmt[:len(stmt)-1]
	stmt += `
ON CONFLICT ON CONSTRAINT unique_provider DO UPDATE
	SET host_uri = excluded.host_uri, 
		attributes = excluded.attributes,
		info = excluded.info,
	    height = excluded.height 
WHERE akash_provider.height <= excluded.height`

	_, err := db.Sql.Exec(stmt, params...)
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
