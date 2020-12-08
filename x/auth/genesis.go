package auth

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/forbole/bdjuno/database"
	"github.com/rs/zerolog/log"
	"github.com/tendermint/tendermint/types"
)

// Handler handles the genesis state of the x/auth module in order to store the initial values
// of the different accounts.
func Handler(
	genDoc *types.GenesisDoc, appState map[string]json.RawMessage, codec *codec.Codec, db *database.BigDipperDb,
) error {
	log.Debug().Str("module", "auth").Msg("parsing genesis")

	var authState auth.GenesisState
	if err := codec.UnmarshalJSON(appState[auth.ModuleName], &authState); err != nil {
		return err
	}

	// Store the accounts
	accounts := make([]exported.Account, len(authState.Accounts))
	for index, account := range authState.Accounts {
		accounts[index] = account.(exported.Account)
	}
	if err := db.SaveAccounts(accounts, 0, genDoc.GenesisTime); err != nil {
		return err
	}

	return nil
}
