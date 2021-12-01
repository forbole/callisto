package blocks

import (
	"encoding/hex"
	"fmt"
	"strings"

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
	"github.com/forbole/juno/v2/types/config"
	"github.com/spf13/cobra"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
)

// blocksCmd returns a Cobra command that allows to fix missing blocks in database
func blocksCmd(parseConfig *parse.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "refetch",
		Short: "Fix missing blocks and transactions in database from the start height",
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

			// Register modules
			bankModule := bank.NewModule(junomessages.BankMessagesParser, sources.BankSource, parseCtx.EncodingConfig.Marshaler, db)
			distrModule := distribution.NewModule(config.Cfg, sources.DistrSource, bankModule, db)
			historyModule := history.NewModule(config.Cfg.Chain, junomessages.BankMessagesParser, parseCtx.EncodingConfig.Marshaler, db)
			stakingModule := staking.NewModule(sources.StakingSource, bankModule, distrModule, historyModule, nil, parseCtx.EncodingConfig.Marshaler, db)
			govModule := gov.NewModule(parseCtx.EncodingConfig.Marshaler, sources.GovSource, nil, bankModule, distrModule, nil, nil, stakingModule, db)

			// Build the consensus module
			consensusModule := consensus.NewModule(bankModule, distrModule, govModule, stakingModule, db)

			// Get latest height
			height, err := parseCtx.Node.LatestHeight()
			if err != nil {
				return fmt.Errorf("error while getting chain latest block height: %s", err)
			}

			k := config.Cfg.Parser.StartHeight
			fmt.Printf("Refetching missing blocks and transactions from height %d ... \n", k)
			for ; k <= height; k++ {
				missingBlock := consensusModule.IsBlockMissing(k)
				if missingBlock {
					fmt.Printf("Refetching block %d ... \n", k)
					err = refreshBlock(parseCtx, k, consensusModule)
					if err != nil {
						return err
					}
				}
			}

			return nil
		},
	}
}

func refreshBlock(parseCtx *parse.Context, blockHeight int64, consensusModule *consensus.Module) error {
	// Get the block details
	block, err := utils.QueryBlock(parseCtx.Node, blockHeight)
	if err != nil {
		return err
	}

	err = consensusModule.UpdateBlock(block)

	if len(block.Block.Txs) != 0 {
		for _, tx := range block.Block.Txs {
			fmt.Printf("Refetching tx %v ... \n", strings.ToUpper(hex.EncodeToString(tx.Hash())))
			err = refreshTxs(parseCtx, consensusModule, block)
			if err != nil {
				return fmt.Errorf("error when updatig tx %s", err)
			}
		}
	}
	if err != nil {
		return fmt.Errorf("error while updating block %d: %s", blockHeight, err)
	}

	return nil
}

func refreshTxs(parseCtx *parse.Context, consensusModule *consensus.Module, block *tmctypes.ResultBlock) error {

	for _, tx := range block.Block.Txs {
		// Get the tx details
		txDetails, err := parseCtx.Node.Tx(hex.EncodeToString(tx.Hash()))
		if err != nil {
			return fmt.Errorf("error while encoding tx to string %s", err)
		}

		err = consensusModule.HandleMessages(txDetails)
		if err != nil {
			return fmt.Errorf("error when updating tx message tx %s", err)
		}
	}

	return nil
}
