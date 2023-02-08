package database

import (
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/x/auth/vesting/exported"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	"github.com/gogo/protobuf/proto"
	"github.com/lib/pq"

	dbtypes "github.com/forbole/bdjuno/v3/database/types"
	dbutils "github.com/forbole/bdjuno/v3/database/utils"

	"github.com/forbole/bdjuno/v3/types"
)

// SaveAccounts saves the given accounts inside the database
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
	_, err := db.SQL.Exec(stmt, params...)
	if err != nil {
		return fmt.Errorf("error while storing accounts: %s", err)
	}

	return nil
}

// SaveVestingAccounts saves the given vesting accounts inside the database
func (db *Db) SaveVestingAccounts(vestingAccounts []exported.VestingAccount) error {
	if len(vestingAccounts) == 0 {
		return nil
	}

	for _, account := range vestingAccounts {
		switch vestingAccount := account.(type) {
		case *vestingtypes.ContinuousVestingAccount, *vestingtypes.DelayedVestingAccount:
			_, err := db.storeVestingAccount(account)
			if err != nil {
				return err
			}

		case *vestingtypes.PeriodicVestingAccount:
			vestingAccountRowID, err := db.storeVestingAccount(account)
			if err != nil {
				return err
			}
			err = db.storeVestingPeriods(vestingAccountRowID, vestingAccount.VestingPeriods)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (db *Db) storeVestingAccount(account exported.VestingAccount) (int, error) {
	stmt := `
	INSERT INTO vesting_account (type, address, original_vesting, end_time, start_time) 
	VALUES ($1, $2, $3, $4, $5)
	ON CONFLICT (address) DO UPDATE 
		SET original_vesting = excluded.original_vesting, 
			end_time = excluded.end_time, 
			start_time = excluded.start_time
			RETURNING id `

	var vestingAccountRowID int
	err := db.SQL.QueryRow(stmt,
		proto.MessageName(account),
		account.GetAddress().String(),
		pq.Array(dbtypes.NewDbCoins(account.GetOriginalVesting())),
		time.Unix(account.GetEndTime(), 0),
		time.Unix(account.GetStartTime(), 0),
	).Scan(&vestingAccountRowID)

	if err != nil {
		return vestingAccountRowID, fmt.Errorf("error while saving Vesting Account of type %v: %s", proto.MessageName(account), err)
	}

	return vestingAccountRowID, nil
}

func (db *Db) StoreBaseVestingAccountFromMsg(bva *vestingtypes.BaseVestingAccount, txTimestamp time.Time) error {
	stmt := `
	INSERT INTO vesting_account (type, address, original_vesting, start_time, end_time) 
	VALUES ($1, $2, $3, $4, $5)
	ON CONFLICT (address) DO UPDATE 
		SET type = excluded.type,
			original_vesting = excluded.original_vesting, 
			start_time = excluded.start_time, 
			end_time = excluded.end_time`

	_, err := db.SQL.Exec(stmt,
		proto.MessageName(bva),
		bva.GetAddress().String(),
		pq.Array(dbtypes.NewDbCoins(bva.OriginalVesting)),
		txTimestamp,
		time.Unix(bva.EndTime, 0))
	if err != nil {
		return fmt.Errorf("error while storing vesting account: %s", err)
	}
	return nil
}

// storeVestingPeriods handles storing the vesting periods of PeriodicVestingAccount type
func (db *Db) storeVestingPeriods(id int, vestingPeriods []vestingtypes.Period) error {
	// Delete already existing periods
	stmt := `DELETE FROM vesting_period WHERE vesting_account_id = $1`
	_, err := db.SQL.Exec(stmt, id)
	if err != nil {
		return fmt.Errorf("error while deleting vesting period: %s", err)
	}

	// Store the new periods
	stmt = `
INSERT INTO vesting_period (vesting_account_id, period_order, length, amount) 
VALUES `

	var params []interface{}
	for i, period := range vestingPeriods {
		ai := i * 4
		stmt += fmt.Sprintf("($%d,$%d,$%d,$%d),", ai+1, ai+2, ai+3, ai+4)

		order := i
		amount := pq.Array(dbtypes.NewDbCoins(period.Amount))
		params = append(params, id, order, period.Length, amount)
	}
	stmt = stmt[:len(stmt)-1]

	_, err = db.SQL.Exec(stmt, params...)
	if err != nil {
		return fmt.Errorf("error while saving vesting periods: %s", err)
	}

	return nil
}

// GetAccounts returns all the accounts that are currently stored inside the database.
func (db *Db) GetAccounts() ([]string, error) {
	var rows []string
	err := db.Sqlx.Select(&rows, `SELECT address FROM account`)
	return rows, err
}
