package database

import (
	"database/sql"
	"fmt"
)

type DbTx struct {
	*sql.Tx
}

func (db *Db) ExecuteTx(callback func(*DbTx) error) error {
	tx, err := db.Sqlx.Begin()
	if err != nil {
		return err
	}

	if err := callback(&DbTx{tx}); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return fmt.Errorf("error while executing tx: %s\nerror while making rollback: %s", err, rollbackErr)
		}
		return err
	}

	return tx.Commit()
}
