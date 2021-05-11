package common

import (
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/rs/zerolog/log"

	bigdipperdb "github.com/forbole/bdjuno/database/bigdipper"
	"github.com/forbole/bdjuno/modules/common/staking"
)

// UpdateDelegationsAndReplaceExisting updates the delegations of the given delegator by querying them at the
// required height, and then stores them inside the database by replacing all existing ones.
func UpdateDelegationsAndReplaceExisting(
	height int64, delegator string, client stakingtypes.QueryClient, db *bigdipperdb.Db,
) error {
	// Remove existing delegations
	err := db.DeleteDelegatorDelegations(delegator)
	if err != nil {
		return err
	}

	return staking.UpdateDelegations(height, delegator, client, db)
}

// RefreshDelegations returns a function that when called updates the delegations of the provided delegator.
// In order to properly update the data, it removes all the existing delegations and stores new ones querying the gRPC
func RefreshDelegations(height int64, delegator string, client stakingtypes.QueryClient, db *bigdipperdb.Db) func() {
	return func() {
		err := UpdateDelegationsAndReplaceExisting(height, delegator, client, db)
		if err != nil {
			log.Error().Str("module", "staking").Err(err).
				Str("operation", "refresh delegations").Msg("error while refreshing delegations")
		}
	}
}
