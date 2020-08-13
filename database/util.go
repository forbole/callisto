package database

import (
	"fmt"
	"time"
)

// UpdateEnableModules allows to save enabled module into the database
func (db BigDipperDb) InsertEnableModules(modules map[string]bool) error {
	stmt := `INSERT INTO enabled_modules values(`
	var values string
	for key, value := range modules {
		stmt += fmt.Sprintf("%s,", key)
		values += fmt.Sprintf("%t,", value)
	}
	stmt += "timestamp) VALUES ("
	values += time.Now().String()
	values += ")"
	stmt += values
	_, err := db.Sql.Exec(stmt)
	if err != nil {
		return err
	}
	return nil
}
