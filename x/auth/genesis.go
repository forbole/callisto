package auth

import (
	"encoding/json"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/rs/zerolog/log"

	"github.com/forbole/bdjuno/database"
)

// Handler handles the genesis state of the x/auth module in order to store the initial values
// of the different accounts.
func Handler(appState map[string]json.RawMessage, cdc codec.Marshaler, db *database.BigDipperDb) error {
	log.Debug().Str("module", "auth").Msg("parsing genesis")

	var authState authtypes.GenesisState
	if err := cdc.UnmarshalJSON(appState[authtypes.ModuleName], &authState); err != nil {
		return err
	}

	// Store the accounts
	for _, account := range authState.Accounts {
		var accountI authtypes.AccountI
		err := cdc.UnpackAny(account, &accountI)
		if err != nil {
			return err
		}

		err = db.SaveAccount(accountI)
		if err != nil {
			return err
		}
	}

	return nil
}
