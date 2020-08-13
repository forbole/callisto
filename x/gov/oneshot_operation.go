package gov

import (
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/desmos-labs/juno/config"
	"github.com/desmos-labs/juno/db"
	"github.com/desmos-labs/juno/parse/client"
	"github.com/forbole/bdjuno/database"
	ops "github.com/forbole/bdjuno/x/gov/operations"
	"github.com/forbole/bdjuno/x/gov/types"
	"github.com/rs/zerolog/log"
)

// OneShotOperation represents a method that is called only once when setting up bdjuno
func OneShotOperation(_ config.Config, _ *codec.Codec, cp client.ClientProxy, db db.Database) error {
	bigDipperDb, ok := db.(database.BigDipperDb)
	if !ok {
		return fmt.Errorf("provided database is not a BigDipper database")
	}

	// Update existing proposal from LCD
	if err := getExistingProposal(cp, bigDipperDb); err != nil {
		return err
	}

	return nil
}

// getExistingProposal reads from the LCD the current proposals and stores its value inside the database
func getExistingProposal(cp client.ClientProxy, db database.BigDipperDb) error {
	log.Debug().
		Str("module", "gov").
		Str("operation", "existingProposal").
		Msg("getting staking pool")

	var proposals gov.Proposals
	_, err := cp.QueryLCDWithHeight("/gov/proposals", &proposals)
	if err != nil {
		return err
	}

	bdproposals := make([]types.Proposal, len(proposals))
	bdTallyResult := make([]types.TallyResult, len(proposals))
	bdDeposit := make([]types.Deposit, len(proposals))
	for _, proposal := range proposals {
		//since there is not possible to get the proposer, set it to nil
		bdproposals = append(bdproposals, types.NewProposal(proposal.GetTitle(), proposal.GetDescription(), proposal.ProposalRoute(), proposal.ProposalType(), proposal.ProposalID, proposal.Status,
			proposal.SubmitTime, proposal.DepositEndTime, proposal.VotingStartTime, proposal.VotingEndTime, sdk.AccAddress{}))

		bdTallyResult = append(bdTallyResult, types.NewTallyResult(proposal.ProposalID, proposal.FinalTallyResult.Yes.Int64(), proposal.FinalTallyResult.Abstain.Int64(), proposal.FinalTallyResult.No.Int64(),
			proposal.FinalTallyResult.NoWithVeto.Int64(), 0, time.Now()))

		bdDeposit = append(bdDeposit, types.NewDeposit(proposal.ProposalID, sdk.AccAddress{}, proposal.TotalDeposit, proposal.TotalDeposit, 0, time.Now()))

		update := ops.UpdateProposal(proposal.ProposalID, cp, db)
		if proposal.Status.String() == "VotingPeriod" {
			time.AfterFunc(time.Since(proposal.VotingEndTime), update)
		} else if proposal.Status.String() == "DepositPeriod" {
			time.AfterFunc(time.Since(proposal.DepositEndTime), update)
		}
	}
	if err := db.SaveProposals(bdproposals); err != nil {
		return nil
	}

	if err := db.SaveDeposits(bdDeposit); err != nil {
		return nil
	}

	return db.SaveTallyResults(bdTallyResult)
}
