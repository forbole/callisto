package database

import (
	"encoding/json"
	"fmt"

	"github.com/forbole/bdjuno/v3/types"
)

// SaveMarkerParams allows to store the given params inside the database
func (db *Db) SaveMarkerParams(params *types.MarkerParams) error {
	paramsBz, err := json.Marshal(&params.Params)
	if err != nil {
		return fmt.Errorf("error while marshaling marker params: %s", err)
	}

	stmt := `
INSERT INTO marker_params (params, height) 
VALUES ($1, $2)
ON CONFLICT (one_row_id) DO UPDATE 
    SET params = excluded.params,
        height = excluded.height
WHERE marker_params.height <= excluded.height`

	_, err = db.Sql.Exec(stmt, string(paramsBz), params.Height)
	if err != nil {
		return fmt.Errorf("error while storing marker params: %s", err)
	}

	return nil
}

// SaveMarkersAccounts allows to store the markers accounts for the given block height
func (db *Db) SaveMarkersAccounts(markersList []types.MarkerAccount, height int64) error {
	if len(markersList) == 0 {
		return nil
	}

	stmt := `
	INSERT INTO marker_account (address, access_control, allow_governance_control, 
		denom, marker_type, status, supply, height)
	VALUES `
	var accounts []types.Account
	var markerParams []interface{}

	for i, marker := range markersList {

		accessControl, err := json.Marshal(&marker.AccessControl)
		if err != nil {
			return err
		}

		supplyValue, err := json.Marshal(&marker.Supply)
		if err != nil {
			return err
		}
		// Prepare the account query
		accounts = append(accounts, types.NewAccount(marker.Address))

		// Prepare the marker query
		vi := i * 8
		stmt += fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d),",
			vi+1, vi+2, vi+3, vi+4, vi+5, vi+6, vi+7, vi+8)

		markerParams = append(markerParams,
			marker.Address,
			string(accessControl),
			marker.AllowGovernanceControl,
			marker.Denom,
			marker.MarkerType.String(),
			marker.Status,
			string(supplyValue),
			marker.Height,
		)
	}

	// Store the accounts
	err := db.SaveAccounts(accounts)
	if err != nil {
		return fmt.Errorf("error while storing markers accounts: %s", err)
	}

	// Store the markers accounts
	stmt = stmt[:len(stmt)-1] // Remove trailing ","
	stmt += " ON CONFLICT DO NOTHING"
	_, err = db.Sql.Exec(stmt, markerParams...)
	if err != nil {
		return fmt.Errorf("error while storing markers list: %s", err)
	}

	return nil

}
