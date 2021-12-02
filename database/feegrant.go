package database

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	feegranttypes "github.com/cosmos/cosmos-sdk/x/feegrant"
	"github.com/forbole/bdjuno/v2/types"
)

// SaveFeeGrantAllowance allows to store the fee grant allowances for the given block height
func (db *Db) SaveFeeGrantAllowance(allowance *feegranttypes.MsgGrantAllowance, height int64) error {
	stmt := `
INSERT INTO fee_grant_allowance(grantee, granter, allowance, height) 
VALUES ($1, $2, $3, $4) 
ON CONFLICT DO NOTHING`

	allowanceJSON, err := codec.ProtoMarshalJSON(allowance.Allowance, nil)
	if err != nil {
		return fmt.Errorf("error while marshaling grant allowance: %s", err)
	}

	// Store the accounts
	var accounts []types.Account
	accounts = append(accounts, types.NewAccount(allowance.Grantee), types.NewAccount(allowance.Granter))
	err = db.SaveAccounts(accounts)
	if err != nil {
		return fmt.Errorf("error while storing fee grant allowance accounts: %s", err)
	}

	_, err = db.Sql.Exec(stmt, allowance.Grantee, allowance.Granter, allowanceJSON, height)
	if err != nil {
		return fmt.Errorf("error while storing fee grant allowance: %s", err)
	}

	return nil
}

// RevokeFeeGrantAllowance removes the fee grant allowances data from the database
func (db *Db) RevokeFeeGrantAllowance(grantee string, granter string) error {
	_, err := db.Sql.Exec(`DELETE FROM fee_grant_allowance WHERE grantee = $1 AND granter = $2`, grantee, granter)
	if err != nil {
		return fmt.Errorf("error while revoking grant allowance: %s", err)
	}
	return nil
}
