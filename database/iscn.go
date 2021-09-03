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


// func (db *Db) SaveIscnRecord(records types.IscnRecord) error {
// 	iscn_data, err := json.Marshal(&records.Data)
// 	if err != nil {
// 		return err
// 	}

// 	stmt := `
// 	INSERT INTO iscn_record (owner_address, iscn_id, latest_version, ipld, stakeholders, iscn_type, context, record_notes, content_metadata, height)
// 	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
// 	ON CONFLICT (iscn_id) DO UPDATE
// 			SET latest_version = excluded.latest_version,
// 			ipld = excluded.ipld,
// 			stakeholders = excluded.stakeholders,
// 			context = excluded.context,
// 			record_notes = excluded.record_notes,
// 			content_metadata = excluded.content_metadata,
// 			height = excluded.height
// 	WHERE iscn_record.height <= excluded.height`

// 	_, err = db.Sql.Exec(stmt, string(records.Owner), records.IscnId, records.LatestVersion, string(records.Ipld), records.Data.stakeholders, iscn_data, iscn_data["@context"], iscn_data.record_notes, iscn_data.content_metadata, records.Height)
// 	if err != nil {
// 		return err
// 	}
// 	return err
// }

// func (db *Db) UpdateIscnRecordOwnership(records types.IscnRecord) error {
// 	stmt := `
// 	INSERT INTO iscn_record (owner_address, iscn_id, height)
// 	VALUES ($1, $2, $3)
// 	ON CONFLICT (iscn_id) DO UPDATE
// 			SET owner_address = excluded.owner_address,
// 			height = excluded.height
// 	WHERE iscn_record.height <= excluded.height`

// 	_, err := db.Sql.Exec(stmt, string(records.Owner), records.IscnId, records.Height)
// 	if err != nil {
// 		return err
// 	}
// 	return err
// }
	
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