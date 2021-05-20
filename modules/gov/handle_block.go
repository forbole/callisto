package gov

import (
	"context"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/rs/zerolog/log"

	"github.com/forbole/bdjuno/modules/utils"
	"github.com/forbole/bdjuno/types"

	"github.com/forbole/bdjuno/database"
	govutils "github.com/forbole/bdjuno/modules/gov/utils"
)

// HandleBlock handles a new block by updating any eventually open proposal's status and tally result
func HandleBlock(height int64, govClient govtypes.QueryClient, bankClient banktypes.QueryClient, db *database.Db) error {
	err := updateProposals(govClient, bankClient, db)
	if err != nil {
		log.Error().Str("module", "gov").Int64("height", height).
			Err(err).Msg("error while updating proposals")
	}

	err = updateParams(height, govClient, db)
	if err != nil {
		log.Error().Str("module", "gov").Int64("height", height).
			Err(err).Msg("error while updating params")
	}

	return nil
}

// updateParams updates the governance parameters for the given height
func updateParams(height int64, govClient govtypes.QueryClient, db *database.Db) error {
	depositRes, err := govClient.Params(
		context.Background(),
		&govtypes.QueryParamsRequest{ParamsType: govtypes.ParamDeposit},
		utils.GetHeightRequestHeader(height),
	)
	if err != nil {
		return err
	}

	votingRes, err := govClient.Params(
		context.Background(),
		&govtypes.QueryParamsRequest{ParamsType: govtypes.ParamVoting},
		utils.GetHeightRequestHeader(height),
	)
	if err != nil {
		return err
	}

	tallyRes, err := govClient.Params(
		context.Background(),
		&govtypes.QueryParamsRequest{ParamsType: govtypes.ParamTallying},
		utils.GetHeightRequestHeader(height),
	)
	if err != nil {
		return err
	}

	return db.SaveGovParams(types.NewGovParams(
		govtypes.NewParams(
			votingRes.GetVotingParams(),
			tallyRes.GetTallyParams(),
			depositRes.GetDepositParams(),
		),
		height,
	))
}

// updateProposals updates the proposals
func updateProposals(govClient govtypes.QueryClient, bankClient banktypes.QueryClient, db *database.Db) error {
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
