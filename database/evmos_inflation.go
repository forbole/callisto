package database

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	dbtypes "github.com/forbole/bdjuno/v4/database/types"

	"github.com/forbole/bdjuno/v4/types"
	"github.com/lib/pq"
)

// SaveEvmosInflationParams allows to store the given params inside the database
func (db *Db) SaveEvmosInflationParams(params *types.EvmosInflationParams) error {
	stmt := `
INSERT INTO evmos_inflation_params (params, height) 
VALUES ($1, $2)
ON CONFLICT (one_row_id) DO UPDATE 
    SET params = excluded.params,
        height = excluded.height
WHERE evmos_inflation_params.height <= excluded.height`

	paramsBz, err := json.Marshal(&params.Params)
	if err != nil {
		return fmt.Errorf("error while marshaling evmos inflation params: %s", err)
	}

	_, err = db.SQL.Exec(stmt, string(paramsBz), params.Height)
	if err != nil {
		return fmt.Errorf("error while storing evmos inflation params: %s", err)
	}

	return nil
}

// SaveEvmosInflationData allows to store the given inflation data inside the database
func (db *Db) SaveEvmosInflationData(data *types.EvmosInflationData) error {
	stmt := `
INSERT INTO evmos_inflation_data (circulating_supply, epoch_mint_provision, inflation_rate, inflation_period, skipped_epochs, height) 
VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT (one_row_id) DO UPDATE 
    SET circulating_supply = excluded.circulating_supply,
		epoch_mint_provision = excluded.epoch_mint_provision,
		inflation_rate = excluded.inflation_rate,
		inflation_period = excluded.inflation_period,
		skipped_epochs = excluded.skipped_epochs,
        height = excluded.height
WHERE evmos_inflation_data.height <= excluded.height`

	// Inflation rate is missing '%' sign so it needs to be multiplied by 0.01
	// Example: when inflaionRate is 150, it's actually 1.5 (150%)
	inflationRate := data.InflationRate.Mul(sdk.NewDecWithPrec(1, 2)).String()

	_, err := db.SQL.Exec(stmt,
		pq.Array(dbtypes.NewDbDecCoins(data.CirculatingSupply)),
		pq.Array(dbtypes.NewDbDecCoins(data.EpochMintProvision)),
		inflationRate, data.InflationPeriod, data.SkippedEpochs, data.Height,
	)

	if err != nil {
		return fmt.Errorf("error while storing evmos inflation data: %s", err)
	}

	return nil
}
