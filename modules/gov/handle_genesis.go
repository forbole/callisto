package gov

import (
	"encoding/json"
	"fmt"

	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/types"

	"github.com/cosmos/cosmos-sdk/codec"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/rs/zerolog/log"
)

func HandleGenesis(appState map[string]json.RawMessage, cdc codec.Marshaler, db *database.Db) error {
	log.Debug().Str("module", "gov").Msg("parsing genesis")

	// Read the genesis state
	var genState govtypes.GenesisState
	err := cdc.UnmarshalJSON(appState[govtypes.ModuleName], &genState)
	if err != nil {
		return fmt.Errorf("error while reading gov genesis data: %s", err)
	}

	// Save the proposals
	err = saveProposals(genState.Proposals, db)
	if err != nil {
		return fmt.Errorf("error while storing genesis governance proposals: %s", err)
	}

	return nil
}

// saveProposals save proposals from genesis file
func saveProposals(slice govtypes.Proposals, db *database.Db) error {
	proposals := make([]types.Proposal, len(slice))
	tallyResults := make([]types.TallyResult, len(slice))
	deposits := make([]types.Deposit, len(slice))

	for index, proposal := range slice {
		// Since it's not possible to get the proposer, set it to nil
		proposals[index] = types.NewProposal(
			proposal.ProposalId,
			proposal.ProposalRoute(),
			proposal.ProposalType(),
			proposal.GetContent(),
			proposal.Status.String(),
			proposal.SubmitTime,
			proposal.DepositEndTime,
			proposal.VotingStartTime,
			proposal.VotingEndTime,
			"",
		)

		tallyResults[index] = types.NewTallyResult(
			proposal.ProposalId,
			proposal.FinalTallyResult.Yes.Int64(),
			proposal.FinalTallyResult.Abstain.Int64(),
			proposal.FinalTallyResult.No.Int64(),
			proposal.FinalTallyResult.NoWithVeto.Int64(),
			1,
		)

		deposits[index] = types.NewDeposit(
			proposal.ProposalId,
			"",
			proposal.TotalDeposit,
			1,
		)
	}

	// Save the proposals
	err := db.SaveProposals(proposals)
	if err != nil {
		return err
	}

	// Save the deposits
	err = db.SaveDeposits(deposits)
	if err != nil {
		return err
	}

	// Save the tally results
	return db.SaveTallyResults(tallyResults)
}
