package ccv

import (
	"fmt"

	"github.com/forbole/bdjuno/v4/database"
	provider "github.com/forbole/bdjuno/v4/modules/ccv/provider"
	modulestypes "github.com/forbole/bdjuno/v4/modules/types"
	parsecmdtypes "github.com/forbole/juno/v4/cmd/parse/types"
	"github.com/forbole/juno/v4/types/config"
	"github.com/spf13/cobra"
)

// ccvAllChainsCmd returns a Cobra command that allows to refresh consumer chains info in database
func ccvAllChainsCmd(parseConfig *parsecmdtypes.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "chains",
		Short: "Fix the information about active consumer chains",
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
			ccvProviderModule := provider.NewModule(sources.CcvProviderSource, parseCtx.EncodingConfig.Marshaler, db)

			// Get latest height
			height, err := parseCtx.Node.LatestHeight()
			if err != nil {
				return fmt.Errorf("error while getting latest block height: %s", err)
			}

			err = ccvProviderModule.UpdateAllConsumerChains(height)
			if err != nil {
				return err
			}

			return nil
		},
	}
}
