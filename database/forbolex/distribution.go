package forbolex

import (
	"fmt"

	dbtypes "github.com/forbole/bdjuno/database/types"

	"github.com/lib/pq"

	forbolexdbtypes "github.com/forbole/bdjuno/database/forbolex/types"
	"github.com/forbole/bdjuno/modules/common/distribution"
	"github.com/forbole/bdjuno/types"
)

var (
	_ distribution.DB = &Db{}
)

// GetDelegators implements distribution.DB
func (db *Db) GetDelegators() ([]string, error) {
	var rows []string
	err := db.Sqlx.Select(&rows, `SELECT delegator_address FROM delegation_history`)
	return rows, err
}

// GetValidatorsInfo implements distribution.DB
func (db *Db) GetValidatorsInfo() ([]types.ValidatorInfo, error) {
	var rows []forbolexdbtypes.ValidatorInfoRow
	err := db.Sqlx.Select(&rows, `SELECT * FROM validator_info`)
	if err != nil {
		return nil, err
	}

	var vals = make([]types.ValidatorInfo, len(rows))
	for index, row := range rows {
		vals[index] = types.NewValidatorInfo(row.ValidatorOperAddr, row.ValidatorSelfDelegateAddr)
	}

	return vals, nil
}

// SaveDelegatorsRewardsAmounts implements distribution.DB
func (db *Db) SaveDelegatorsRewardsAmounts(amounts []types.DelegatorRewardAmount) error {
	stmt := `
INSERT INTO delegation_reward_history 
    (validator_address, delegator_address, withdraw_address, amount, height) 
VALUES `
	var params []interface{}

	for i, amount := range amounts {
		ai := i * 5
		stmt += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d),", ai+1, ai+2, ai+3, ai+4, ai+5)

		coins := pq.Array(dbtypes.NewDbDecCoins(amount.Amount))
		params = append(params,
			amount.ValidatorOperAddr, amount.DelegatorAddress, amount.WithdrawAddress, coins, amount.Height)
	}

	stmt = stmt[:len(stmt)-1] // Remove trailing ,
	stmt += `
ON CONFLICT ON CONSTRAINT validator_delegator_unique DO UPDATE 
	SET withdraw_address = excluded.withdraw_address,
		amount = excluded.amount`
	_, err := db.Sql.Exec(stmt, params...)
	return err
}

// SaveValidatorCommissionAmount implements distribution.DB
func (db *Db) SaveValidatorCommissionAmount(amount types.ValidatorCommissionAmount) error {
	stmt := `
INSERT INTO validator_commission_amount_history(self_delegate_address, amount, height) 
VALUES ($1, $2, $3) 
ON CONFLICT ON CONSTRAINT commission_height_unique DO UPDATE 
    SET amount = excluded.amount`

	_, err := db.Sql.Exec(stmt,
		amount.ValidatorSelfDelegateAddr, pq.Array(dbtypes.NewDbDecCoins(amount.Amount)), amount.Height)
	return err
}
