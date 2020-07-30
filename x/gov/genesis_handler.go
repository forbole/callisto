package gov

import (
	"encoding/json"
	"fmt"
	"time"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/desmos-labs/juno/parse/worker"
	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/x/gov/types"
	"github.com/rs/zerolog/log"
	tmtypes "github.com/tendermint/tendermint/types"
)
//genesisDoc.GenesisTime

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

	if err := saveProposals(genState.Proposals, genesisDoc, bigDipperDb); err != nil {
		return err
	}
	return nil
}

func saveProposals(proposals gov.Proposals, genesisDoc *tmtypes.GenesisDoc, db database.BigDipperDb)error {
	bdproposals := make([]types.Proposal,len(proposals))
	bdTallyResult := make([]types.TallyResult,len(proposals))
	for _,proposal :=range(proposals){
		submitTime,err := time.Parse(time.RFC3339,proposal.SubmitTime.String())
		if err !=nil{
			return err
		}
		depositEndTime,err := time.Parse(time.RFC3339,proposal.DepositEndTime.String())
		if err !=nil{
			return err
		}
		votingStartTime,err := time.Parse(time.RFC3339,proposal.VotingStartTime.String())
		if err !=nil{
			return err
		}
		votingEndTime,err := time.Parse(time.RFC3339,proposal.VotingEndTime.String())
		if err!=nil{
			return err
		}

		bdproposals = append(bdproposals,types.NewProposal(proposal.GetTitle(),proposal.GetDescription(),proposal.ProposalRoute(),proposal.ProposalType,proposal.ProposalID,proposal.Status.String(),
							submitTime,depositEndTime,proposal.TotalDeposit,votingStartTime,votingEndTime))
		//see if the current status of proposal 
		
		bdTallyResult = append(bdTallyResult,types.NewTallyResult(proposal.ProposalID,proposal.FinalTallyResult.Yes.BigInt(),proposal.FinalTallyResult.Abstain.BigInt(),proposal.FinalTallyResult.No.BigInt(),
								proposal.FinalTallyResult.NoWithVeto.BigInt(),"0",timestamp))
	}
	return nil
}
