package database

import (
	"encoding/json"
	"fmt"

	dbtypes "github.com/forbole/bdjuno/v2/database/types"

	"github.com/forbole/bdjuno/v2/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lib/pq"
)

// SaveCommunityPool allows to save for the given height the given total amount of coins
func (db *Db) SaveCommunityPool(coin sdk.DecCoins, height int64) error {
	query := `
INSERT INTO community_pool(coins, height) 
VALUES ($1, $2) 
ON CONFLICT (one_row_id) DO UPDATE 
    SET coins = excluded.coins,
        height = excluded.height
WHERE community_pool.height <= excluded.height`
	_, err := db.Sql.Exec(query, pq.Array(dbtypes.NewDbDecCoins(coin)), height)
	if err != nil {
		return fmt.Errorf("error while storing community pool: %s", err)
	}

	return nil
}

// -------------------------------------------------------------------------------------------------------------------

// SaveDistributionParams allows to store the given distribution parameters inside the database
func (db *Db) SaveDistributionParams(params *types.DistributionParams) error {
	paramsBz, err := json.Marshal(&params.Params)
	if err != nil {
		return fmt.Errorf("error while marshaling params: %s", err)
	}

	stmt := `
INSERT INTO distribution_params (params, height) 
VALUES ($1, $2)
ON CONFLICT (one_row_id) DO UPDATE 
    SET params = excluded.params,
      	height = excluded.height
WHERE distribution_params.height <= excluded.height`
	_, err = db.Sql.Exec(stmt, string(paramsBz), params.Height)
	if err != nil {
		return fmt.Errorf("error while storing distribution params: %s", err)
	}

	return nil
}

// -------------------------------------------------------------------------------------------------------------------

// SaveValidatorCommissionAmount allows to store the given validator commission amount as the most updated one
func (db *Db) SaveValidatorCommissionAmount(amount types.ValidatorCommissionAmount) error {

	consAddr, err := db.GetValidatorConsensusAddress(amount.ValidatorOperAddr)
	if err != nil {
		return fmt.Errorf("error while getting validator consensus address: %s", err)
	}

	err = db.storeUpToDateValidatorCommissionAmount(amount, consAddr)
	if err != nil {
		return fmt.Errorf("error while storing up-to-date validator commission amount: %s", err)
	}

	return nil
}

// storeUpToDateValidatorCommissionAmount allows to store the given amount as the most up-to-date one
func (db *Db) storeUpToDateValidatorCommissionAmount(amount types.ValidatorCommissionAmount, consAddr sdk.ConsAddress) error {
	stmt := `
INSERT INTO validator_commission_amount(validator_address, amount, height) 
VALUES ($1, $2, $3) 
ON CONFLICT (validator_address) DO UPDATE 
    SET amount = excluded.amount, 
        height = excluded.height
WHERE validator_commission_amount.height <= excluded.height`

	_, err := db.Sql.Exec(stmt, consAddr.String(), pq.Array(dbtypes.NewDbDecCoins(amount.Amount)), amount.Height)
	if err != nil {
		return fmt.Errorf("error while storing validator commission amount: %s", err)
	}

	return nil
}

// GetUserValidatorCommissionAmount returns the current amount of the validator commission
// that is associated with the user having the given address
func (db *Db) GetUserValidatorCommissionAmount(address string) (sdk.DecCoins, error) {
	stmt := `
SELECT validator_commission_amount.*
FROM validator_commission_amount 
    INNER JOIN validator_info vi on validator_commission_amount.validator_address = vi.consensus_address
WHERE vi.self_delegate_address = $1`

	var rows []*dbtypes.ValidatorCommissionAmountRow
	err := db.Sqlx.Select(&rows, stmt, address)
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return sdk.DecCoins{}, nil
	}

	return rows[0].Amount.ToDecCoins(), nil
}

// -------------------------------------------------------------------------------------------------------------------

// SaveDelegatorsRewardsAmounts allows to store the given delegator reward amounts as the most updated ones
func (db *Db) SaveDelegatorsRewardsAmounts(amounts []types.DelegatorRewardAmount) error {
	if len(amounts) == 0 {
		return nil
	}

	err := db.storeUpToDateDelegatorsRewardsAmounts(amounts)
	if err != nil {
		return fmt.Errorf("error while storing up-to-date delegator rewards amounts: %s", err)
	}

	return nil
}

// storeUpToDateDelegatorsRewardsAmounts allows to store the given amounts has the most up-to-date ones
func (db *Db) storeUpToDateDelegatorsRewardsAmounts(amounts []types.DelegatorRewardAmount) error {
	if len(amounts) == 0 {
		return nil
	}

	stmt := `INSERT INTO delegation_reward(validator_address, delegator_address, withdraw_address, amount, height) VALUES `
	var params []interface{}

	for i, amount := range amounts {
		ai := i * 5
		stmt += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d),", ai+1, ai+2, ai+3, ai+4, ai+5)

		// Get the validator consensus address
		consAddr, err := db.GetValidatorConsensusAddress(amount.ValidatorOperAddr)
		if err != nil {
			return fmt.Errorf("error while getting validator consensus address: %s", err)
		}

		coins := pq.Array(dbtypes.NewDbDecCoins(amount.Amount))
		params = append(params,
			consAddr.String(), amount.DelegatorAddress, amount.WithdrawAddress, coins, amount.Height)
	}

	stmt = stmt[:len(stmt)-1] // Remove trailing ,
	stmt += `
ON CONFLICT ON CONSTRAINT delegation_reward_validator_delegator_unique DO UPDATE 
	SET withdraw_address = excluded.withdraw_address,
		amount = excluded.amount,
		height = excluded.height
WHERE delegation_reward.height <= excluded.height`
	_, err := db.Sql.Exec(stmt, params...)
	if err != nil {
		return fmt.Errorf("error while storing delegation reward: %s", err)
	}

	return nil
}

// GetUserDelegatorRewardsAmount returns the amount of rewards that the given user has currently associated
func (db *Db) GetUserDelegatorRewardsAmount(address string) (sdk.DecCoins, error) {
	stmt := `SELECT * FROM delegation_reward WHERE delegator_address = $1`

	var rows []*dbtypes.DelegationRewardRow
	err := db.Sqlx.Select(&rows, stmt, address)
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return sdk.DecCoins{}, nil
	}

	var rewardsAmount = sdk.DecCoins{}
	for _, row := range rows {
		rewardsAmount = rewardsAmount.Add(row.Amount.ToDecCoins()...)
	}
	return rewardsAmount, nil
}

// DeleteDelegatorRewardsAmount deletes the rewards for the given delegator-validator tuple at the given height
func (db *Db) DeleteDelegatorRewardsAmount(delegatorAddr string, height int64) error {
	stmt := `DELETE FROM delegation_reward WHERE delegator_address = $1 AND height <= $2`
	_, err := db.Sql.Exec(stmt, delegatorAddr, height)
	if err != nil {
		return fmt.Errorf("error while deleting delegation reward: %s", err)
	}

	return nil
}

// HasDelegatorRewards checks if the database contains any delegation reward
func (db *Db) HasDelegatorRewards() (bool, error) {
	stmt := `SELECT COUNT(*) FROM delegation_reward`

	var count int
	err := db.Sql.QueryRow(stmt).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("error while getting delegation reward: %s", err)

	}

	return count > 0, nil
}

// HasValidatorCommission checks if the database contains any validators commission
func (db *Db) HasValidatorCommission() (bool, error) {
	stmt := `SELECT COUNT(*) FROM validator_commission_amount`

	var count int
	err := db.Sql.QueryRow(stmt).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("error while getting validator commission: %s", err)

	}

	return count > 0, nil
}
