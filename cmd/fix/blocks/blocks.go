package blocks

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/forbole/juno/v2/cmd/parse"
	"github.com/forbole/juno/v2/types/config"
	"github.com/spf13/cobra"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
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
	junomessages "github.com/forbole/juno/v2/modules/messages"
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

			var k int64 = consensusModule.GetStartingHeight()
			fmt.Printf("Starting height is %v ... \n", k)
			for ; k <= height; k++ {
				missingBlock := consensusModule.IsBlockMissing(k)
				if missingBlock {
					fmt.Printf("Refetching block %v ... \n", k)
					err = refreshBlock(parseCtx, k, consensusModule, sources, db)
					if err != nil {
						return err
					}
				}
			}

			return nil
		},
	}
}

func refreshBlock(parseCtx *parse.Context, blockHeight int64, consensusModule *consensus.Module, sources *modules.Sources, db *database.Db) error {
	// Get the block details
	block, blockResults, err := utils.QueryBlock(parseCtx.Node, blockHeight)
	if err != nil {
		return err
	}

	err = consensusModule.UpdateBlock(block, blockResults)

	if len(block.Block.Txs) != 0 {
		for _, tx := range block.Block.Txs {
			fmt.Printf("Refetching tx %v ... \n", strings.ToUpper(hex.EncodeToString(tx.Hash())))
			err = refreshTxs(parseCtx, consensusModule, block, blockResults, sources, db)
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

func refreshTxs(parseCtx *parse.Context, consensusModule *consensus.Module, block *tmctypes.ResultBlock, blockResults *tmctypes.ResultBlockResults, sources *modules.Sources, db *database.Db) error {

	for _, tx := range block.Block.Txs {
		// Get the tx details
		txDetails, err := parseCtx.Node.Tx(hex.EncodeToString(tx.Hash()))
		if err != nil {
			return err
		}

		// Handle messages
		for index, msg := range txDetails.GetMsgs() {
			// Register modules
			bankModule := bank.NewModule(junomessages.BankMessagesParser, sources.BankSource, parseCtx.EncodingConfig.Marshaler, db)
			govModule := gov.NewModule(parseCtx.EncodingConfig.Marshaler, sources.GovSource, nil, nil, nil, db)
			distrModule := distribution.NewModule(config.Cfg, sources.DistrSource, bankModule, db)
			historyModule := history.NewModule(config.Cfg.Chain, junomessages.BankMessagesParser, parseCtx.EncodingConfig.Marshaler, db)
			stakingModule := staking.NewModule(sources.StakingSource, bankModule, distrModule, historyModule, parseCtx.EncodingConfig.Marshaler, db)

			msg := msg
			mssg, err := json.Marshal(&msg)
			if err != nil {
				return fmt.Errorf("error while marshaling messages: %s", err)
			}

			var message []string

			message = append(message, string(mssg))

			fee, err := json.Marshal(&txDetails.AuthInfo.Fee)
			if err != nil {
				return fmt.Errorf("error while marshaling fees: %s", err)
			}

			logs, err := json.Marshal(&txDetails.Logs)
			if err != nil {
				return fmt.Errorf("error while marshaling logs: %s", err)
			}

			var signatures []string
			for _, signature := range txDetails.Signatures {
				signature := signature
				eachSignature, err := json.Marshal(&signature)
				if err != nil {
					return fmt.Errorf("error while marshaling signatures: %s", err)
				}
				signatures = append(signatures, string(eachSignature))
			}

			signer := txDetails.GetSigners()
			signers, err := json.Marshal(&signer)
			if err != nil {
				return fmt.Errorf("error while marshaling signers: %s  ERR: %s", signer, err)
			}

			switch msg.(type) {
			case *banktypes.MsgSend:
				blockErr := bankModule.HandleBlock(block, blockResults, nil, nil)
				if blockErr != nil {
					return fmt.Errorf("error when updatig MsgSend Handle Block: %s", err)
				}
				messageErr := bankModule.HandleMsg(index, msg, txDetails)
				if messageErr != nil {
					return fmt.Errorf("error when updatig MsgSend Handle Message: %s", err)
				}
				err = consensusModule.UpdateTxs(txDetails.TxHash, txDetails.Height, txDetails.Successful(), message, txDetails.GetBody().Memo, signatures, signers, string(fee), txDetails.GasWanted, txDetails.GasUsed, txDetails.RawLog, logs)
				if err != nil {
					return fmt.Errorf("error when updatig MsgSend \ntx: %v , txDetails: %v \ntxdetails gin: %v \nmsg: %s  : %s", tx, txDetails, message, message, err)
				}
			case *banktypes.MsgMultiSend:
				blockErr := bankModule.HandleBlock(block, blockResults, nil, nil)
				if blockErr != nil {
					return fmt.Errorf("error when updatig MsgMultiSend Handle Block: %s", err)
				}
				messageErr := bankModule.HandleMsg(index, msg, txDetails)
				if messageErr != nil {
					return fmt.Errorf("error when updatig MsgMultiSend Handle Message: %s", err)
				}
				err = consensusModule.UpdateTxs(txDetails.TxHash, txDetails.Height, txDetails.Successful(), message, txDetails.GetBody().Memo, signatures, signers, string(fee), txDetails.GasWanted, txDetails.GasUsed, txDetails.RawLog, logs)
				if err != nil {
					return fmt.Errorf("error when updatig MsgMultiSend: %s", err)
				}
			case *govtypes.MsgDeposit:
				blockErr := govModule.HandleBlock(block, blockResults, nil, nil)
				if blockErr != nil {
					return fmt.Errorf("error when updatig MsgDeposit Handle Block: %s", err)
				}
				messageErr := govModule.HandleMsg(index, msg, txDetails)
				if messageErr != nil {
					return fmt.Errorf("error when updatig MsgDeposit Handle Message: %s", err)
				}
				err = consensusModule.UpdateTxs(txDetails.TxHash, txDetails.Height, txDetails.Successful(), message, txDetails.GetBody().Memo, signatures, signers, string(fee), txDetails.GasWanted, txDetails.GasUsed, txDetails.RawLog, logs)
				if err != nil {
					return fmt.Errorf("error when updatig MsgDeposit: %s", err)
				}
			case *govtypes.MsgVote:
				blockErr := govModule.HandleBlock(block, blockResults, nil, nil)
				if blockErr != nil {
					return fmt.Errorf("error when updatig MsgVote Handle Block: %s", err)
				}
				messageErr := govModule.HandleMsg(index, msg, txDetails)
				if messageErr != nil {
					return fmt.Errorf("error when updatig MsgVote Handle Message: %s", err)
				}
				err = consensusModule.UpdateTxs(txDetails.TxHash, txDetails.Height, txDetails.Successful(), message, txDetails.GetBody().Memo, signatures, signers, string(fee), txDetails.GasWanted, txDetails.GasUsed, txDetails.RawLog, logs)
				if err != nil {
					return fmt.Errorf("error when updatig MsgVote: %s", err)
				}
			case *govtypes.MsgSubmitProposal:
				blockErr := govModule.HandleBlock(block, blockResults, nil, nil)
				if blockErr != nil {
					return fmt.Errorf("error when updatig MsgSubmitProposal Handle Block: %s", err)
				}
				messageErr := govModule.HandleMsg(index, msg, txDetails)
				if messageErr != nil {
					return fmt.Errorf("error when updatig MsgSubmitProposal Handle Message: %s", err)
				}
				err = consensusModule.UpdateTxs(txDetails.TxHash, txDetails.Height, txDetails.Successful(), message, txDetails.GetBody().Memo, signatures, signers, string(fee), txDetails.GasWanted, txDetails.GasUsed, txDetails.RawLog, logs)
				if err != nil {
					return fmt.Errorf("error when updatig MsgSubmitProposal: %s", err)
				}

			// case *slashingtypes.MsgUnjail:
			// 	slashingModule := slashing.NewModule(sources.SlashingSource, nil, nil)
			// 	err = slashingModule.HandleMsg(index, msg, txDetails)
			case *stakingtypes.MsgCreateValidator:
				messageErr := stakingModule.HandleMsg(index, msg, txDetails)
				if messageErr != nil {
					return fmt.Errorf("error when updatig MsgCreateValidator Handle Message: %s", err)
				}
				err = consensusModule.UpdateTxs(txDetails.TxHash, txDetails.Height, txDetails.Successful(), message, txDetails.GetBody().Memo, signatures, signers, string(fee), txDetails.GasWanted, txDetails.GasUsed, txDetails.RawLog, logs)
				if err != nil {
					return fmt.Errorf("error when updatig MsgCreateValidator tx %s", err)
				}
			case *stakingtypes.MsgBeginRedelegate:
				messageErr := stakingModule.HandleMsg(index, msg, txDetails)
				if messageErr != nil {
					return fmt.Errorf("error when updatig MsgBeginRedelegate Handle Message: %s", err)
				}
				err = consensusModule.UpdateTxs(txDetails.TxHash, txDetails.Height, txDetails.Successful(), message, txDetails.GetBody().Memo, signatures, signers, string(fee), txDetails.GasWanted, txDetails.GasUsed, txDetails.RawLog, logs)
				if err != nil {
					return fmt.Errorf("error when updatig MsgBeginRedelegate tx %s", err)
				}
			case *stakingtypes.MsgDelegate:
				messageErr := stakingModule.HandleMsg(index, msg, txDetails)
				if messageErr != nil {
					return fmt.Errorf("error when updatig MsgDelegate Handle Message: %s", err)
				}
				err = consensusModule.UpdateTxs(txDetails.TxHash, txDetails.Height, txDetails.Successful(), message, txDetails.GetBody().Memo, signatures, signers, string(fee), txDetails.GasWanted, txDetails.GasUsed, txDetails.RawLog, logs)
				if err != nil {
					return fmt.Errorf("error when updatig MsgDelegate tx %s", err)
				}
			case *stakingtypes.MsgEditValidator:
				messageErr := stakingModule.HandleMsg(index, msg, txDetails)
				if messageErr != nil {
					return fmt.Errorf("error when updatig MsgEditValidator Handle Message: %s", err)
				}
				err = consensusModule.UpdateTxs(txDetails.TxHash, txDetails.Height, txDetails.Successful(), message, txDetails.GetBody().Memo, signatures, signers, string(fee), txDetails.GasWanted, txDetails.GasUsed, txDetails.RawLog, logs)
				if err != nil {
					return fmt.Errorf("error when updatig MsgEditValidator tx %s", err)
				}
			case *stakingtypes.MsgUndelegate:
				messageErr := stakingModule.HandleMsg(index, msg, txDetails)
				if messageErr != nil {
					return fmt.Errorf("error when updatig MsgUndelegate Handle Message: %s", err)
				}
				err = consensusModule.UpdateTxs(txDetails.TxHash, txDetails.Height, txDetails.Successful(), message, txDetails.GetBody().Memo, signatures, signers, string(fee), txDetails.GasWanted, txDetails.GasUsed, txDetails.RawLog, logs)
				if err != nil {
					return fmt.Errorf("error when updatig MsgUndelegate tx %s", err)
				}
				// case *distrtypes.MsgFundCommunityPool:
			// 	fmt.Printf("%v Module name is %s", index, cosmosMsg)
			// case *distrtypes.MsgSetWithdrawAddress:
			// 	fmt.Printf("%v Module name is %s", index, cosmosMsg)
			// case *distrtypes.MsgWithdrawDelegatorReward:
			// 	fmt.Printf("%v Module name is %s", index, cosmosMsg)
			// case *distrtypes.MsgWithdrawValidatorCommission:
			// 	fmt.Printf("%v Module name is %s", index, cosmosMsg)
			default:
				return nil

			}
		}
	}
	return nil
}
