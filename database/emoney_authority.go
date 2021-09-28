package database

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	dbtypes "github.com/forbole/bdjuno/database/types"
	"github.com/forbole/bdjuno/types"
	"github.com/lib/pq"
)

// SaveEMoneyGasPrices allows to save the gas prices set by the authority(key)
func (db *Db) SaveEMoneyGasPrices(gasPrices types.EMoneyGasPrices) error {
	stmt := `
INSERT INTO emoney_gas_prices (authority_key, gas_prices, height) 
VALUES ($1, $2, $3) 
ON CONFLICT (authority_key) DO UPDATE 
    SET gas_prices = excluded.gas_prices,
		height = excluded.height
WHERE emoney_gas_prices.height <= excluded.height`

	_, err := db.Sql.Exec(stmt, gasPrices.AuthorityKey, pq.Array(dbtypes.NewDbDecCoins(gasPrices.GasPrices)), gasPrices.Height)
	if err != nil {
		return fmt.Errorf("error while storing e-Money gas prices: %s", err)
	}

	return nil
}

// UpdateEMoneyGasPrices allows to update the gas prices obtained from periodic operations
func (db *Db) UpdateEMoneyGasPrices(gasPrices sdk.DecCoins, height int64) error {
	stmt := `UPDATE emoney_gas_prices SET gas_prices = $1, height = $2`

	_, err := db.Sql.Exec(stmt, pq.Array(dbtypes.NewDbDecCoins(gasPrices)), height)
	if err != nil {
		return fmt.Errorf("error while updating e-Money gas prices: %s", err)
	}

	return nil
}
