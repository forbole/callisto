package types

import (
	"database/sql/driver"
	"fmt"

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
