package common

import (
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/rs/zerolog/log"

	"github.com/forbole/bdjuno/modules/common/bank"

	bigdipperdb "github.com/forbole/bdjuno/database/bigdipper"
)

// RefreshBalance returns a function that when called refreshes the balance of the user having the given address
func RefreshBalance(address string, client banktypes.QueryClient, db *bigdipperdb.Db) func() {
	return func() {
		height, err := db.GetLastBlockHeight()
		if err != nil {
			log.Error().Err(err).Str("module", "bank").
				Str("operation", "refresh balance").Msg("error while getting latest block height")
			return
		}

		err = bank.UpdateBalances([]string{address}, height, client, db)
		if err != nil {
			log.Error().Err(err).Str("module", "bank").
				Str("operation", "refresh balance").Msg("error while updating balance")
		}
	}
}
