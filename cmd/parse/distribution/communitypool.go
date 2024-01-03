package distribution

import (
	"fmt"

	parsecmdtypes "github.com/forbole/juno/v5/cmd/parse/types"
	"github.com/forbole/juno/v5/types/config"
	"github.com/spf13/cobra"

	"github.com/forbole/bdjuno/v4/database"
	"github.com/forbole/bdjuno/v4/modules/distribution"
	modulestypes "github.com/forbole/bdjuno/v4/modules/types"
)

// communityPoolCmd returns the Cobra command allowing to refresh community pool
func communityPoolCmd(parseConfig *parsecmdtypes.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "community-pool",
		Short: "Refresh community pool",
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

			// Build distribution module
			distrModule := distribution.NewModule(sources.DistrSource, parseCtx.EncodingConfig.Codec, db)

			err = distrModule.GetLatestCommunityPool()
			if err != nil {
				return fmt.Errorf("error while updating community pool: %s", err)
			}

			return nil
		},
	}
}
