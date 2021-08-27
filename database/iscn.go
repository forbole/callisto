package database

import (
	"encoding/json"

	"github.com/forbole/bdjuno/types"

)

// SaveRecord allows to store iscn record for the given block height and timestamp
func (db *Db) SaveRecord(record types.IscnRecord) error {
	stmt := `
	INSERT INTO iscn_record (ipld, context, record_id, record_route, record_type, content_fingerprints, 
			content_metadata, record_notes, record_timestamp, record_version, stakeholders, height) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) 
	ON CONFLICT ON CONSTRAINT one_row_id DO UPDATE
		SET ipld = excluded.ipld, context = excluded.context, record_id = excluded.record_id, 
			record_route = excluded.record_route, record_type = excluded.record_type, 
			content_fingerprints = excluded.content_fingerprints, content_metadata = excluded.content_metadata, 
			record_notes = excluded.record_notes, record_timestamp = excluded.record_timestamp, 
			record_version = excluded.record_version, stakeholders = excluded.stakeholders,
			height = excluded.height
	WHERE iscn_record.height <= excluded.height`
	_, err := db.Sql.Exec(stmt,
		record.Ipld,
		record.Context,
		record.RecordID,
		record.RecordRoute,
		record.RecordType,
		record.ContentFingerprints,
		record.ContentMetadata,
		record.RecordNotes,
		record.RecordTimestamp,
		record.RecordVersion,
		record.Stakeholders,
		record.Height,
	)
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
