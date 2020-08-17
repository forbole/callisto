package gov

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/desmos-labs/juno/parse/client"
	"github.com/desmos-labs/juno/parse/worker"
	"github.com/forbole/bdjuno/database"
	ops "github.com/forbole/bdjuno/x/gov/operations"
	"github.com/forbole/bdjuno/x/gov/types"
	"github.com/rs/zerolog/log"
	tmtypes "github.com/tendermint/tendermint/types"
)

func GenesisHandler(codec *codec.Codec, genesisDoc *tmtypes.GenesisDoc, appState map[string]json.RawMessage, w worker.Worker) error {
	log.Debug().Str("module", "gov").Msg("parsing genesis")

	bigDipperDb, ok := w.Db.(database.BigDipperDb)
	if !ok {
		return fmt.Errorf("given database instance is not a BigDipperDb")
	}
	// Read the genesis state
	var genState gov.GenesisState
	if err := codec.UnmarshalJSON(appState[gov.ModuleName], &genState); err != nil {
		return err
	}

	if err := saveProposals(genState.Proposals, genesisDoc, bigDipperDb, w.ClientProxy); err != nil {
		return err
	}
	return nil
}

// saveProposals save proposals from genesis file
func saveProposals(proposals gov.Proposals, genesisDoc *tmtypes.GenesisDoc, db database.BigDipperDb, cp client.ClientProxy) error {
	bdproposals := make([]types.Proposal, len(proposals))
	bdTallyResult := make([]types.TallyResult, len(proposals))
	bdDeposit := make([]types.Deposit, len(proposals))
	for _, proposal := range proposals {
		submitTime, err := time.Parse(time.RFC3339, proposal.SubmitTime.String())
		if err != nil {
			return err
		}
		depositEndTime, err := time.Parse(time.RFC3339, proposal.DepositEndTime.String())
		if err != nil {
			return err
		}
		votingStartTime, err := time.Parse(time.RFC3339, proposal.VotingStartTime.String())
		if err != nil {
			return err
		}
		votingEndTime, err := time.Parse(time.RFC3339, proposal.VotingEndTime.String())
		if err != nil {
			return err
		}
		genesisTime, err := time.Parse(time.RFC3339, genesisDoc.GenesisTime.String())
		if err != nil {
			return err
		}

		//since there is not possible to get the proposer, set it to nil
		bdproposals = append(bdproposals, types.NewProposal(proposal.GetTitle(), proposal.GetDescription(), proposal.ProposalRoute(), proposal.ProposalType(), proposal.ProposalID, proposal.Status,
			submitTime, depositEndTime, votingStartTime, votingEndTime, sdk.AccAddress{}))

		bdTallyResult = append(bdTallyResult, types.NewTallyResult(proposal.ProposalID, proposal.FinalTallyResult.Yes.Int64(), proposal.FinalTallyResult.Abstain.Int64(), proposal.FinalTallyResult.No.Int64(),
			proposal.FinalTallyResult.NoWithVeto.Int64(), 0, genesisTime))

		bdDeposit = append(bdDeposit, types.NewDeposit(proposal.ProposalID, sdk.AccAddress{}, proposal.TotalDeposit, proposal.TotalDeposit, 0, genesisTime))

		update := ops.UpdateProposal(proposal.ProposalID, cp, db)
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
