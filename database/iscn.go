package database

import (
	"encoding/json"

	"github.com/forbole/bdjuno/types"
)

func (db *Db) SaveIscnRecord(records types.IscnRecord) error {
	iscn_data, err := json.Marshal(&records.Data)
	if err != nil {
		return err
	}

	stmt := `
	INSERT INTO iscn_record (owner_address, iscn_id, latest_version, ipld, iscn_data, height)
	VALUES ($1, $2, $3, $4, $5, $6)
	ON CONFLICT DO NOTHING`

	_, err = db.Sql.Exec(stmt, string(records.Owner), records.IscnId, records.LatestVersion, string(records.Ipld), string(iscn_data), records.Height)
	if err != nil {
		return err
	}
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