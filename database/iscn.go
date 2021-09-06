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
	ON CONFLICT (iscn_id) DO UPDATE
			SET latest_version = excluded.latest_version,
			ipld = excluded.ipld,
			iscn_data = excluded.iscn_data,
			height = excluded.height
	WHERE iscn_record.height <= excluded.height`

	_, err = db.Sql.Exec(stmt, string(records.Owner), records.IscnId, records.LatestVersion, string(records.Ipld), string(iscn_data), records.Height)
	if err != nil {
		return err
	}
	return err
}


func (db *Db) UpdateIscnRecord(records types.IscnRecord) error {
	iscn_data, err := json.Marshal(&records.Data)
	if err != nil {
		return err
	}

	stmt := `
	UPDATE iscn_record (owner_address, latest_version, ipld, iscn_data, height)
	VALUES ($1, $2, $3, $4, $5)
	WHERE iscn_id = $6
	ON CONFLICT DO UPDATE`

	_, err = db.Sql.Exec(stmt, string(records.Owner), records.LatestVersion, string(records.Ipld), string(iscn_data), records.Height, records.IscnId)
	if err != nil {
		return err
	}
	return err
}


func (db *Db) UpdateIscnRecordOwnership(records types.IscnChangeOwnership) error {
	
	stmt := `UPDATE iscn_record SET owner_address = $1 where iscn_id = $2`

	_, err := db.Sql.Exec(stmt, string(records.NewOwner), records.IscnId)
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
