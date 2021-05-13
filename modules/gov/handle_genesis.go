package gov

import (
	"encoding/json"
	"fmt"
	"time"

	govutils "github.com/forbole/bdjuno/modules/gov/utils"

	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/types"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/cosmos/cosmos-sdk/codec"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/rs/zerolog/log"
)

func HandleGenesis(
	appState map[string]json.RawMessage,
	govClient govtypes.QueryClient, authClient authtypes.QueryClient, bankClient banktypes.QueryClient,
	cdc codec.Marshaler, db *database.Db,
) error {
	log.Debug().Str("module", "gov").Msg("parsing genesis")

	// Read the genesis state
	var genState govtypes.GenesisState
	err := cdc.UnmarshalJSON(appState[govtypes.ModuleName], &genState)
	if err != nil {
		return fmt.Errorf("error while reading gov genesis data: %s", err)
	}

	// Save the proposals
	err = saveProposals(genState.Proposals, govClient, authClient, bankClient, cdc, db)
	if err != nil {
		return fmt.Errorf("error while storing genesis governance proposals: %s", err)
	}

	return nil
}

// saveProposals save proposals from genesis file
func saveProposals(
	p govtypes.Proposals,
	govClient govtypes.QueryClient, authClient authtypes.QueryClient, bankClient banktypes.QueryClient,
	cdc codec.Marshaler, db *database.Db,
) error {
	proposals := make([]types.Proposal, len(p))
	tallyResults := make([]types.TallyResult, len(p))
	deposits := make([]types.Deposit, len(p))

	for index, proposal := range p {
		// Since it's not possible to get the proposer, set it to nil
		proposals[index] = types.NewProposal(
			proposal.ProposalId,
			proposal.ProposalRoute(),
			proposal.ProposalType(),
			proposal.GetContent(),
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
		update := govutils.UpdateProposal(proposal.ProposalId, govClient, authClient, bankClient, cdc, db)
		if proposal.Status == govtypes.StatusVotingPeriod {
			time.AfterFunc(time.Until(proposal.VotingEndTime), update)
		} else if proposal.Status == govtypes.StatusDepositPeriod {
			time.AfterFunc(time.Until(proposal.DepositEndTime), update)
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
