package operations

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/desmos-labs/juno/parse/client"
	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/x/gov/types"
	stakingops "github.com/forbole/bdjuno/x/staking/operations"
)

// UpdateProposal return a function for time.AfterFunc() to save the to update the proposal status in given time
func UpdateProposal(id uint64, cp client.ClientProxy, db database.BigDipperDb) func() {
	return func() { updateProposalStatuses(id, cp, db) }
}

func updateProposalStatuses(id uint64, cp client.ClientProxy, db database.BigDipperDb) error {
	//update status, voting start time, end time
	var s gov.Proposals
	_, err := cp.QueryLCDWithHeight(fmt.Sprintf("/gov/proposals/%d", id), &s)
	if err != nil {
		return err
	}

	for _, proposal := range s {

		if proposal.Status.String() == "VotingPeriod" {
			update := UpdateProposal(proposal.ProposalID, cp, db)
			time.AfterFunc(time.Since(proposal.VotingEndTime), update)
		}
		//get the voting power in each update
		if err = stakingops.UpdateValidatorVotingPower(cp, db); err != nil {
			return err
		}
		//no metter votingEndTime or votingStarttime it need to update status
		if err = db.UpdateProposal(types.NewProposal(proposal.GetTitle(), proposal.GetDescription(), proposal.ProposalRoute(), proposal.ProposalType(), proposal.ProposalID, proposal.Status,
			proposal.SubmitTime, proposal.DepositEndTime, proposal.VotingStartTime, proposal.VotingEndTime, sdk.AccAddress{})); err != nil {
			return err
		}
	}

	return nil
}
