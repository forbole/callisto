package database

import (
	"encoding/json"
	"fmt"

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
