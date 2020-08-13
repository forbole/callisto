package gov

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/desmos-labs/juno/config"
	"github.com/desmos-labs/juno/db"
	"github.com/desmos-labs/juno/parse/client"
	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/x/gov/types"
	"github.com/rs/zerolog/log"
)

// OneShotOperation represents a method that is called only once when setting up bdjuno
func OneShotOperation(_ config.Config, _ *codec.Codec, cp client.ClientProxy, db db.Database) error {
	bigDipperDb, ok := db.(database.BigDipperDb)
	if !ok {
		return fmt.Errorf("provided database is not a BigDipper database")
	}

	// Update the staking pool
	if err := getExistingProposal(cp, bigDipperDb); err != nil {
		return err
	}

	return nil
}

// updateStakingPool reads from the LCD the current staking pool and stores its value inside the database
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
		submitTime := proposal.SubmitTime
		
		depositEndTime:=proposal.DepositEndTime
		votingStartTime:= proposal.VotingStartTime
		votingEndTime:=proposal.VotingEndTime
		genesisTime:=time.Now()
		//since there is not possible to get the proposer, set it to nil
		bdproposals = append(bdproposals, types.NewProposal(proposal.GetTitle(), proposal.GetDescription(), proposal.ProposalRoute(), proposal.ProposalType(), proposal.ProposalID, proposal.Status,
			submitTime, depositEndTime, votingStartTime, votingEndTime, sdk.AccAddress{}))

		bdTallyResult = append(bdTallyResult, types.NewTallyResult(proposal.ProposalID, proposal.FinalTallyResult.Yes.Int64(), proposal.FinalTallyResult.Abstain.Int64(), proposal.FinalTallyResult.No.Int64(),
			proposal.FinalTallyResult.NoWithVeto.Int64(), 0, genesisTime))

		bdDeposit = append(bdDeposit, types.NewDeposit(proposal.ProposalID, sdk.AccAddress{}, proposal.TotalDeposit, proposal.TotalDeposit, 0, genesisTime))

		update := UpdateProposal(proposal.ProposalID, cp, db)
		if proposal.Status.String() == "VotingPeriod" {
			time.AfterFunc(time.Since(votingEndTime), update)
		} else if proposal.Status.String() == "DepositPeriod" {
			time.AfterFunc(time.Since(depositEndTime), update)
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
