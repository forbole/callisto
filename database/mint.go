package database

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SaveInflation allows to store the inflation for the given block height as well as timestamp
func (db *BigDipperDb) SaveInflation(inflation sdk.Dec) error {
	stmt := `DELETE FROM inflation WHERE TRUE`
	_, err := db.Sql.Exec(stmt)
	if err != nil {
		return err
	}

	stmt = `INSERT INTO inflation (value) VALUES ($1)`
	_, err = db.Sql.Exec(stmt, inflation.String())
	return err
}
