package database

import (
	"fmt"
	"time"
)

// UpdateEnableModules allows to save enabled module into the database
func (db BigDipperDb) InsertEnableModules(modules map[string]bool) error {
	stmt := `INSERT INTO modules (`
	input := `(`
	var values []interface{}
	count := 1
	for key, value := range modules {
		stmt += fmt.Sprintf("%s,", key)
		input += fmt.Sprintf("$%d,", count)
		values = append(values, value)
		count++
	}
	stmt += "timestamp) VALUES "
	input += fmt.Sprintf("$%d)", count)

	now, err := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	if err != nil {
		return err
	}
	values = append(values, now)
	
	stmt += input
	_, err = db.Sql.Exec(stmt, values...)
	if err != nil {
		return err
	}
	return nil
}
