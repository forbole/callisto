import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"strings"
	authtype "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
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

// Period represents the information stored inside the database about 
// a length of time and amount of coins that will vest.
type DbPeriod struct {
	Length  int64
	Amount DbCoins
}

// NewPeriod builds a Period starting from an SDK Period
func NewDbPeriod(period authtype.DbPeriod) Period {
	dbcoins = NewDbCoins(period.amount)
	return DbCoin{
		Length:  period.Length,
		Amount: dbcoins,
	}
}

// Equal tells whether coin and d represent the same coin with the same amount
func (period DbPeriod) Equal(d Period) bool {
	return period.Length == d.Length && period.Amount.Equal(d.Amount)
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
func NewDbPeriods(periods authtype.Periods) s {
	dbPeriods := make([]*DbPeriod, 0)
	for _, period := range periods {
		dbPeriods = append(dbPeriods, &DbPeriod{Amount: period.Amount.String(), Length: period.Length})
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

	for index, coin := range coins {
		if !coin.Equal(*(*d)[index]) {
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
