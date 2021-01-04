package gov

import (
	"fmt"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/forbole/bdjuno/x/staking"

	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/desmos-labs/juno/client"

	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/x/gov/types"
)

// UpdateProposal return a function for time.AfterFunc() to update the proposal status on a given time
func UpdateProposal(id uint64, cp *client.Proxy, db *database.BigDipperDb) func() {
	return func() {
		err := updateProposalStatuses(id, cp, db)
		if err != nil {
			log.Fatal().Err(err).Send()
		}
	}
}

func updateProposalStatuses(id uint64, cp *client.Proxy, db *database.BigDipperDb) error {
	// Get the proposal
	var proposal gov.Proposal
	_, err := cp.QueryLCDWithHeight(fmt.Sprintf("/gov/proposals/%d", id), &proposal)
	if err != nil {
		return err
	}

	if proposal.Status.String() == "VotingPeriod" {
		update := UpdateProposal(proposal.ProposalID, cp, db)
		time.AfterFunc(time.Since(proposal.VotingEndTime), update)
	}

	// Update the validators voting powers
	err = staking.UpdateValidatorVotingPower(cp, db)
	if err != nil {
		return err
	}

	// Update the proposal to update the status
	return db.UpdateProposal(types.NewProposal(
		proposal.GetTitle(),
		proposal.GetDescription(),
		proposal.ProposalRoute(),
		proposal.ProposalType(),
		proposal.ProposalID,
		proposal.Status,
		proposal.SubmitTime,
		proposal.DepositEndTime,
		proposal.VotingStartTime,
		proposal.VotingEndTime,
		"",
	))
}
