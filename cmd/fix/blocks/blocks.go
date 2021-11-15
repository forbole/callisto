package blocks

import (
	"encoding/hex"
	"fmt"
	"strings"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/forbole/bdjuno/v2/database"
	"github.com/forbole/bdjuno/v2/modules"
	"github.com/forbole/bdjuno/v2/modules/bank"
	"github.com/forbole/bdjuno/v2/modules/consensus"
	"github.com/forbole/bdjuno/v2/modules/distribution"
	"github.com/forbole/bdjuno/v2/modules/gov"
	"github.com/forbole/bdjuno/v2/modules/history"
	"github.com/forbole/bdjuno/v2/modules/staking"
	"github.com/forbole/bdjuno/v2/utils"
	"github.com/forbole/juno/v2/cmd/parse"
	junomessages "github.com/forbole/juno/v2/modules/messages"
	juno "github.com/forbole/juno/v2/types"
	"github.com/forbole/juno/v2/types/config"
	"github.com/spf13/cobra"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
)

// blocksCmd returns a Cobra command that allows to fix missing blocks in database
func blocksCmd(parseConfig *parse.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "refetch",
		Short: "Fix missing blocks in database from the latest known height",
		RunE: func(cmd *cobra.Command, args []string) error {
			parseCtx, err := parse.GetParsingContext(parseConfig)
			if err != nil {
				return err
			}

			sources, err := modules.BuildSources(config.Cfg.Node, parseCtx.EncodingConfig)
			if err != nil {
				return err
			}

			// Get the database
			db := database.Cast(parseCtx.Database)

			// Build the consensus module
			consensusModule := consensus.NewModule(config.Cfg, db)

			// Get latest height
			height, err := parseCtx.Node.LatestHeight()
			if err != nil {
				return fmt.Errorf("error while getting chain latest block height: %s", err)
			}

			k := consensusModule.GetStartingHeight()
			fmt.Printf("Starting height is %v ... \n", k)
			for ; k <= height; k++ {
				missingBlock := consensusModule.IsBlockMissing(k)
				if missingBlock {
					fmt.Printf("Refetching block %v ... \n", k)
					err = refreshBlock(parseCtx, sources, k, consensusModule, db)
					if err != nil {
						return err
					}
				}
			}

			return nil
		},
	}
}

func refreshBlock(parseCtx *parse.Context, sources *modules.Sources, blockHeight int64, consensusModule *consensus.Module, db *database.Db) error {
	// Get the block details
	block, blockResults, err := utils.QueryBlock(parseCtx.Node, blockHeight)
	if err != nil {
		return err
	}

	err = consensusModule.UpdateBlock(block, blockResults)

	if len(block.Block.Txs) != 0 {
		for _, tx := range block.Block.Txs {
			fmt.Printf("Refetching tx %v ... \n", strings.ToUpper(hex.EncodeToString(tx.Hash())))
			err = refreshTxs(parseCtx, sources, consensusModule, block, db)
			if err != nil {
				return fmt.Errorf("error when updatig tx %s", err)
			}
		}
	}
	if err != nil {
		return fmt.Errorf("error while updating block %v: %s", blockHeight, err)
	}

	return nil
}

func refreshTxs(parseCtx *parse.Context, sources *modules.Sources, consensusModule *consensus.Module, block *tmctypes.ResultBlock, db *database.Db) error {

	// Register modules
	bankModule := bank.NewModule(junomessages.BankMessagesParser, sources.BankSource, parseCtx.EncodingConfig.Marshaler, db)
	govModule := gov.NewModule(parseCtx.EncodingConfig.Marshaler, sources.GovSource, nil, nil, nil, db)
	distrModule := distribution.NewModule(config.Cfg, sources.DistrSource, bankModule, db)
	historyModule := history.NewModule(config.Cfg.Chain, junomessages.BankMessagesParser, parseCtx.EncodingConfig.Marshaler, db)
	stakingModule := staking.NewModule(sources.StakingSource, bankModule, distrModule, historyModule, parseCtx.EncodingConfig.Marshaler, db)

	for _, tx := range block.Block.Txs {
		var txDetails *juno.Tx
		// Get the tx details
		txDetails, err := parseCtx.Node.Tx(hex.EncodeToString(tx.Hash()))
		if err != nil {
			return fmt.Errorf("error while encoding tx to string %s", err)
		}

		// Handle messages
		for index, msg := range txDetails.GetMsgs() {
			switch msg.(type) {
			case *banktypes.MsgSend:
				messageErr := bankModule.HandleMsg(index, msg, txDetails)
				if messageErr != nil {
					return fmt.Errorf("error when updatig MsgSend Handle Message: %s", err)
				}
			case *banktypes.MsgMultiSend:
				messageErr := bankModule.HandleMsg(index, msg, txDetails)
				if messageErr != nil {
					return fmt.Errorf("error when updatig MsgMultiSend Handle Message: %s", err)
				}
			case *govtypes.MsgDeposit:
				messageErr := govModule.HandleMsg(index, msg, txDetails)
				if messageErr != nil {
					return fmt.Errorf("error when updatig MsgDeposit Handle Message: %s", err)
				}
			case *govtypes.MsgVote:
				messageErr := govModule.HandleMsg(index, msg, txDetails)
				if messageErr != nil {
					return fmt.Errorf("error when updatig MsgVote Handle Message: %s", err)
				}
			case *govtypes.MsgSubmitProposal:
				messageErr := govModule.HandleMsg(index, msg, txDetails)
				if messageErr != nil {
					return fmt.Errorf("error when updatig MsgSubmitProposal Handle Message: %s", err)
				}
			// case *slashingtypes.MsgUnjail:
			// 	err = consensusModule.UpdateTxs2(index, txDetails)
			// 	if err != nil {
			// 		return fmt.Errorf("error when updatig MsgUnjail: %s", err)
			// 	}
			case *stakingtypes.MsgCreateValidator:
				messageErr := stakingModule.HandleMsg(index, msg, txDetails)
				if messageErr != nil {
					return fmt.Errorf("error when updatig MsgCreateValidator Handle Message: %s", err)
				}
			case *stakingtypes.MsgBeginRedelegate:
				messageErr := stakingModule.HandleMsg(index, msg, txDetails)
				if messageErr != nil {
					return fmt.Errorf("error when updatig MsgBeginRedelegate Handle Message: %s", err)
				}
			case *stakingtypes.MsgDelegate:
				messageErr := stakingModule.HandleMsg(index, msg, txDetails)
				if messageErr != nil {
					return fmt.Errorf("error when updatig MsgDelegate Handle Message: %s", err)
				}
			case *stakingtypes.MsgEditValidator:
				messageErr := stakingModule.HandleMsg(index, msg, txDetails)
				if messageErr != nil {
					return fmt.Errorf("error when updatig MsgEditValidator Handle Message: %s", err)
				}
			case *stakingtypes.MsgUndelegate:
				messageErr := stakingModule.HandleMsg(index, msg, txDetails)
				if messageErr != nil {
					return fmt.Errorf("error when updatig MsgUndelegate Handle Message: %s", err)
				}
			case *distrtypes.MsgFundCommunityPool:
				messageErr := distrModule.HandleMsg(index, msg, txDetails)
				if messageErr != nil {
					return fmt.Errorf("error when updatig MsgFundCommunityPool Handle Message: %s", err)
				}
			case *distrtypes.MsgSetWithdrawAddress:
				messageErr := distrModule.HandleMsg(index, msg, txDetails)
				if messageErr != nil {
					return fmt.Errorf("error when updatig MsgSetWithdrawAddress Handle Message: %s", err)
				}
			case *distrtypes.MsgWithdrawDelegatorReward:
				messageErr := distrModule.HandleMsg(index, msg, txDetails)
				if messageErr != nil {
					return fmt.Errorf("error when updatig MsgWithdrawDelegatorReward Handle Message: %s", err)
				}
			case *distrtypes.MsgWithdrawValidatorCommission:
				messageErr := distrModule.HandleMsg(index, msg, txDetails)
				if messageErr != nil {
					return fmt.Errorf("error when updatig MsgWithdrawValidatorCommission Handle Message: %s", err)
				}
			default:
				return nil
			}

			err = consensusModule.UpdateTxs(index, txDetails)
			if err != nil {
				return fmt.Errorf("error when updatig transactions tx %s", err)
			}
		}
	}
	return nil
}
