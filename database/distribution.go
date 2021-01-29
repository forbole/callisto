package database

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lib/pq"

	dbtypes "github.com/forbole/bdjuno/database/types"
	bdistrtypes "github.com/forbole/bdjuno/x/distribution/types"
)

// SaveCommunityPool allows to save for the given height the given total amount of coins
func (db *BigDipperDb) SaveCommunityPool(coin sdk.DecCoins, height int64) error {
	query := `INSERT INTO community_pool(coins, height) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	_, err := db.Sql.Exec(query, pq.Array(dbtypes.NewDbDecCoins(coin)), height)
	return err
}

// SaveValidatorCommissionAmounts saves the given validator commission amounts for the given height
func (db *BigDipperDb) SaveValidatorCommissionAmounts(amounts []bdistrtypes.ValidatorCommissionAmount, height int64) error {
	stmt := `INSERT INTO validator_commission_amount(validator_address, amount, height) VALUES `
	var params []interface{}

	for i, amount := range amounts {
		ai := i * 3
		stmt += fmt.Sprintf("($%d, $%d, $%d),", ai+1, ai+2, ai+3)
		params = append(params, amount.ValidatorAddress, pq.Array(dbtypes.NewDbDecCoins(amount.Amount)), height)
	}

	stmt = stmt[:len(stmt)-1] // Remove trailing ,
	stmt += "ON CONFLICT DO NOTHING"
	_, err := db.Sql.Exec(stmt, params...)
	return err
}

// SaveDelegatorsCommissionAmounts saves the given delegator commission amounts for the provided height
func (db *BigDipperDb) SaveDelegatorsCommissionAmounts(amounts []bdistrtypes.DelegatorCommissionAmount, height int64) error {
	stmt := `INSERT INTO delegation_reward(validator_address, delegator_address, amount, height) VALUES `
	var params []interface{}

	for i, amount := range amounts {
		ai := i * 4
		stmt += fmt.Sprintf("($%d, $%d, $%d, $%d),", ai+1, ai+2, ai+3, ai+4)

		coins := pq.Array(dbtypes.NewDbDecCoins(amount.Amount))
		params = append(params,
			amount.ValidatorAddress, amount.DelegatorAddress, coins, height)
	}

	stmt = stmt[:len(stmt)-1] // Remove trailing ,
	stmt += "ON CONFLICT DO NOTHING"
	_, err := db.Sql.Exec(stmt, params...)
	return err
}
