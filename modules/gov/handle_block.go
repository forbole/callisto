package gov

import (
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/rs/zerolog/log"

	"github.com/forbole/bdjuno/database"
	govutils "github.com/forbole/bdjuno/modules/gov/utils"
)

// HandleBlock handles a new block by updating any eventually open proposal's status and tally result
func HandleBlock(govClient govtypes.QueryClient, bankClient banktypes.QueryClient, db *database.Db) error {
	ids, err := db.GetOpenProposalsIds()
	if err != nil {
		log.Error().Err(err).Str("module", "gov").Msg("error while getting open ids")
	}

	for _, id := range ids {
		err = govutils.UpdateProposal(id, govClient, bankClient, db)
		if err != nil {
			return err
		}
	}

	return nil
}
