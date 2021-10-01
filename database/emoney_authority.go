package database

import (
	"fmt"

	dbtypes "github.com/forbole/bdjuno/database/types"
	"github.com/forbole/bdjuno/types"
	"github.com/lib/pq"
)

// SaveEMoneyGasPrices allows to save the gas prices set by the authority(key)
func (db *Db) SaveEMoneyGasPrices(gasPrices types.EMoneyGasPrices) error {
	stmt := `
INSERT INTO emoney_gas_prices (gas_prices, height) 
VALUES ($1, $2) 
ON CONFLICT (one_row_id) DO UPDATE 
    SET gas_prices = excluded.gas_prices,
		height = excluded.height
WHERE emoney_gas_prices.height <= excluded.height`

	_, err := db.Sql.Exec(stmt, pq.Array(dbtypes.NewDbDecCoins(gasPrices.GasPrices)), gasPrices.Height)
	if err != nil {
		return fmt.Errorf("error while storing e-Money gas prices: %s", err)
	}

	return nil
}
