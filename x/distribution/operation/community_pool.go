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
		Str("module", "staking").
		Str("operation", " tokens").
		Msg("getting total token supply")
	var s sdk.Coins
	height, err := cp.QueryLCDWithHeight("/distribution/community_pool", &s)
	if err != nil {
		return err
	}
	// Store the signing infos into the database
	err = db.SaveCommunityPool(s,
		height,
	)
	if err != nil {
		return err
	}
	return nil
}
