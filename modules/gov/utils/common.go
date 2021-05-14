package utils

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/forbole/bdjuno/database"
	authutils "github.com/forbole/bdjuno/modules/auth/utils"
	"github.com/forbole/bdjuno/modules/bank/utils"
	"github.com/forbole/bdjuno/types"

	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

// RefreshProposal return a function for time.AfterFunc() to update the proposal status on a given time
func RefreshProposal(
	id uint64,
	govClient govtypes.QueryClient, authClient authtypes.QueryClient, bankClient banktypes.QueryClient,
	cdc codec.Marshaler, db *database.Db,
) func() {
	return func() {
		err := UpdateProposal(id, govClient, authClient, bankClient, cdc, db)
		if err != nil {
			log.Error().Str("module", "gov").Err(err).Uint64("proposal_id", id).
				Msg("error while refreshing proposal")
		}
	}
}

func UpdateProposal(
	id uint64,
	govClient govtypes.QueryClient, authClient authtypes.QueryClient, bankClient banktypes.QueryClient,
	cdc codec.Marshaler, db *database.Db,
) error {
	// Get the proposal
	res, err := govClient.Proposal(context.Background(), &govtypes.QueryProposalRequest{ProposalId: id})
	if err != nil {
		return fmt.Errorf("error while getting proposal: %s", err)
	}

	err = updateProposalStatus(res.Proposal, govClient, authClient, bankClient, cdc, db)
	if err != nil {
		return fmt.Errorf("error while updating proposal status: %s", err)
	}

	err = updateAccount(res.Proposal, bankClient, db)
	if err != nil {
		return fmt.Errorf("error while updating account: %s", err)
	}

	return nil
}

func updateProposalStatus(
	proposal govtypes.Proposal,
	govClient govtypes.QueryClient, authClient authtypes.QueryClient, bankClient banktypes.QueryClient,
	cdc codec.Marshaler, db *database.Db,
) error {
	if proposal.Status == govtypes.StatusVotingPeriod {
		update := RefreshProposal(proposal.ProposalId, govClient, authClient, bankClient, cdc, db)
		time.AfterFunc(time.Until(proposal.VotingEndTime), update)
	}

	// Update the proposal to update the status
	return db.UpdateProposal(
		types.NewProposalUpdate(
			proposal.ProposalId,
			proposal.Status,
			proposal.VotingStartTime,
			proposal.VotingEndTime,
		),
	)
}

func updateAccount(proposal govtypes.Proposal, bankClient banktypes.QueryClient, db *database.Db) error {
	content, ok := proposal.Content.GetCachedValue().(*distrtypes.CommunityPoolSpendProposal)
	if ok {
		height, err := db.GetLastBlockHeight()
		if err != nil {
			return err
		}

		addresses := []string{content.Recipient}

		err = authutils.UpdateAccounts(addresses, db)
		if err != nil {
			return err
		}

		return utils.UpdateBalances(addresses, height, bankClient, db)
	}
	return nil
}
