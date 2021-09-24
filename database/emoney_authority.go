package database

import (
	"fmt"

	dbtypes "github.com/forbole/bdjuno/database/types"
	"github.com/forbole/bdjuno/types"
	"github.com/lib/pq"
)

// SaveEmoneyGasPrices allows to save the gas prices set by the authority(key)
func (db *Db) SaveEmoneyGasPrices(gasPrice types.EmoneyGasPrice) error {
	stmt := `
INSERT INTO emoney_gas_prices (authority_key, gas_prices, height) 
VALUES ($1, $2, $3) 
ON CONFLICT (authority_key) DO UPDATE 
    SET gas_prices = excluded.gas_prices,
		height = excluded.height,
WHERE emoney_gas_prices.height <= excluded.height`

	fmt.Println("SaveEmoneyGasPrices Executed")
	_, err := db.Sql.Exec(stmt, gasPrice.AuthorityKey, pq.Array(dbtypes.NewDbDecCoins(gasPrice.GasPrices)), gasPrice.Height)
	if err != nil {
		return fmt.Errorf("error while storing emoney gas prices: %s", err)
	}

	return nil
}
