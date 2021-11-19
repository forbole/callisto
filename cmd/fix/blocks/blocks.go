package blocks

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/forbole/bdjuno/v2/database"
	"github.com/forbole/bdjuno/v2/modules"
	"github.com/forbole/bdjuno/v2/modules/consensus"
	"github.com/forbole/bdjuno/v2/utils"
	"github.com/forbole/juno/v2/cmd/parse"
	juno "github.com/forbole/juno/v2/types"
	"github.com/forbole/juno/v2/types/config"
	"github.com/spf13/cobra"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
)

// blocksCmd returns a Cobra command that allows to fix missing blocks in database
func blocksCmd(parseConfig *parse.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "refetch",
		Short: "Fix missing blocks and transactions in database from the latest known height",
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
			consensusModule := consensus.NewModule(config.Cfg, nil, nil, nil, nil, nil, db)

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
			err = refreshTxs(parseCtx, consensusModule, block)
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

func refreshTxs(parseCtx *parse.Context, consensusModule *consensus.Module, block *tmctypes.ResultBlock) error {

	for index, tx := range block.Block.Txs {
		var txDetails *juno.Tx
		// Get the tx details
		txDetails, err := parseCtx.Node.Tx(hex.EncodeToString(tx.Hash()))
		if err != nil {
			return fmt.Errorf("error while encoding tx to string %s", err)
		}

		err = consensusModule.HandleMessages(txDetails)
		if err != nil {
			return fmt.Errorf("error when updating tx message tx %s", err)
		}
		err = consensusModule.UpdateTxs(index, txDetails)
		if err != nil {
			return fmt.Errorf("error when updatig transactions tx %s", err)
		}
	}

	return nil
}
