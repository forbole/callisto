package database

import (
	"encoding/json"
	"fmt"

	"github.com/forbole/bdjuno/types"

)

func (db *Db) SaveRecord(records []types.IscnRecord) error {
	if len(records) == 0 {
		return nil
	}

	stmt := `INSERT INTO iscn_record(owner, latest_version, records, height) VALUES `
	var recordList []interface{}

	for i, record := range records {
		vi := i * 4
		stmt += fmt.Sprintf("($%d,$%d,$%d,$%d),", vi+1, vi+2, vi+3, vi+4)
		recordList = append(recordList, record.Owner, record.LatestVersion, record.Records, record.Height)
	}

	stmt = stmt[:len(stmt)-1] // Remove trailing ,
	stmt += `
	ON CONFLICT ON CONSTRAINT one_row_id) DO UPDATE 
		SET records = excluded.records,
			height = excluded.height
	WHERE iscn_record.height <= excluded.height`
	_, err := db.Sql.Exec(stmt, recordList...)
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
