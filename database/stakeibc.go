package database

import (
	"encoding/json"
	"fmt"

	"github.com/forbole/bdjuno/v3/types"
)

// SaveStakeIBCParams allows to store the given params inside the database
func (db *Db) SaveStakeIBCParams(params *types.StakeIBCParams) error {
	paramsBz, err := json.Marshal(&params.Params)
	if err != nil {
		return fmt.Errorf("error while marshaling stakeibc params: %s", err)
	}

	stmt := `
INSERT INTO stakeibc_params (params, height) 
VALUES ($1, $2)
ON CONFLICT (one_row_id) DO UPDATE 
    SET params = excluded.params,
        height = excluded.height
WHERE stakeibc_params.height <= excluded.height`

	_, err = db.Sql.Exec(stmt, string(paramsBz), params.Height)
	if err != nil {
		return fmt.Errorf("error while storing stakeibc params: %s", err)
	}

	return nil
}
