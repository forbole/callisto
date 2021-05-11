package auth

import (
	"encoding/json"
	"fmt"

	"github.com/forbole/bdjuno/modules/common/auth/utils"

	"github.com/cosmos/cosmos-sdk/codec"
)

// Handler handles the genesis state of the x/auth module in order to store the initial values
// of the different accounts.
func Handler(appState map[string]json.RawMessage, cdc codec.Marshaler, db DB) error {
	accounts, err := utils.GetGenesisAccounts(appState, cdc)
	if err != nil {
		return err
	}

	err = db.SaveAccounts(accounts)
	if err != nil {
		return fmt.Errorf("error while storing genesis accounts: %s", err)
	}

	return nil
}
