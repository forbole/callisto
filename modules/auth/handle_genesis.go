package auth

import (
	"encoding/json"
	"fmt"

	"github.com/forbole/bdjuno/database"

	"github.com/cosmos/cosmos-sdk/codec"

	authutils "github.com/forbole/bdjuno/modules/auth/utils"
)

// Handler handles the genesis state of the x/auth module in order to store the initial values
// of the different accounts.
func Handler(appState map[string]json.RawMessage, cdc codec.Marshaler, db *database.Db) error {
	accounts, err := authutils.GetGenesisAccounts(appState, cdc)
	if err != nil {
		return err
	}

	err = db.SaveAccounts(accounts)
	if err != nil {
		return fmt.Errorf("error while storing genesis accounts: %s", err)
	}

	return nil
}
