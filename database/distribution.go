package database

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lib/pq"

	dbtypes "github.com/forbole/bdjuno/database/types"
	bdistrtypes "github.com/forbole/bdjuno/x/distribution/types"
)

// SaveCommunityPool allows to save for the given height the given total amount of coins
func (db *BigDipperDb) SaveCommunityPool(coin sdk.DecCoins) error {
	stmt := `DELETE FROM community_pool WHERE TRUE`
	_, err := db.Sql.Exec(stmt)
	if err != nil {
		return err
	}

	stmt = `INSERT INTO community_pool(coins) VALUES ($1)`
	_, err = db.Sql.Exec(stmt, pq.Array(dbtypes.NewDbDecCoins(coin)))
	return err
}

// SaveValidatorCommissionAmount saves the given validator commission amounts for the given height
func (db *BigDipperDb) SaveValidatorCommissionAmount(amount bdistrtypes.ValidatorCommissionAmount) error {
	stmt := `
INSERT INTO validator_commission_amount(validator_address, amount) 
VALUES ($1, $2) 
ON CONFLICT (validator_address) DO UPDATE 
    SET amount = excluded.amount`

	_, err := db.Sql.Exec(stmt,
		amount.ValidatorConsAddress, pq.Array(dbtypes.NewDbDecCoins(amount.Amount)))
	return err
}

// SaveDelegatorsRewardsAmounts saves the given delegator commission amounts for the provided height
func (db *BigDipperDb) SaveDelegatorsRewardsAmounts(amounts []bdistrtypes.DelegatorReward) error {
	stmt := `INSERT INTO delegation_reward(validator_address, delegator_address, withdraw_address, amount) VALUES `
	var params []interface{}

	for i, amount := range amounts {
		ai := i * 4
		stmt += fmt.Sprintf("($%d, $%d, $%d, $%d),", ai+1, ai+2, ai+3, ai+4)

		coins := pq.Array(dbtypes.NewDbDecCoins(amount.Amount))
		params = append(params,
			amount.ValidatorConsAddress, amount.DelegatorAddress, amount.WithdrawAddress, coins)
	}

	stmt = stmt[:len(stmt)-1] // Remove trailing ,
	stmt += `ON CONFLICT ON CONSTRAINT validator_delegator_unique DO UPDATE 
SET amount = excluded.amount, 
    withdraw_address = excluded.withdraw_address`

	_, err := db.Sql.Exec(stmt, params...)
	return err
}
