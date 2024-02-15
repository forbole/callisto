package gov

import (
	"encoding/json"
	"fmt"

	tmtypes "github.com/cometbft/cometbft/types"

	"github.com/forbole/callisto/v4/types"

	gov "github.com/cosmos/cosmos-sdk/x/gov/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	"github.com/rs/zerolog/log"
)

// HandleGenesis implements modules.Module
func (m *Module) HandleGenesis(doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	log.Debug().Str("module", "gov").Msg("parsing genesis")

	// Read the genesis state
	var genStatev1beta1 govtypesv1.GenesisState
	err := m.cdc.UnmarshalJSON(appState[gov.ModuleName], &genStatev1beta1)
	if err != nil {
		return fmt.Errorf("error while reading gov genesis data: %s", err)
	}

	// Save the proposals
	err = m.saveGenesisProposals(genStatev1beta1.Proposals, doc)
	if err != nil {
		return fmt.Errorf("error while storing genesis governance proposals: %s", err)
	}

	// Save the params
	err = m.db.SaveGovParams(types.NewGovParams(genStatev1beta1.Params, doc.InitialHeight))
	if err != nil {
		return fmt.Errorf("error while storing genesis governance params: %s", err)
	}

	return nil
}

// saveGenesisProposals save proposals from genesis file
func (m *Module) saveGenesisProposals(slice govtypesv1.Proposals, genDoc *tmtypes.GenesisDoc) error {
	proposals := make([]types.Proposal, len(slice))
	tallyResults := make([]types.TallyResult, len(slice))
	deposits := make([]types.Deposit, len(slice))

	for index, proposal := range slice {
		// Since it's not possible to get the proposer, set it to nil
		proposals[index] = types.NewProposal(
			proposal.Id,
			proposal.Title,
			proposal.Summary,
			proposal.Metadata,
			proposal.Messages,
			proposal.Status.String(),
			*proposal.SubmitTime,
			*proposal.DepositEndTime,
			proposal.VotingStartTime,
			proposal.VotingEndTime,
			"",
		)

		tallyResults[index] = types.NewTallyResult(
			proposal.Id,
			proposal.FinalTallyResult.YesCount,
			proposal.FinalTallyResult.AbstainCount,
			proposal.FinalTallyResult.NoCount,
			proposal.FinalTallyResult.NoWithVetoCount,
			genDoc.InitialHeight,
		)

		deposits[index] = types.NewDeposit(
			proposal.Id,
			"",
			proposal.TotalDeposit,
			genDoc.GenesisTime,
			"",
			genDoc.InitialHeight,
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
