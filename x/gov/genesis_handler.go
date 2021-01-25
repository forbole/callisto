package gov

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/rs/zerolog/log"

	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/x/gov/types"
)

func HandleGenesis(
	appState map[string]json.RawMessage, cdc codec.Marshaler, govClient govtypes.QueryClient, db *database.BigDipperDb,
) error {
	log.Debug().Str("module", "gov").Msg("parsing genesis")

	// Read the genesis state
	var genState govtypes.GenesisState
	err := cdc.UnmarshalJSON(appState[govtypes.ModuleName], &genState)
	if err != nil {
		return fmt.Errorf("error while reading gov genesis data: %s", err)
	}

	// Save the proposals
	err = saveProposals(genState.Proposals, govClient, db)
	if err != nil {
		return fmt.Errorf("error while storing genesis governance proposals: %s", err)
	}

	return nil
}

// saveProposals save proposals from genesis file
func saveProposals(p govtypes.Proposals, govClient govtypes.QueryClient, db *database.BigDipperDb) error {
	proposals := make([]types.Proposal, len(p))
	tallyResults := make([]types.TallyResult, len(p))
	deposits := make([]types.Deposit, len(p))

	for index, proposal := range p {
		// Since it's not possible to get the proposer, set it to nil
		proposals[index] = types.NewProposal(
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

		// Update the proposal status when the voting period or deposit period ends
		update := UpdateProposal(proposal.ProposalId, govClient, db)
		if proposal.Status == govtypes.StatusVotingPeriod {
			time.AfterFunc(time.Since(proposal.VotingEndTime), update)
		} else if proposal.Status == govtypes.StatusDepositPeriod {
			time.AfterFunc(time.Since(proposal.DepositEndTime), update)
		}
	}

	// Save the proposals
	err := db.SaveProposals(proposals)
	if err != nil {
		return nil
	}

	// Save the deposits
	err = db.SaveDeposits(deposits)
	if err != nil {
		return nil
	}

	// Save the tally results
	return db.SaveTallyResults(tallyResults)
}
