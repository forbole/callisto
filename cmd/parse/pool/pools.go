package pool

import (
	"fmt"

	parsecmdtypes "github.com/forbole/juno/v4/cmd/parse/types"
	"github.com/forbole/juno/v4/types/config"
	"github.com/spf13/cobra"

	"github.com/forbole/bdjuno/v4/database"
	"github.com/forbole/bdjuno/v4/modules/pool"
	modulestypes "github.com/forbole/bdjuno/v4/modules/types"
)

// poolCmd returns the Cobra command allowing to refresh x/pool values
func poolCmd(parseConfig *parsecmdtypes.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "all",
		Short: "Refresh all pools information in database",
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

			// Build pool module
			poolModule := pool.NewModule(sources.PoolSource, parseCtx.EncodingConfig.Codec, db)

			height, err := parseCtx.Node.LatestHeight()
			if err != nil {
				return fmt.Errorf("error while getting chain latest block height: %s", err)
			}

			err = poolModule.UpdatePools(height)
			if err != nil {
				return fmt.Errorf("error while updating pools: %s", err)
			}

			return nil
		},
	}
}
