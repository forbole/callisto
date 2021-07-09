package types

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"strings"

	authtype "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
)

// _________________________________________________________

// Period represents the information stored inside the database about 
// a length of time and amount of coins that will vest.
type DbPeriod struct {
	Length  int64
	Amount DbCoins
}

// NewPeriod builds a Period starting from an SDK Period
func NewDbPeriod(period authtype.Period) DbPeriod {
	dbcoins := NewDbCoins(period.Amount)
	return DbPeriod{
		Length:  period.Length,
		Amount: dbcoins,
	}
}

// Equal tells whether coin and d represent the same coin with the same amount
func (period DbPeriod) Equal(d DbPeriod) bool {
	return period.Length == d.Length && period.Amount.Equal(&d.Amount)
}

// Value implements driver.Valuer
func (period *DbPeriod) Value() (driver.Value, error) {
	return fmt.Sprintf("(%s,%s)", period.Length, period.Amount), nil
}

// Scan implements sql.Scanner
func (period *DbPeriod) Scan(src interface{}) error {
	/* strValue := string(src.([]byte))
	strValue = strings.ReplaceAll(strValue, `"`, "")
	strValue = strings.ReplaceAll(strValue, "{", "")
	strValue = strings.ReplaceAll(strValue, "}", "")
	strValue = strings.ReplaceAll(strValue, "(", "")
	strValue = strings.ReplaceAll(strValue, ")", "")

	values := strings.Split(strValue, ",")

	*coin = Periods{Length: values[0], Amount: values[1]} */
	return nil
}

// _________________________________________________________

// DbPeriods represents an array of DbPeriod
type DbPeriods []*DbPeriod

// NewDbPeriod build a new DbPeriods object starting from an array of DbPeriod
func NewDbPeriods(periods authtype.Periods) DbPeriods {
	dbPeriods := make([]*DbPeriod, 0)
	for _, period := range periods {
		amount := NewDbCoins(period.Amount) 
		dbPeriods = append(dbPeriods, &DbPeriod{Amount: amount, Length: period.Length})
	}
	return dbPeriods
}

// Equal tells whether c and d contain the same items in the same order
func (periods DbPeriods) Equal(d *DbPeriods) bool {
	if d == nil {
		return false
	}

	if len(periods) != len(*d) {
		return false
	}

	for index, period := range periods {
		if !period.Equal(*(*d)[index]) {
			return false
		}
	}

	return true
}

// Scan implements sql.Scanner
func (coins *DbPeriods) Scan(src interface{}) error {
	/* strValue := string(src.([]byte))
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

	*coins = coinsV */
	return nil
}
