package database

import (
	"encoding/json"
	"fmt"

	dbtypes "github.com/forbole/bdjuno/v3/database/types"
	dbutils "github.com/forbole/bdjuno/v3/database/utils"
	"github.com/forbole/bdjuno/v3/types"
	"github.com/lib/pq"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SaveMarketParams allows to store the given params inside the database
func (db *Db) SaveMarketParams(params *types.MarketParams) error {
	paramsBz, err := json.Marshal(&params.Params)
	if err != nil {
		return fmt.Errorf("error while marshaling market params: %s", err)
	}

	stmt := `
INSERT INTO market_params (params, height) 
VALUES ($1, $2)
ON CONFLICT (one_row_id) DO UPDATE 
    SET params = excluded.params,
        height = excluded.height
WHERE market_params.height <= excluded.height`

	_, err = db.Sql.Exec(stmt, string(paramsBz), params.Height)
	if err != nil {
		return fmt.Errorf("error while storing market params: %s", err)
	}

	return nil
}

// SaveLease allows to store a single x/market lease inside the database
func (db *Db) SaveLease(lease *types.MarketLease, height int64) error {
	return db.SaveLeases([]*types.MarketLease{lease}, height)
}

// SaveLeases allows to store the given x/market leases inside the database
func (db *Db) SaveLeases(leases []*types.MarketLease, height int64) error {
	paramsCount := 10
	slices := dbutils.SplitLeases(leases, paramsCount)

	for _, leases := range slices {
		if len(leases) == 0 {
			continue
		}

		// Store up-to-date data
		err := db.saveLeases(paramsCount, leases, height)
		if err != nil {
			return fmt.Errorf("error while storing akash x/market leases: %s", err)
		}
	}

	return nil
}

func (db *Db) saveLeases(paramsCount int, leases []*types.MarketLease, height int64) error {
	stmt := `INSERT INTO akash_lease (owner, d_seq, g_seq, o_seq, provider, lease_state, price, created_at, closed_on, height) VALUES `

	var params []interface{}
	for i, lease := range leases {
		ii := i * paramsCount

		stmt += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d),",
			ii+1, ii+2, ii+3, ii+4, ii+5, ii+6, ii+7, ii+8, ii+9, ii+10)

		params = append(params,
			lease.Owner, lease.DSeq, lease.GSeq, lease.OSeq, lease.Provider,
			lease.State,
			pq.Array(dbtypes.NewDbDecCoins([]sdk.DecCoin{lease.Price})),
			lease.CreatedAt,
			lease.ClosedOn,
			height,
		)

	}
	stmt = stmt[:len(stmt)-1]
	stmt += ` ON CONFLICT (owner, d_seq, g_seq, o_seq, provider) DO UPDATE 
    SET lease_state = excluded.lease_state,
		price = excluded.price,
		created_at = excluded.created_at,
		closed_on = excluded.closed_on,
    	height = excluded.height 
WHERE akash_lease.height <= excluded.height`

	fmt.Println("params lens: ", len(params))

	_, err := db.Sql.Exec(stmt, params...)
	if err != nil {
		return fmt.Errorf("error while storing akash leases: %s", err)
	}

	return nil
}
