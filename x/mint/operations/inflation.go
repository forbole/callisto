package operations

import (
	"fmt"

	"github.com/rs/zerolog/log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/desmos-labs/juno/parse/client"
	"github.com/forbole/bdjuno/database"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
)

// UpdateInflation fetches from the REST APIs the latest value for the
// inflation, and saves it inside the database.
func UpdateInflation(cp client.ClientProxy, db database.BigDipperDb) error {
	log.Debug().
		Str("module", "mint").
		Str("operation", "inflation").
		Msg("getting inflation data")

	var block tmctypes.ResultBlock
	err := cp.QueryLCD("/blocks/latest", &block)
	if err != nil {
		return err
	}

	var inflation sdk.Dec
	endpoint := fmt.Sprintf("/mint/inflation?height=%d", block.Block.Height)
	height, err := cp.QueryLCDWithHeight(endpoint, &inflation)
	if err != nil {
		return err
	}

	return db.SaveInflation(inflation, height, block.Block.Time)
}
