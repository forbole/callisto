package utils

import (
	"context"

	juno "github.com/desmos-labs/juno/types"
	authoritytypes "github.com/e-money/em-ledger/x/authority/types"
	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/types"
	"github.com/rs/zerolog/log"
)

// StoreSetGasPricesFromMessage handles storing new gas prices inside the database
func StoreSetGasPricesFromMessage(
	height int64, tx *juno.Tx, msg *authoritytypes.MsgSetGasPrices, authorityClient authoritytypes.QueryClient, db *database.Db,
) error {
	log.Debug().Str("module", "authority/min gas price").Int64("height", height).Msg("updating min gas prices")

	res, err := authorityClient.GasPrices(
		context.Background(),
		&authoritytypes.QueryGasPricesRequest{},
	)

	if err != nil {
		return err
	}

	newEMoneyGasPrices := types.NewEMoneyGasPrices(msg.GetAuthority(), res.GetMinGasPrices(), height)
	return db.SaveEMoneyGasPrices(newEMoneyGasPrices)
}
