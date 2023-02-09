package gov

import (
	"encoding/json"
	"fmt"

	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/forbole/bdjuno/v4/types"

	gov "github.com/cosmos/cosmos-sdk/x/gov/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	"github.com/rs/zerolog/log"
)

// HandleGenesis implements modules.Module
func (m *Module) HandleGenesis(doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	log.Debug().Str("module", "gov").Msg("parsing genesis")

	// Read v1 genesis state
	var genStateV1 govtypesv1.GenesisState
	err := m.cdc.UnmarshalJSON(appState[gov.ModuleName], &genStateV1)
	if err != nil {
		return fmt.Errorf("error while reading gov genesis data v1: %s", err)
	}

	// Save the proposals
	err = m.saveGenesisProposals(genStateV1.Proposals, doc)
	if err != nil {
		return fmt.Errorf("error while storing genesis governance proposals: %s", err)
	}

	// Save the params
	err = m.db.SaveGovParams(types.NewGovParams(
		types.NewVotingParams(genStateV1.VotingParams),
		types.NewDepositParam(genStateV1.DepositParams),
		types.NewTallyParams(genStateV1.TallyParams),
		doc.InitialHeight,
	))
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
		// Set proposal route, type, content and proposer values to nil,
		// since it's not possible to parse them correctly (sdk v0.46.7)
		proposals[index] = types.NewProposal(
			proposal.GetId(),
			"",
			"",
			nil,
			proposal.Status.String(),
			*proposal.SubmitTime,
			*proposal.DepositEndTime,
			*proposal.VotingStartTime,
			*proposal.VotingEndTime,
			"",
		)
		tallyResults[index] = types.NewTallyResult(
			proposal.GetId(),
			proposal.FinalTallyResult.YesCount,
			proposal.FinalTallyResult.AbstainCount,
			proposal.FinalTallyResult.NoCount,
			proposal.FinalTallyResult.NoWithVetoCount,
			genDoc.InitialHeight,
		)

		deposits[index] = types.NewDeposit(
			proposal.GetId(),
			"",
			proposal.TotalDeposit,
			genDoc.GenesisTime,
			genDoc.InitialHeight,
		)
	}

	// Save the proposals
	err := m.db.SaveGenesisProposals(proposals)
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
