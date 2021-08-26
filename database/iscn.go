package database

import (
	"encoding/json"

	"github.com/forbole/bdjuno/types"
	iscntypes "github.com/likecoin/likechain/x/iscn/types"

)

// SaveRecord allows to store iscn record for the given block height and timestamp
func (db *Db) SaveRecord(owner string, latestVersion uint64, records []iscntypes.QueryResponseRecord, height int64) error {
	stmt := `INSERT INTO iscn_record(owner, latestVersion, records, height) 
VALUES ($1, $2, $3, $4)
ON CONFLICT (one_row_id) DO UPDATE 
	SET owner = excluded.owner, latest_version = excluded.latestVersion, records = excluded.records,
		height = excluded.height
WHERE iscn_record.height <= excluded.height`
	_, err := db.Sql.Exec(stmt, string(owner), uint64(latestVersion), records, height)
	return err
}

// SaveIscnParams allows to store iscn params inside the database
func (db *Db) SaveIscnParams(params types.IscnParams) error {
	paramsBz, err := json.Marshal(&params.Params)
	if err != nil {
		return err
	}

	stmt := `
INSERT INTO iscn_params (params, height) 
VALUES ($1, $2)
ON CONFLICT (one_row_id) DO UPDATE 
    SET params = excluded.params,
        height = excluded.height
WHERE iscn_params.height <= excluded.height`
	_, err = db.Sql.Exec(stmt, string(paramsBz), params.Height)
	return err
}
