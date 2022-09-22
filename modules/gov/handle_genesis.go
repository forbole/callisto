package gov

import (
	"encoding/json"
	"fmt"

	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/forbole/bdjuno/v3/types"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/rs/zerolog/log"
)

// HandleGenesis implements modules.Module
func (m *Module) HandleGenesis(doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	log.Debug().Str("module", "gov").Msg("parsing genesis")

	// Read the genesis state
	var genState govtypes.GenesisState
	err := m.cdc.UnmarshalJSON(appState[govtypes.ModuleName], &genState)
	if err != nil {
		return fmt.Errorf("error while reading gov genesis data: %s", err)
	}

	// Save the proposals
	err = m.saveProposals(genState.Proposals)
	if err != nil {
		return fmt.Errorf("error while storing genesis governance proposals: %s", err)
	}

	// Save the params
	err = m.db.SaveGovParams(types.NewGovParams(
		types.NewVotingParams(genState.VotingParams),
		types.NewDepositParam(genState.DepositParams),
		types.NewTallyParams(genState.TallyParams),
		doc.InitialHeight,
	))
	if err != nil {
		return fmt.Errorf("error while storing genesis governance params: %s", err)
	}

	return nil
}

// saveProposals save proposals from genesis file
func (m *Module) saveProposals(slice govtypes.Proposals) error {
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
			proposal.FinalTallyResult.Yes.String(),
			proposal.FinalTallyResult.Abstain.String(),
			proposal.FinalTallyResult.No.String(),
			proposal.FinalTallyResult.NoWithVeto.String(),
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
	err := m.db.SaveProposals(proposals)
	if err != nil {
		return err
	}

	// Save the deposits
	err = m.db.SaveDeposits(deposits)
	if err != nil {
		return err
	}

	// Save the tally results
	return m.db.SaveTallyResults(tallyResults)
}
