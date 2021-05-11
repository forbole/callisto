package auth

import (
	"github.com/forbole/bdjuno/modules/common/auth/utils"
)

// UpdateAccounts takes the given addresses and for each one queries the chain
// retrieving the account data and stores it inside the database.
func UpdateAccounts(addresses []string, db DB) error {
	accounts := utils.GetAccounts(addresses)
	return db.SaveAccounts(accounts)
}
