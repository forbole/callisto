package periodictask

import (
	parsecmdtypes "github.com/forbole/juno/v3/cmd/parse/types"
	"github.com/forbole/juno/v3/types/config"
	"github.com/spf13/cobra"

	"github.com/forbole/bdjuno/v3/database"
	"github.com/forbole/bdjuno/v3/modules/distribution"
	modulestypes "github.com/forbole/bdjuno/v3/modules/types"
)

// distributionCmd returns the Cobra command allowing to refresh data that's obtained from x/distribution periodic tasks
func distributionCmd(parseConfig *parsecmdtypes.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "distribution",
		Short: "Run x/distribution periodic task",
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
			distrModule := distribution.NewModule(sources.DistrSource, parseCtx.EncodingConfig.Marshaler, db)

			err = distrModule.GetLatestCommunityPool()
			if err != nil {
				return err
			}

			return nil
		},
	}
}
