package gov

import (
	"context"
	"time"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/rs/zerolog/log"

	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/x/gov/types"
)

// UpdateProposal return a function for time.AfterFunc() to update the proposal status on a given time
func UpdateProposal(id uint64, govClient govtypes.QueryClient, db *database.BigDipperDb) func() {
	return func() {
		err := updateProposalStatuses(id, govClient, db)
		if err != nil {
			log.Error().Str("module", "gov").Err(err).Uint64("proposal_id", id).
				Msg("error while updating proposal")
		}
	}
}

func updateProposalStatuses(id uint64, govClient govtypes.QueryClient, db *database.BigDipperDb) error {
	// Get the proposal
	res, err := govClient.Proposal(context.Background(), &govtypes.QueryProposalRequest{ProposalId: id})
	if err != nil {
		return err
	}

	proposal := res.Proposal
	if proposal.Status == govtypes.StatusVotingPeriod {
		update := UpdateProposal(proposal.ProposalId, govClient, db)
		time.AfterFunc(time.Since(proposal.VotingEndTime), update)
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
