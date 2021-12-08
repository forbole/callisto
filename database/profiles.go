package database

import (
	"encoding/json"
	"fmt"

	"github.com/forbole/bdjuno/v2/types"
)

// SaveProfilesParams save the params of profiles module in the database
func (db *Db) SaveProfilesParams(params *types.ProfilesParams) error {
	paramsBz, err := json.Marshal(&params.Params)
	if err != nil {
		return fmt.Errorf("error while marshaling profiles params: %s", err)
	}

	stmt := `
INSERT INTO profiles_params (params, height) 
VALUES ($1, $2) 
ON CONFLICT (one_row_id) DO UPDATE 
    SET params = excluded.params,
      	height = excluded.height
WHERE profiles_params.height <= excluded.height`

	_, err = db.Sql.Exec(stmt, string(paramsBz), params.Height)

	if err != nil {
		return fmt.Errorf("error while storing profiles params: %s", err)
	}

	return nil
}
