package periodictask

import (
	"fmt"

	modulestypes "github.com/forbole/bdjuno/v3/modules/types"

	parsecmdtypes "github.com/forbole/juno/v3/cmd/parse/types"
	"github.com/forbole/juno/v3/modules/messages"
	"github.com/forbole/juno/v3/types/config"
	"github.com/spf13/cobra"

	"github.com/forbole/bdjuno/v3/database"
	"github.com/forbole/bdjuno/v3/modules/bank"
)

// bankCmd returns the Cobra command allowing to refresh data that's obtained from x/bank periodic tasks
func bankCmd(parseConfig *parsecmdtypes.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "bank",
		Short: "Run x/bank periodic task",
		RunE: func(cmd *cobra.Command, args []string) error {
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

			// Build bank module
			msgParser := messages.JoinMessageParsers(messages.CosmosMessageAddressesParser)
			bankModule := bank.NewModule(msgParser, sources.BankSource, parseCtx.EncodingConfig.Marshaler, db)

			err = bankModule.UpdateSupply()
			if err != nil {
				return fmt.Errorf("error while getting latest bank supply: %s", err)
			}

			return nil
		},
	}
}
