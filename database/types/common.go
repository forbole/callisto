package types

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func ToString(value sql.NullString) string {
	if value.Valid {
		return value.String
	}
	return ""
}

func ToNullString(value string) sql.NullString {
	value = strings.TrimSpace(value)
	return sql.NullString{
		Valid:  value != "",
		String: value,
	}
}

// _________________________________________________________

// DbCoin represents the information stored inside the database about a single coin
type DbCoin struct {
	Denom  string
	Amount string
}

// NewCoin builds a DbCoin starting from an SDK Coin
func NewDbCoin(coin sdk.Coin) DbCoin {
	return DbCoin{
		Denom:  coin.Denom,
		Amount: coin.Amount.String(),
	}
}

// Equal tells whether coin and d represent the same coin with the same amount
func (coin DbCoin) Equal(d DbCoin) bool {
	return coin.Denom == d.Denom && coin.Amount == d.Amount
}

// Value implements driver.Valuer
func (coin *DbCoin) Value() (driver.Value, error) {
	return fmt.Sprintf("(%s,%s)", coin.Denom, coin.Amount), nil
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

	*coin = DbCoin{Denom: values[0], Amount: values[1]}
	return nil
}

// _________________________________________________________

// DbCoins represents an array of coins
type DbCoins []*DbCoin

// NewDbCoins build a new DbCoins object starting from an array of coins
func NewDbCoins(coins sdk.Coins) DbCoins {
	dbCoins := make([]*DbCoin, 0)
	for _, coin := range coins {
		dbCoins = append(dbCoins, &DbCoin{Amount: coin.Amount.String(), Denom: coin.Denom})
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

		coin := DbCoin{Denom: v[0], Amount: v[1]}
		coinsV[index] = &coin
	}

	*coins = coinsV
	return nil
}

//_______________________________________________________

// DbDecCoin represents the information stored inside the database about a single coin
type DbDecCoin struct {
	Denom  string
	Amount string
}

// NewCoin builds a DbDecCoin starting from an SDK Coin
func NewDbDecCoin(coin sdk.DecCoin) DbDecCoin {
	return DbDecCoin{
		Denom:  coin.Denom,
		Amount: coin.Amount.String(),
	}
}

// Equal tells whether coin and d represent the same coin with the same amount
func (coin DbDecCoin) Equal(d DbDecCoin) bool {
	return coin.Denom == d.Denom && coin.Amount == d.Amount
}

// Value implements driver.Valuer
func (coin *DbDecCoin) Value() (driver.Value, error) {
	return fmt.Sprintf("(%s,%s)", coin.Denom, coin.Amount), nil
}

// Scan implements sql.Scanner
func (coin *DbDecCoin) Scan(src interface{}) error {
	strValue := string(src.([]byte))
	strValue = strings.ReplaceAll(strValue, `"`, "")
	strValue = strings.ReplaceAll(strValue, "{", "")
	strValue = strings.ReplaceAll(strValue, "}", "")
	strValue = strings.ReplaceAll(strValue, "(", "")
	strValue = strings.ReplaceAll(strValue, ")", "")

	values := strings.Split(strValue, ",")
	*coin = DbDecCoin{Denom: values[0], Amount: values[1]}
	return nil
}

// _________________________________________________________

// DbDecCoins represents an array of coins
type DbDecCoins []*DbDecCoin

// NewDbDecCoins build a new DbDecCoins object starting from an array of coins
func NewDbDecCoins(coins sdk.DecCoins) DbDecCoins {
	DbDecCoins := make([]*DbDecCoin, 0)
	for _, coin := range coins {
		DbDecCoins = append(DbDecCoins, &DbDecCoin{Amount: coin.Amount.String(), Denom: coin.Denom})
	}
	return DbDecCoins
}

// Equal tells whether c and d contain the same items in the same order
func (coins DbDecCoins) Equal(d *DbDecCoins) bool {
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
func (coins *DbDecCoins) Scan(src interface{}) error {
	strValue := string(src.([]byte))
	strValue = strings.ReplaceAll(strValue, `"`, "")
	strValue = strings.ReplaceAll(strValue, "{", "")
	strValue = strings.ReplaceAll(strValue, "}", "")
	strValue = strings.ReplaceAll(strValue, "),(", ") (")
	strValue = strings.ReplaceAll(strValue, "(", "")
	strValue = strings.ReplaceAll(strValue, ")", "")

	values := strings.Split(strValue, " ")

	coinsV := make(DbDecCoins, len(values))
	for index, value := range values {
		v := strings.Split(value, ",") // Split the values

		coin := DbDecCoin{Denom: v[0], Amount: v[1]}
		coinsV[index] = &coin
	}

	*coins = coinsV
	return nil
}
