package database

import (
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/x/auth/vesting/exported"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	dbtypes "github.com/forbole/bdjuno/v2/database/types"
	dbutils "github.com/forbole/bdjuno/v2/database/utils"
	"github.com/lib/pq"

	"github.com/forbole/bdjuno/v2/types"
)

// SaveAccounts saves the given account addresses inside the database
func (db *Db) SaveAccounts(accounts []types.Account) error {
	paramsNumber := 1
	slices := dbutils.SplitAccounts(accounts, paramsNumber)

	for _, accounts := range slices {
		if len(accounts) == 0 {
			continue
		}

		// Store up-to-date data
		err := db.saveAccounts(paramsNumber, accounts)
		if err != nil {
			return fmt.Errorf("error while storing accounts: %s", err)
		}
	}

	return nil
}

func (db *Db) saveAccounts(paramsNumber int, accounts []types.Account) error {
	if len(accounts) == 0 {
		return nil
	}

	stmt := `INSERT INTO account (address) VALUES `
	var params []interface{}

	for i, account := range accounts {
		ai := i * paramsNumber
		stmt += fmt.Sprintf("($%d),", ai+1)
		params = append(params, account.Address)
	}

	stmt = stmt[:len(stmt)-1]
	stmt += " ON CONFLICT DO NOTHING"
	_, err := db.Sql.Exec(stmt, params...)
	if err != nil {
		return fmt.Errorf("error while storing accounts: %s", err)
	}

	return nil
}

// SaveVestingAccounts saves the given vesting account details inside the database, including:
// type, address, original vesting, endTime, startTime, and vesting periods
func (db *Db) SaveVestingAccounts(vestingAccounts []exported.VestingAccount) error {
	if len(vestingAccounts) == 0 {
		return nil
	}

	for _, account := range vestingAccounts {
		switch vestingAccount := account.(type) {
		case *vestingtypes.ContinuousVestingAccount:
			ContinuousVestingAccount := types.NewContinuousVestingAccount(*vestingAccount)
			err := db.storeContinuousVestingAccount(ContinuousVestingAccount)
			if err != nil {
				return fmt.Errorf("error while storing Continuous Vesting Account: %s", err)
			}

		case *vestingtypes.DelayedVestingAccount:
			DelayedVestingAccount := types.NewDelayedVestingAccount(*vestingAccount)
			err := db.storeDelayedVestingAccount(DelayedVestingAccount)
			if err != nil {
				return fmt.Errorf("error while storing Delayed Vesting Account: %s", err)
			}

		case *vestingtypes.PeriodicVestingAccount:
			PeriodicVestingAccount := types.NewPeriodicVestingAccount(*vestingAccount)
			err := db.storePeriodicVestingAccount(PeriodicVestingAccount)
			if err != nil {
				return fmt.Errorf("error while storing Periodic Vesting Account: %s", err)
			}
		}
	}

	return nil
}

//storeContinuousVestingAccount stores the vesting account details of type ContinuousVestingAccount into the database
func (db *Db) storeContinuousVestingAccount(account types.ContinuousVestingAccount) error {
	stmt := `
INSERT INTO vesting_account (type, address, original_vesting, end_time, start_time) 
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (address) DO UPDATE 
    SET original_vesting = excluded.original_vesting, 
		end_time = excluded.end_time, 
		start_time = excluded.start_time`

	_, err := db.Sql.Exec(stmt,
		account.Type,
		account.Address,
		pq.Array(dbtypes.NewDbCoins(account.OriginalVesting)),
		time.Unix(account.EndTime, 0).Format(time.RFC3339),
		time.Unix(account.StartTime, 0).Format(time.RFC3339),
	)
	if err != nil {
		return err
	}
	return nil
}

//storeDelayedVestingAccount stores the vesting account details of type DelayedVestingAccount into the database
func (db *Db) storeDelayedVestingAccount(account types.DelayedVestingAccount) error {
	stmt := `
INSERT INTO vesting_account (type, address, original_vesting, end_time) 
VALUES ($1, $2, $3, $4)
ON CONFLICT (address) DO UPDATE 
    SET original_vesting = excluded.original_vesting, 
		end_time = excluded.end_time`

	_, err := db.Sql.Exec(
		stmt,
		account.Type,
		account.Address,
		pq.Array(dbtypes.NewDbCoins(account.OriginalVesting)),
		time.Unix(account.EndTime, 0).Format(time.RFC3339),
	)
	if err != nil {
		return err
	}

	return nil
}

//storePeriodicVestingAccount stores the vesting account details of type PeriodicVestingAccount into the database
func (db *Db) storePeriodicVestingAccount(account types.PeriodicVestingAccount) error {
	stmt := `
INSERT INTO vesting_account (type, address, original_vesting, end_time, start_time) 
VALUES ($1, $2, $3, $4, $5) 
ON CONFLICT (address) DO UPDATE 
	SET type = excluded.type,
		original_vesting = excluded.original_vesting,
		end_time = excluded.end_time,
		start_time = excluded.start_time
		RETURNING id `

	var vestingAccountId int
	err := db.Sql.QueryRow(
		stmt,
		account.Type,
		account.Address,
		pq.Array(dbtypes.NewDbCoins(account.OriginalVesting)),
		time.Unix(account.EndTime, 0).Format(time.RFC3339),
		time.Unix(account.StartTime, 0).Format(time.RFC3339),
	).Scan(&vestingAccountId)

	if err != nil {
		return fmt.Errorf("error while saving Periodic Vesting Account: %s", err)
	}

	err = db.storeVestingPeriods(vestingAccountId, account.VestingPeriods)
	if err != nil {
		return fmt.Errorf("error while storing vesting periods: %s", err)
	}

	return nil
}

//storePeriodicVestingAccount stores the vesting periods of type PeriodicVestingAccount into the database
func (db *Db) storeVestingPeriods(vestingAccountId int, vestingPeriods []vestingtypes.Period) error {
	stmt := `
INSERT INTO vesting_period (vesting_account_id, period_order, length, amount) 
VALUES `

	var params []interface{}
	for i, period := range vestingPeriods {
		ai := i * 4
		stmt += fmt.Sprintf("($%d,$%d,$%d,$%d),", ai+1, ai+2, ai+3, ai+4)

		order := i
		amount := pq.Array(dbtypes.NewDbCoins(period.Amount))
		params = append(params, vestingAccountId, order, period.Length, amount)
	}
	stmt = stmt[:len(stmt)-1]

	_, err := db.Sql.Exec(stmt, params...)
	if err != nil {
		return err
	}

	return nil
}

// GetAccounts returns all the accounts that are currently stored inside the database.
func (db *Db) GetAccounts() ([]string, error) {
	var rows []string
	err := db.Sqlx.Select(&rows, `SELECT address FROM account`)
	return rows, err
}
