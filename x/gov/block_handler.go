package gov

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/desmos-labs/juno/parse/client"
	"github.com/desmos-labs/juno/parse/worker"
	juno "github.com/desmos-labs/juno/types"
	"github.com/forbole/bdjuno/database"
	ops "github.com/forbole/bdjuno/x/gov/operations"
	"github.com/forbole/bdjuno/x/gov/types"
	"github.com/rs/zerolog/log"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
)

// OneShotOperation represents a method that is called only once when setting up bdjuno
func BlockHandler(block *tmctypes.ResultBlock, _ []juno.Tx, _ *tmctypes.ResultValidators, w worker.Worker) error {
	bigDipperDb, ok := w.Db.(database.BigDipperDb)
	if !ok {
		return fmt.Errorf("provided database is not a BigDipper database")
	}

	// Update existing proposal from LCD
	if err := getHistoricalProposal(block.Block.Height, w.ClientProxy, bigDipperDb); err != nil {
		return err
	}

	return nil
}

// getHistoricalProposal reads from the LCD the current proposals and stores its value inside the database
func getHistoricalProposal(height int64, cp client.ClientProxy, db database.BigDipperDb) error {
	log.Debug().
		Str("module", "gov").
		Str("operation", "existingProposal").
		Msg("getting staking pool")

	var proposals gov.Proposals
	_, err := cp.QueryLCDWithHeight(fmt.Sprintf("/gov/proposals?height=%d", height), &proposals)
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

		if proposal.Status.String() == "VotingPeriod" && proposal.VotingEndTime.After(time.Now()) {
			time.AfterFunc(time.Since(proposal.VotingEndTime), ops.UpdateProposal(proposal.ProposalID, cp, db))
		} else if proposal.Status.String() == "DepositPeriod" && proposal.DepositEndTime.After(time.Now()) {
			time.AfterFunc(time.Since(proposal.DepositEndTime), ops.UpdateProposal(proposal.ProposalID, cp, db))
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
