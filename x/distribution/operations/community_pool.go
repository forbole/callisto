package operations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/desmos-labs/juno/parse/client"
	"github.com/forbole/bdjuno/database"
	"github.com/rs/zerolog/log"
)

// UpdateCommunityPool fetch total amount of coins in the system from RPC and store it into database
func UpdateCommunityPool(cp client.ClientProxy, db database.BigDipperDb) error {
	log.Debug().
		Str("module", "distribution").
		Str("operation", "community pool").
		Msg("getting community pool")
	var s sdk.DecCoins
	height, err := cp.QueryLCDWithHeight("/distribution/community_pool", &s)
	if err != nil {
		return err
	}
	// Store the signing infos into the database
	err = db.SaveCommunityPool(s, height)
	return err
}
