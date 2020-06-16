package database

import (
	"database/sql/driver"
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/lib/pq"
)

type AccountRow struct {
	Address   string    `db:"address"`
	Coins     DbCoin    `db:"coins"`
	Height    int64     `db:"height"`
	Timestamp time.Time `db:"height"`
}

// DbCoin represents the information stored inside the database about a single coin
type DbCoin struct {
	Denom  string
	Amount int64
}

// Value implements driver.Valuer
func (coin *DbCoin) Value() (driver.Value, error) {
	return fmt.Sprintf("(%s,%d)", coin.Denom, coin.Amount), nil
}

// DbCoins represents an array of coins
type DbCoins []*DbCoin

// NewDbCoins build a new DbCoins object starting from an array of coins
func NewDbCoins(coins sdk.Coins) DbCoins {
	dbCoins := make([]*DbCoin, 0)
	for _, coin := range coins {
		dbCoins = append(dbCoins, &DbCoin{Amount: coin.Amount.Int64(), Denom: coin.Denom})
	}
	return dbCoins
}

// SaveAccount saves the given account information for the given block height and timestamp
func (db BigDipperDb) SaveAccount(account exported.Account, height int64, timestamp time.Time) error {
	statement := `INSERT INTO account (address, coins, height, timestamp) 
				  VALUES ($1, $2::coin[], $3, $4) 
				  ON CONFLICT DO NOTHING`

	_, err := db.Sql.Exec(statement,
		account.GetAddress().String(), pq.Array(NewDbCoins(account.GetCoins())), height, timestamp)
	return err
}

// SaveAccount saves the given accounts information for the given block height and timestamp
func (db BigDipperDb) SaveAccounts(accounts []exported.Account, height int64, timestamp time.Time) error {
	var insertParams []interface{}

	queryInsert := "INSERT INTO account (address, coins, height, timestamp) VALUES "
	for i, account := range accounts {
		p1 := i * 4 // Starting position for insert params

		queryInsert += fmt.Sprintf("($%d,$%d,$%d,$%d),", p1+1, p1+2, p1+3, p1+4)
		insertParams = append(insertParams,
			account.GetAddress().String(), pq.Array(NewDbCoins(account.GetCoins())), height, timestamp)
	}

	queryInsert = queryInsert[:len(queryInsert)-1] // Remove trailing ","
	queryInsert += " ON CONFLICT DO NOTHING"
	_, err := db.Sql.Exec(queryInsert, insertParams...)
	return err
}

// GetPollByPostID returns the poll row associated to the post having the specified id.
// If the post with the same id has no poll associated to it, nil is returned instead.
func (db BigDipperDb) GetAccounts() ([]sdk.AccAddress, error) {
	sqlStmt := `SELECT DISTINCT address from account`

	var rows []AccountRow
	err := db.sqlx.Select(&rows, sqlStmt)
	if err != nil {
		return nil, err
	}

	addresses := make([]sdk.AccAddress, len(rows))
	for index, row := range rows {
		address, err := sdk.AccAddressFromBech32(row.Address)
		if err != nil {
			return nil, err
		}

		addresses[index] = address
	}

	return addresses, nil
}
