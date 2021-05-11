package bigdipper

import (
	"fmt"

	dbtypes "github.com/forbole/bdjuno/database/types"

	"github.com/forbole/bdjuno/modules/common/distribution"

	"github.com/forbole/bdjuno/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/lib/pq"
)

var (
	_ distribution.DB = &Db{}
)

// GetValidatorsInfo implements distribution.DB
func (db *Db) GetValidatorsInfo() ([]types.ValidatorInfo, error) {
	validators, err := db.GetValidators()
	if err != nil {
		return nil, err
	}

	var valInfo = make([]types.ValidatorInfo, len(validators))
	for index, val := range validators {
		valInfo[index] = types.NewValidatorInfo(val.ValAddress, val.SelfDelegateAddress)
	}

	return valInfo, nil
}

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

// SaveValidatorCommissionAmount implements distribution.DB
func (db *Db) SaveValidatorCommissionAmount(amount types.ValidatorCommissionAmount) error {
	stmt := `
INSERT INTO validator_commission_amount(validator_address, amount, height) 
VALUES ($1, $2, $3) 
ON CONFLICT (validator_address) DO UPDATE 
    SET amount = excluded.amount, 
        height = excluded.height
WHERE validator_commission_amount.height <= excluded.height`

	consAddr, err := db.GetValidatorConsensusAddress(amount.ValidatorOperAddr)
	if err != nil {
		return err
	}

	_, err = db.Sql.Exec(stmt, consAddr.String(), pq.Array(dbtypes.NewDbDecCoins(amount.Amount)), amount.Height)
	return err
}

// SaveDelegatorsRewardsAmounts implements distribution.DB
func (db *Db) SaveDelegatorsRewardsAmounts(amounts []types.DelegatorRewardAmount) error {
	stmt := `INSERT INTO delegation_reward(validator_address, delegator_address, withdraw_address, amount, height) VALUES `
	var params []interface{}

	for i, amount := range amounts {
		ai := i * 5
		stmt += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d),", ai+1, ai+2, ai+3, ai+4, ai+5)

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
ON CONFLICT ON CONSTRAINT validator_delegator_unique DO UPDATE 
	SET withdraw_address = excluded.withdraw_address,
		amount = excluded.amount,
		height = excluded.height
WHERE delegation_reward.height <= excluded.height`
	_, err := db.Sql.Exec(stmt, params...)
	return err
}
