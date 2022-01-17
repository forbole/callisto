package blocks

import (
	"encoding/hex"
	"fmt"

	ibchanneltypes "github.com/cosmos/ibc-go/modules/core/04-channel/types"
	"github.com/forbole/juno/v2/modules/messages"

	"github.com/forbole/bdjuno/v2/database"
	"github.com/forbole/bdjuno/v2/modules"
	"github.com/forbole/bdjuno/v2/modules/bank"
	"github.com/forbole/bdjuno/v2/utils"

	"github.com/forbole/juno/v2/cmd/parse"

	"github.com/spf13/cobra"

	"github.com/forbole/juno/v2/types/config"
)

// receivedCmd returns a Cobra command that allows to fix balances for accounts that received incoming IBC transfers
func receivedCmd(parseConfig *parse.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "received",
		Short: "Fix balances for accounts that received incoming IBC transfers",
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

			bankModule := bank.NewModule(messages.CosmosMessageAddressesParser, sources.BankSource, parseCtx.EncodingConfig.Marshaler, db)
			err = refreshIBCReceivePacket(parseCtx, bankModule)
			if err != nil {
				return err
			}

			return nil
		},
	}
}

func refreshIBCReceivePacket(parseCtx *parse.Context, bankModule *bank.Module) error {
	txs, err := utils.QueryTxs(parseCtx.Node, "recv_packet.packet_dst_port = 'transfer'")
	if err != nil {
		return err
	}

	for _, transaction := range txs {
		// Get the tx details
		tx, err := parseCtx.Node.Tx(hex.EncodeToString(transaction.Tx.Hash()))
		if err != nil {
			return err
		}

		// Handle the MsgSubmitProposal messages
		for index, msg := range tx.GetMsgs() {
			if _, ok := msg.(*ibchanneltypes.MsgRecvPacket); !ok {
				continue
			}

			err = bankModule.HandleMsg(index, msg, tx)
			if err != nil {
				return fmt.Errorf("error while handling MsgSubmitProposal: %s", err)
			}
		}
	}

	return nil
}
