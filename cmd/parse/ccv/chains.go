package ccv

import (
	"fmt"

	"github.com/forbole/bdjuno/v4/database"
	consumer "github.com/forbole/bdjuno/v4/modules/ccv/consumer"
	modulestypes "github.com/forbole/bdjuno/v4/modules/types"
	parsecmdtypes "github.com/forbole/juno/v4/cmd/parse/types"
	"github.com/forbole/juno/v4/types/config"
	"github.com/spf13/cobra"
)

// consumerChainsCmd returns a Cobra command that allows to refresh consumer chains info in database
func consumerChainsCmd(parseConfig *parsecmdtypes.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "consumer-chains",
		Short: "Fix the consumer chains details stored in database",
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

			// Build the ccv provider module
			ccvProviderModule := consumer.NewModule(sources.CcvConsumerSource, parseCtx.EncodingConfig.Marshaler, db)

			// Get latest height
			height, err := parseCtx.Node.LatestHeight()
			if err != nil {
				return fmt.Errorf("error while getting latest block height: %s", err)
			}

			err = ccvProviderModule.UpdateAllConsumerChains(height)
			if err != nil {
				return fmt.Errorf("error while updating all consumer chains at height: %s", err)
			}

			return nil
		},
	}
}
