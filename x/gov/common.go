package gov

import (
	"context"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"

	bauthutils "github.com/forbole/bdjuno/x/auth/common"
	bbankutils "github.com/forbole/bdjuno/x/bank/common"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/rs/zerolog/log"

	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/x/gov/types"
)

// UpdateProposal return a function for time.AfterFunc() to update the proposal status on a given time
func UpdateProposal(
	id uint64,
	govClient govtypes.QueryClient, authClient authtypes.QueryClient, bankClient banktypes.QueryClient,
	cdc codec.Marshaler, db *database.BigDipperDb,
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

		err = updateAccount(res.Proposal, authClient, bankClient, cdc, db)
		if err != nil {
			log.Error().Str("module", "gov").Err(err).Uint64("proposal_id", id).
				Msg("error while updating proposal related accounts balances")
		}
	}
}

func updateProposalStatuses(
	proposal govtypes.Proposal,
	govClient govtypes.QueryClient, authClient authtypes.QueryClient, bankClient banktypes.QueryClient,
	cdc codec.Marshaler, db *database.BigDipperDb,
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

func updateAccount(
	proposal govtypes.Proposal, authClient authtypes.QueryClient, bankClient banktypes.QueryClient,
	cdc codec.Marshaler, db *database.BigDipperDb,
) error {
	content, ok := proposal.Content.GetCachedValue().(*distrtypes.CommunityPoolSpendProposal)
	if ok {
		height, err := db.GetLastBlockHeight()
		if err != nil {
			return err
		}

		addresses := []string{content.Recipient}

		err = bauthutils.UpdateAccounts(addresses, height, authClient, cdc, db)
		if err != nil {
			return err
		}

		return bbankutils.UpdateBalances(addresses, height, bankClient, db)
	}
	return nil
}
