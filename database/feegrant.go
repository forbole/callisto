package database

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/forbole/bdjuno/v2/types"
)

// SaveFeeGrantAllowance allows to store the fee grant allowances for the given block height
func (db *Db) SaveFeeGrantAllowance(allowance types.FeeGrant) error {

	// Store the accounts
	var accounts []types.Account
	accounts = append(accounts, types.NewAccount(allowance.Grantee), types.NewAccount(allowance.Granter))
	err := db.SaveAccounts(accounts)
	if err != nil {
		return fmt.Errorf("error while storing fee grant allowance accounts: %s", err)
	}

	stmt := `
INSERT INTO fee_grant_allowance(grantee, granter, allowance, height) 
VALUES ($1, $2, $3, $4) 
ON CONFLICT ON CONSTRAINT unique_fee_grant_allowance DO UPDATE 
    SET allowance = excluded.allowance,
        height = excluded.height
WHERE fee_grant_allowance.height <= excluded.height`

	allowanceJSON, err := codec.ProtoMarshalJSON(allowance.Allowance, nil)
	if err != nil {
		return fmt.Errorf("error while marshaling grant allowance: %s", err)
	}

	_, err = db.Sql.Exec(stmt, allowance.Grantee, allowance.Granter, allowanceJSON, allowance.Height)
	if err != nil {
		return fmt.Errorf("error while saving fee grant allowance: %s", err)
	}

	return nil
}

// DeleteFeeGrantAllowance removes the fee grant allowance data from the database
func (db *Db) DeleteFeeGrantAllowance(allowance types.GrantRemoval) error {
	stmt := `DELETE FROM fee_grant_allowance WHERE grantee = $1 AND granter = $2 AND height <= $3`
	_, err := db.Sql.Exec(stmt, allowance.Grantee, allowance.Granter, allowance.Height)

	if err != nil {
		return fmt.Errorf("error while deleting grant allowance: %s", err)
	}
	return nil
}
