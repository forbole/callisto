package auth

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/rs/zerolog/log"

	"github.com/forbole/bdjuno/database"
)

// Handler handles the genesis state of the x/auth module in order to store the initial values
// of the different accounts.
func Handler(appState map[string]json.RawMessage, codec *codec.Codec, db *database.BigDipperDb) error {
	log.Debug().Str("module", "auth").Msg("parsing genesis")

	var authState auth.GenesisState
	if err := codec.UnmarshalJSON(appState[auth.ModuleName], &authState); err != nil {
		return err
	}

	// Store the accounts
	for _, account := range authState.Accounts {
		err := db.SaveAccount(account)
		if err != nil {
			return err
		}

		err = db.SaveAccountBalance(account.GetAddress().String(), account.GetCoins(), 1)
		if err != nil {
			return err
		}
	}

	return nil
}
