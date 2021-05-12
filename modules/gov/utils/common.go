package utils

import (
	"context"

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

	"github.com/rs/zerolog/log"
)

// UpdateProposal return a function for time.AfterFunc() to update the proposal status on a given time
func UpdateProposal(
	id uint64,
	govClient govtypes.QueryClient, authClient authtypes.QueryClient, bankClient banktypes.QueryClient,
	cdc codec.Marshaler, db *database.Db,
) func() {
	return func() {
		// Get the proposal
		res, err := govClient.Proposal(context.Background(), &govtypes.QueryProposalRequest{ProposalId: id})
		if err != nil {
			log.Error().Str("module", "gov").Err(err).Uint64("proposal_id", id).
				Msg("error while getting proposal")
			return
		}

		err = updateProposalStatuses(res.Proposal, govClient, authClient, bankClient, cdc, db)
		if err != nil {
			log.Error().Str("module", "gov").Err(err).Uint64("proposal_id", id).
				Msg("error while updating proposal")
			return
		}

		err = updateAccount(res.Proposal, bankClient, db)
		if err != nil {
			log.Error().Str("module", "gov").Err(err).Uint64("proposal_id", id).
				Msg("error while updating proposal related accounts balances")
		}
	}
}

func updateProposalStatuses(
	proposal govtypes.Proposal,
	govClient govtypes.QueryClient, authClient authtypes.QueryClient, bankClient banktypes.QueryClient,
	cdc codec.Marshaler, db *database.Db,
) error {
	if proposal.Status == govtypes.StatusVotingPeriod {
		update := UpdateProposal(proposal.ProposalId, govClient, authClient, bankClient, cdc, db)
		time.AfterFunc(time.Until(proposal.VotingEndTime), update)
	}

	// Update the proposal to update the status
	return db.UpdateProposal(types.NewProposal(
		proposal.GetTitle(),
		proposal.GetContent().GetDescription(),
		proposal.ProposalRoute(),
		proposal.ProposalType(),
		proposal.ProposalId,
		proposal.Status,
		proposal.SubmitTime,
		proposal.DepositEndTime,
		proposal.VotingStartTime,
		proposal.VotingEndTime,
		"",
	))
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
