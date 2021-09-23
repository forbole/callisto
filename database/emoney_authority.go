package database

import (
	"fmt"

	dbtypes "github.com/forbole/bdjuno/database/types"
	"github.com/forbole/bdjuno/types"
	"github.com/lib/pq"
)

// SaveSupply allows to save for the given height the given total amount of coins
func (db *Db) SaveGasPrices(gasPrice types.EmoneyGasPrice) error {
	stmt := `
INSERT INTO emoney_gas_prices (authority, gas_prices, height) 
VALUES ($1, $2, $3) 
ON CONFLICT (authority) DO UPDATE 
    SET gas_prices = excluded.gas_prices,
		height = excluded.height,
WHERE emoney_gas_prices.height <= excluded.height`

	_, err := db.Sql.Exec(stmt, gasPrice.Authority, pq.Array(dbtypes.NewDbDecCoins(gasPrice.GasPrices)), gasPrice.Height)
	if err != nil {
		return fmt.Errorf("error while storing emoney gas prices: %s", err)
	}

	return nil
}
