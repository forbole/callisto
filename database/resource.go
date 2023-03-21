package database

import (
	"encoding/json"
	"fmt"

	"github.com/forbole/bdjuno/v4/types"
)

// SaveResource allows to store the given resource inside the database
func (db *Db) SaveResource(resource *types.Resource) error {
	dataBz, err := json.Marshal(&resource.Data)
	if err != nil {
		return fmt.Errorf("error while marshaling resource data: %s", err)
	}

	alsoKnownAsBz, err := json.Marshal(&resource.AlsoKnownAs)
	if err != nil {
		return fmt.Errorf("error while marshaling resource: %s", err)
	}

	stmt := `
INSERT INTO resource (id, collection_id, data, name, version,
	resource_type, also_known_as, from_address, height) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
ON CONFLICT (id) DO UPDATE 
    SET collection_id = excluded.collection_id,
        data = excluded.data,
        name = excluded.name,
        version = excluded.version,
        resource_type = excluded.resource_type,
        also_known_as = excluded.also_known_as,
		from_address = excluded.from_address,
        height = excluded.height
WHERE resource.height <= excluded.height`

	_, err = db.SQL.Exec(stmt, resource.ID, resource.CollectionID, string(dataBz),
		resource.Name, resource.Version, resource.ResourceType, string(alsoKnownAsBz),
		resource.FromAddress, resource.Height)
	if err != nil {
		return fmt.Errorf("error while storing resource in db: %s", err)
	}

	return nil
}
