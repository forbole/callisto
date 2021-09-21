package gov

import (
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/desmos-labs/juno/client"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/rs/zerolog/log"

	"github.com/forbole/bdjuno/types"

	"github.com/forbole/bdjuno/database"
	govutils "github.com/forbole/bdjuno/modules/gov/utils"
)

// HandleBlock handles a new block by updating any eventually open proposal's status and tally result
func HandleBlock(
	height int64, blockVals *tmctypes.ResultValidators,
	govClient govtypes.QueryClient, bankClient banktypes.QueryClient, stakingClient stakingtypes.QueryClient,
	cdc codec.Marshaler, db *database.Db,
) error {
	err := updateProposals(height, blockVals, govClient, bankClient, stakingClient, cdc, db)
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
	header := client.GetHeightRequestHeader(height)
	depositRes, err := govClient.Params(
		context.Background(),
		&govtypes.QueryParamsRequest{ParamsType: govtypes.ParamDeposit},
		header,
	)
	if err != nil {
		return fmt.Errorf("error while getting gov deposit params: %s", err)
	}

	votingRes, err := govClient.Params(
		context.Background(),
		&govtypes.QueryParamsRequest{ParamsType: govtypes.ParamVoting},
		header,
	)
	if err != nil {
		return fmt.Errorf("error while getting gov voting params: %s", err)
	}

	tallyRes, err := govClient.Params(
		context.Background(),
		&govtypes.QueryParamsRequest{ParamsType: govtypes.ParamTallying},
		header,
	)
	if err != nil {
		return fmt.Errorf("error while getting gov tally params: %s", err)
	}

	return db.SaveGovParams(types.NewGovParams(
		types.NewVotingParams(votingRes.GetVotingParams()),
		types.NewDepositParam(depositRes.GetDepositParams()),
		types.NewTallyParams(tallyRes.GetTallyParams()),
		height,
	))
}

// updateProposals updates the proposals
func updateProposals(
	height int64, blockVals *tmctypes.ResultValidators,
	govClient govtypes.QueryClient, bankClient banktypes.QueryClient, stakingClient stakingtypes.QueryClient,
	cdc codec.Marshaler, db *database.Db,
) error {
	ids, err := db.GetOpenProposalsIds()
	if err != nil {
		log.Error().Err(err).Str("module", "gov").Msg("error while getting open ids")
	}

	for _, id := range ids {
		err = govutils.UpdateProposal(height, blockVals, id, govClient, bankClient, stakingClient, cdc, db)
		if err != nil {
			return fmt.Errorf("error while updating proposal: %s", err)
		}
	}
	return nil
}
