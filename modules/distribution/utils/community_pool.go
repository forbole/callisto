package utils

import (
	"context"
	"fmt"

	"github.com/forbole/bdjuno/database"

	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/rs/zerolog/log"
)

// UpdateCommunityPool fetch total amount of coins in the system from RPC and store it into database
func UpdateCommunityPool(height int64, client distrtypes.QueryClient, db *database.Db) error {
	log.Debug().Str("module", "distribution").Int64("height", height).Msg("getting community pool")

	res, err := client.CommunityPool(context.Background(), &distrtypes.QueryCommunityPoolRequest{})
	if err != nil {
		return fmt.Errorf("error while getting comminity pool: %s", err)
	}

	// Store the signing infos into the database
	return db.SaveCommunityPool(res.Pool, height)
}
