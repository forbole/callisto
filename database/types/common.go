package types

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DbCoin represents the information stored inside the database about a single coin
type DbCoin struct {
	Denom  string
	Amount int64
}

// NewCoin builds a DbCoin starting from an SDK Coin
func NewDbCoin(coin sdk.Coin) DbCoin {
	return DbCoin{
		Denom:  coin.Denom,
		Amount: coin.Amount.Int64(),
	}
}

// Equal tells whether coin and d represent the same coin with the same amount
func (coin DbCoin) Equal(d DbCoin) bool {
	return coin.Denom == d.Denom && coin.Amount == d.Amount
}

// Value implements driver.Valuer
func (coin *DbCoin) Value() (driver.Value, error) {
	return fmt.Sprintf("(%s,%d)", coin.Denom, coin.Amount), nil
}

// Scan implements sql.Scanner
func (coin *DbCoin) Scan(src interface{}) error {
	strValue := string(src.([]byte))
	strValue = strings.ReplaceAll(strValue, `"`, "")
	strValue = strings.ReplaceAll(strValue, "{", "")
	strValue = strings.ReplaceAll(strValue, "}", "")
	strValue = strings.ReplaceAll(strValue, "(", "")
	strValue = strings.ReplaceAll(strValue, ")", "")

	values := strings.Split(strValue, ",")

	amt, err := strconv.ParseInt(values[1], 10, 64)
	if err != nil {
		return err
	}

	*coin = DbCoin{Denom: values[0], Amount: amt}
	return nil
}

// _________________________________________________________

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

// Equal tells whether c and d contain the same items in the same order
func (coins DbCoins) Equal(d *DbCoins) bool {
	if d == nil {
		return false
	}

	if len(coins) != len(*d) {
		return false
	}

	for index, coin := range coins {
		if !coin.Equal(*(*d)[index]) {
			return false
		}
	}

	return true
}

// Scan implements sql.Scanner
func (coins *DbCoins) Scan(src interface{}) error {
	strValue := string(src.([]byte))
	strValue = strings.ReplaceAll(strValue, `"`, "")
	strValue = strings.ReplaceAll(strValue, "{", "")
	strValue = strings.ReplaceAll(strValue, "}", "")
	strValue = strings.ReplaceAll(strValue, "),(", ") (")
	strValue = strings.ReplaceAll(strValue, "(", "")
	strValue = strings.ReplaceAll(strValue, ")", "")

	values := strings.Split(strValue, " ")

	coinsV := make(DbCoins, len(values))
	for index, value := range values {
		v := strings.Split(value, ",") // Split the values

		am, err := strconv.ParseInt(v[1], 10, 64) // Get the value
		if err != nil {
			return err
		}

		coin := DbCoin{Denom: v[0], Amount: am}
		coinsV[index] = &coin
	}

	*coins = coinsV
	return nil
}
