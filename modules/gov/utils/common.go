package utils

import (
	"context"
	"fmt"
	"github.com/forbole/bdjuno/modules/utils"

	"github.com/rs/zerolog/log"

	"github.com/forbole/bdjuno/database"
	authutils "github.com/forbole/bdjuno/modules/auth/utils"
	bankutils "github.com/forbole/bdjuno/modules/bank/utils"
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

	err = updateProposalTallyResult(res.Proposal, govClient, db)
	if err != nil {
		return fmt.Errorf("error while updating proposal tally result: %s", err)
	}

	err = updateAccounts(res.Proposal, bankClient, db)
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

// updateProposalTallyResult updates the tally result associated with the given proposal
func updateProposalTallyResult(proposal govtypes.Proposal, govClient govtypes.QueryClient, db *database.Db) error {
	height, err := db.GetLastBlockHeight()
	if err != nil {
		return err
	}

	header := utils.GetHeightRequestHeader(height)
	res, err := govClient.TallyResult(
		context.Background(),
		&govtypes.QueryTallyResultRequest{ProposalId: proposal.ProposalId},
		header,
	)
	if err != nil {
		return err
	}

	return db.SaveTallyResults([]types.TallyResult{
		types.NewTallyResult(
			proposal.ProposalId,
			res.Tally.Yes.Int64(),
			res.Tally.Abstain.Int64(),
			res.Tally.No.Int64(),
			res.Tally.NoWithVeto.Int64(),
			height,
		),
	})
}

// updateAccounts updates any account that might be involved in the proposal (eg. fund community recipient)
func updateAccounts(proposal govtypes.Proposal, bankClient banktypes.QueryClient, db *database.Db) error {
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

		return bankutils.UpdateBalances(addresses, height, bankClient, db)
	}
	return nil
}
