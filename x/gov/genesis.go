package gov

import (
	"encoding/json"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/desmos-labs/juno/client"
	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/x/gov/types"
	"github.com/rs/zerolog/log"
	tmtypes "github.com/tendermint/tendermint/types"
)

func HandleGenesis(
	genesisDoc *tmtypes.GenesisDoc, appState map[string]json.RawMessage,
	cdc *codec.Codec, cp *client.Proxy, db *database.BigDipperDb,
) error {
	log.Debug().Str("module", "gov").Msg("parsing genesis")

	// Read the genesis state
	var genState gov.GenesisState
	err := cdc.UnmarshalJSON(appState[gov.ModuleName], &genState)
	if err != nil {
		return err
	}

	// Save the proposals
	err = saveProposals(genesisDoc.GenesisTime, genState.Proposals, cp, db)
	if err != nil {
		return err
	}

	return nil
}

// saveProposals save proposals from genesis file
func saveProposals(genTime time.Time, p gov.Proposals, cp *client.Proxy, db *database.BigDipperDb) error {
	proposals := make([]types.Proposal, len(p))
	tallyResults := make([]types.TallyResult, len(p))
	deposits := make([]types.Deposit, len(p))

	for index, proposal := range p {
		// Since it's not possible to get the proposer, set it to nil
		proposals[index] = types.NewProposal(
			proposal.GetTitle(),
			proposal.GetDescription(),
			proposal.ProposalRoute(),
			proposal.ProposalType(),
			proposal.ProposalID,
			proposal.Status,
			proposal.SubmitTime,
			proposal.DepositEndTime,
			proposal.VotingStartTime,
			proposal.VotingEndTime,
			nil,
		)

		tallyResults[index] = types.NewTallyResult(
			proposal.ProposalID,
			proposal.FinalTallyResult.Yes.Int64(),
			proposal.FinalTallyResult.Abstain.Int64(),
			proposal.FinalTallyResult.No.Int64(),
			proposal.FinalTallyResult.NoWithVeto.Int64(),
			1,
			genTime,
		)

		deposits[index] = types.NewDeposit(
			proposal.ProposalID,
			sdk.AccAddress{},
			proposal.TotalDeposit,
			proposal.TotalDeposit,
			1,
			genTime,
		)

		// Update the proposal status when the voting period or deposit period ends
		update := UpdateProposal(proposal.ProposalID, cp, db)
		if proposal.Status.String() == "VotingPeriod" {
			time.AfterFunc(time.Since(proposal.VotingEndTime), update)
		} else if proposal.Status.String() == "DepositPeriod" {
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
