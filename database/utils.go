package database

import (
	"fmt"
)

// UpdateEnableModules allows to save enabled module into the database
func (db *BigDipperDb) InsertEnableModules(modules []string) error {
	//clear table first
	stmt := "DELETE FROM modules WHERE TRUE"
	_, err := db.Sql.Exec(stmt)
	if err != nil {
		return err
	}

	if len(modules) == 0 {
		return nil
	}

	var values []interface{}
	stmt = `INSERT INTO modules (module_name) VALUES`
	for key, value := range modules {
		stmt += fmt.Sprintf("($%d),", key+1)
		values = append(values, value)
	}
	stmt = stmt[:len(stmt)-1] //remove tailing ","
	stmt += " ON CONFLICT DO NOTHING"
	_, err = db.Sql.Exec(stmt, values...)
	if err != nil {
		return err
	}
	return nil
}
