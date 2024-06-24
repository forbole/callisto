package gov

import (
	"github.com/spf13/cobra"

	parsecmdtypes "github.com/forbole/juno/v6/cmd/parse/types"
	"github.com/forbole/juno/v6/types/config"

	"github.com/forbole/callisto/v4/database"
	"github.com/forbole/callisto/v4/modules/distribution"
	"github.com/forbole/callisto/v4/modules/gov"
	"github.com/forbole/callisto/v4/modules/mint"
	"github.com/forbole/callisto/v4/modules/slashing"
	"github.com/forbole/callisto/v4/modules/staking"
	modulestypes "github.com/forbole/callisto/v4/modules/types"
	"github.com/forbole/callisto/v4/utils"
)

func paramsCmd(parseConfig *parsecmdtypes.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "params",
		Short: "Get the current parameters of the gov module",
		RunE: func(cmd *cobra.Command, args []string) error {
			parseCtx, err := parsecmdtypes.GetParserContext(config.Cfg, parseConfig)
			if err != nil {
				return err
			}

			cdc := utils.GetCodec()
			sources, err := modulestypes.BuildSources(config.Cfg.Node, cdc)
			if err != nil {
				return err
			}

			// Get the database
			db := database.Cast(parseCtx.Database)

			// Build expected modules of gov modules
			distrModule := distribution.NewModule(sources.DistrSource, cdc, db)
			mintModule := mint.NewModule(sources.MintSource, cdc, db)
			slashingModule := slashing.NewModule(sources.SlashingSource, cdc, db)
			stakingModule := staking.NewModule(sources.StakingSource, cdc, db)

			// Build the gov module
			govModule := gov.NewModule(sources.GovSource, distrModule, mintModule, slashingModule, stakingModule, cdc, db)

			height, err := parseCtx.Node.LatestHeight()
			if err != nil {
				return err
			}

			return govModule.UpdateParams(height)
		},
	}
}
