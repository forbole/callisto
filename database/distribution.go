package database

import (
	"fmt"

	dbtypes "github.com/forbole/bdjuno/database/types"

	"github.com/forbole/bdjuno/types"

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
	return err
}

// -------------------------------------------------------------------------------------------------------------------

// SaveDistributionParams allows to store the given distribution parameters inside the database
func (db *Db) SaveDistributionParams(params types.DistributionParams) error {
	stmt := `
INSERT INTO distribution_params (community_tax, base_proposer_reward, bonus_proposer_reward, withdraw_address_enabled, height) 
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (one_row_id) DO UPDATE 
    SET community_tax = excluded.community_tax,
      	base_proposer_reward = excluded.base_proposer_reward,
      	bonus_proposer_reward = excluded.bonus_proposer_reward,
      	withdraw_address_enabled = excluded.withdraw_address_enabled,
      	height = excluded.height
WHERE distribution_params.height <= excluded.height`
	_, err := db.Sql.Exec(stmt,
		params.CommunityTax.String(), params.BaseProposerReward.String(), params.BonusProposerReward.String(),
		params.WithdrawAddrEnabled, params.Height)
	return err
}

// -------------------------------------------------------------------------------------------------------------------

// SaveValidatorCommissionAmount allows to store the given validator commission amount as the most updated one
func (db *Db) SaveValidatorCommissionAmount(amount types.ValidatorCommissionAmount) error {

	consAddr, err := db.GetValidatorConsensusAddress(amount.ValidatorOperAddr)
	if err != nil {
		return err
	}

	err = db.storeUpToDateValidatorCommissionAmount(amount, consAddr)
	if err != nil {
		return fmt.Errorf("error while storing up-to-date validator commission amount: %s", err)
	}

	if db.IsStoreHistoricDataEnabled() {
		err = db.storeHistoricValidatorCommissionAmount(amount, consAddr)
		if err != nil {
			return fmt.Errorf("error while storing validator commission amount history: %s", err)
		}
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
	return err
}

// storeHistoricValidatorCommissionAmount allows to store the given amount has an historic one
func (db *Db) storeHistoricValidatorCommissionAmount(amount types.ValidatorCommissionAmount, consAddr sdk.ConsAddress) error {
	stmt := `
INSERT INTO validator_commission_amount_history(validator_address, amount, height) 
VALUES ($1, $2, $3) 
ON CONFLICT ON CONSTRAINT validator_commission_amount_history_commission_height_unique DO UPDATE 
    SET amount = excluded.amount`

	_, err := db.Sql.Exec(stmt, consAddr.String(), pq.Array(dbtypes.NewDbDecCoins(amount.Amount)), amount.Height)
	return err
}

// -------------------------------------------------------------------------------------------------------------------

// SaveDelegatorsRewardsAmounts allows to store the given delegator reward amounts as the most updated ones
func (db *Db) SaveDelegatorsRewardsAmounts(amounts []types.DelegatorRewardAmount) error {
	err := db.storeUpToDateDelegatorsRewardsAmounts(amounts)
	if err != nil {
		return fmt.Errorf("error while storing up-to-date delegator rewards amounts: %s", err)
	}

	if db.IsStoreHistoricDataEnabled() {
		err = db.storeDelegatorsRewardsAmountsHistory(amounts)
		if err != nil {
			return fmt.Errorf("error while storing delegator rewards amounts history: %s", err)
		}
	}

	return nil
}

// storeUpToDateDelegatorsRewardsAmounts allows to store the given amounts has the most up-to-date ones
func (db *Db) storeUpToDateDelegatorsRewardsAmounts(amounts []types.DelegatorRewardAmount) error {
	stmt := `INSERT INTO delegation_reward(validator_address, delegator_address, withdraw_address, amount, height) VALUES `
	var params []interface{}

	for i, amount := range amounts {
		ai := i * 5
		stmt += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d),", ai+1, ai+2, ai+3, ai+4, ai+5)

		// Get the validator consensus address
		consAddr, err := db.GetValidatorConsensusAddress(amount.ValidatorOperAddr)
		if err != nil {
			return err
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
	return err
}

// storeDelegatorsRewardsAmountsHistory allows to store the given amounts as historic rewards amounts
func (db *Db) storeDelegatorsRewardsAmountsHistory(amounts []types.DelegatorRewardAmount) error {
	stmt := `
INSERT INTO delegation_reward_history 
    (validator_address, delegator_address, withdraw_address, amount, height) 
VALUES `
	var params []interface{}

	for i, amount := range amounts {
		ai := i * 5
		stmt += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d),", ai+1, ai+2, ai+3, ai+4, ai+5)

		// Get the validator consensus address
		consAddr, err := db.GetValidatorConsensusAddress(amount.ValidatorOperAddr)
		if err != nil {
			return err
		}

		coins := pq.Array(dbtypes.NewDbDecCoins(amount.Amount))
		params = append(params,
			consAddr.String(), amount.DelegatorAddress, amount.WithdrawAddress, coins, amount.Height)
	}

	stmt = stmt[:len(stmt)-1] // Remove trailing ,
	stmt += `
ON CONFLICT ON CONSTRAINT delegation_reward_history_validator_delegator_unique DO UPDATE 
	SET withdraw_address = excluded.withdraw_address,
		amount = excluded.amount`
	_, err := db.Sql.Exec(stmt, params...)
	return err
}
