package database

import (
	"fmt"
	"time"
)

// UpdateEnableModules allows to save enabled module into the database
func (db BigDipperDb) InsertEnableModules(modules map[string]bool, time time.Time) error {
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

	values = append(values, time)

	stmt += input
	_, err := db.Sql.Exec(stmt, values...)
	if err != nil {
		return err
	}
	return nil
}
