package gov

import (
	"encoding/hex"
	"fmt"
	"strconv"

	modulestypes "github.com/forbole/bdjuno/v3/modules/types"
	"github.com/rs/zerolog/log"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	parsecmdtypes "github.com/forbole/juno/v3/cmd/parse/types"
	"github.com/forbole/juno/v3/types/config"
	"github.com/spf13/cobra"

	"github.com/forbole/juno/v3/parser"

	"github.com/forbole/bdjuno/v3/database"
	"github.com/forbole/bdjuno/v3/modules/distribution"
	"github.com/forbole/bdjuno/v3/modules/gov"
	"github.com/forbole/bdjuno/v3/modules/mint"
	"github.com/forbole/bdjuno/v3/modules/slashing"
	"github.com/forbole/bdjuno/v3/modules/staking"
	"github.com/forbole/bdjuno/v3/utils"
)

// proposalCmd returns the Cobra command allowing to fix all things related to a proposal
func proposalCmd(parseConfig *parsecmdtypes.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "proposal [id]",
		Short: "Get the description, votes and everything related to a proposal given its id",
		RunE: func(cmd *cobra.Command, args []string) error {
			proposalID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			parseCtx, err := parsecmdtypes.GetParserContext(config.Cfg, parseConfig)
			if err != nil {
				return err
			}

			sources, err := modulestypes.BuildSources(config.Cfg.Node, parseCtx.EncodingConfig)
			if err != nil {
				return err
			}

			// Get the database
			db := database.Cast(parseCtx.Database)

			// Build expected modules of gov modules for handleParamChangeProposal
			distrModule := distribution.NewModule(sources.DistrSource, parseCtx.EncodingConfig.Marshaler, db)
			mintModule := mint.NewModule(sources.MintSource, parseCtx.EncodingConfig.Marshaler, db)
			slashingModule := slashing.NewModule(sources.SlashingSource, parseCtx.EncodingConfig.Marshaler, db)
			stakingModule := staking.NewModule(sources.StakingSource, parseCtx.EncodingConfig.Marshaler, db)

			// Build the gov module
			govModule := gov.NewModule(sources.GovSource, nil, distrModule, mintModule, slashingModule, stakingModule, parseCtx.EncodingConfig.Marshaler, db)

			err = refreshProposalDetails(parseCtx, proposalID, govModule)
			if err != nil {
				return err
			}

			err = refreshProposalDeposits(parseCtx, proposalID, govModule)
			if err != nil {
				return err
			}

			err = refreshProposalVotes(parseCtx, proposalID, govModule)
			if err != nil {
				return err
			}

			// Update the proposal to the latest status
			height, err := parseCtx.Node.LatestHeight()
			if err != nil {
				return fmt.Errorf("error while getting chain latest block height: %s", err)
			}

			err = govModule.UpdateProposal(height, proposalID)
			if err != nil {
				return err
			}

			return nil
		},
	}
}

func refreshProposalDetails(parseCtx *parser.Context, proposalID uint64, govModule *gov.Module) error {
	log.Debug().Msg("refreshing proposal details")

	// Get the tx that created the proposal
	txs, err := utils.QueryTxs(parseCtx.Node, fmt.Sprintf("submit_proposal.proposal_id=%d", proposalID))
	if err != nil {
		return err
	}

	if len(txs) > 1 {
		return fmt.Errorf("expecting only one create proposal transaction, found %d", len(txs))
	}

	// Get the tx details
	tx, err := parseCtx.Node.Tx(hex.EncodeToString(txs[0].Tx.Hash()))
	if err != nil {
		return err
	}

	// Handle the MsgSubmitProposal messages
	for index, msg := range tx.GetMsgs() {
		if _, ok := msg.(*govtypes.MsgSubmitProposal); !ok {
			continue
		}

		err = govModule.HandleMsg(index, msg, tx)
		if err != nil {
			return fmt.Errorf("error while handling MsgSubmitProposal: %s", err)
		}
	}

	return nil
}

func refreshProposalDeposits(parseCtx *parser.Context, proposalID uint64, govModule *gov.Module) error {
	log.Debug().Msg("refreshing proposal deposits")

	// Get the tx that deposited to the proposal
	txs, err := utils.QueryTxs(parseCtx.Node, fmt.Sprintf("proposal_deposit.proposal_id=%d", proposalID))
	if err != nil {
		return err
	}

	for _, tx := range txs {
		// Get the tx details
		junoTx, err := parseCtx.Node.Tx(hex.EncodeToString(tx.Tx.Hash()))
		if err != nil {
			return err
		}

		// Handle the MsgDeposit messages
		for index, msg := range junoTx.GetMsgs() {
			if _, ok := msg.(*govtypes.MsgDeposit); !ok {
				continue
			}

			err = govModule.HandleMsg(index, msg, junoTx)
			if err != nil {
				return fmt.Errorf("error while handling MsgDeposit: %s", err)
			}
		}
	}

	return nil
}

func refreshProposalVotes(parseCtx *parser.Context, proposalID uint64, govModule *gov.Module) error {
	log.Debug().Msg("refreshing proposal votes")

	// Get the tx that voted the proposal
	txs, err := utils.QueryTxs(parseCtx.Node, fmt.Sprintf("proposal_vote.proposal_id=%d", proposalID))
	if err != nil {
		return err
	}

	for _, tx := range txs {
		// Get the tx details
		junoTx, err := parseCtx.Node.Tx(hex.EncodeToString(tx.Tx.Hash()))
		if err != nil {
			return err
		}

		// Handle the MsgVote messages
		for index, msg := range junoTx.GetMsgs() {
			if _, ok := msg.(*govtypes.MsgVote); !ok {
				continue
			}

			err = govModule.HandleMsg(index, msg, junoTx)
			if err != nil {
				return fmt.Errorf("error while handling MsgVote: %s", err)
			}
		}
	}

	return nil
}
