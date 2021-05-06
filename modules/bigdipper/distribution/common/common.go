package common

import (
	"context"

	bigdipperdb "github.com/forbole/bdjuno/database/bigdipper"

	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/rs/zerolog/log"
)

// UpdateCommunityPool fetch total amount of coins in the system from RPC and store it into bigdipper
func UpdateCommunityPool(height int64, client distrtypes.QueryClient, db *bigdipperdb.Db) error {
	log.Debug().Str("module", "distribution").Int64("height", height).Msg("getting community pool")

	res, err := client.CommunityPool(context.Background(), &distrtypes.QueryCommunityPoolRequest{})
	if err != nil {
		return err
	}

	// Store the signing infos into the bigdipper
	return db.SaveCommunityPool(res.Pool, height)
}
