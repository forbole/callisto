package feegrant

import (
	"encoding/hex"
	"fmt"

	"github.com/forbole/bdjuno/v2/modules/feegrant"
	"github.com/forbole/bdjuno/v2/utils"

	"github.com/forbole/bdjuno/v2/database"
	"github.com/forbole/juno/v2/cmd/parse"
	"github.com/forbole/juno/v2/types/config"
	"github.com/spf13/cobra"

	feegranttypes "github.com/cosmos/cosmos-sdk/x/feegrant"
)

// allowanceCmd returns the Cobra command allowing to fix all things related to fee grant allowance
func allowanceCmd(parseConfig *parse.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "allowance",
		Short: "Fix the granted and revoked allowances to the latest height",
		RunE: func(cmd *cobra.Command, args []string) error {
			parseCtx, err := parse.GetParsingContext(parseConfig)
			if err != nil {
				return err
			}

			// Get the database
			db := database.Cast(parseCtx.Database)

			// Build feegrant module
			feegrantModule := feegrant.NewModule(parseCtx.EncodingConfig.Marshaler, db)

			// Get latest height
			height, err := parseCtx.Node.LatestHeight()
			if err != nil {
				return fmt.Errorf("error while getting chain latest block height: %s", err)
			}

			startHeight := config.Cfg.Parser.StartHeight

			// Handle messages realated to feegrant allowance of each block from the start height
			for h := startHeight; h <= height; h++ {
				err = refreshAllowance(parseCtx, h, feegrantModule)
				if err != nil {
					fmt.Printf("error while refreshing allowance: %s \n", err)
					// fmt.Errorf("error while refreshing allowance: %s", err)
				}
			}

			return nil
		},
	}
}

func refreshAllowance(parseCtx *parse.Context, blockHeight int64, feegrantModule *feegrant.Module) error {
	// Get the block details
	block, err := utils.QueryBlock(parseCtx.Node, blockHeight)
	if err != nil {
		return err
	}

	if len(block.Block.Txs) != 0 {
		for _, tx := range block.Block.Txs {
			// Get the tx details
			fmt.Println("hash:", hex.EncodeToString(tx.Hash()))

			junoTx, err := parseCtx.Node.Tx(hex.EncodeToString(tx.Hash()))
			if err != nil {
				return fmt.Errorf("error while getting tx details: %s", err)
			}

			// Handle the MsgDeposit messages
			for _, msg := range junoTx.GetMsgs() {
				if msgGrantAllowance, ok := msg.(*feegranttypes.MsgGrantAllowance); ok {
					fmt.Println("handling MsgGrantAllowance")

					err = feegrantModule.HandleMsgGrantAllowance(junoTx, msgGrantAllowance)
					if err != nil {
						return fmt.Errorf("error while handling MsgGrantAllowance: %s", err)
					}
				}

				if msgRevokeAllowance, ok := msg.(*feegranttypes.MsgRevokeAllowance); ok {
					fmt.Println("handling MsgRevokeAllowance")

					err = feegrantModule.HandleMsgRevokeAllowance(msgRevokeAllowance)
					if err != nil {
						return fmt.Errorf("error while handling msgRevokeAllowance: %s", err)
					}
				}

			}

		}
	}

	return nil
}
