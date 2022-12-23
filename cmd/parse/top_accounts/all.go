package top_accounts

import (
	"fmt"

	modulestypes "github.com/forbole/bdjuno/v3/modules/types"

	parsecmdtypes "github.com/forbole/juno/v3/cmd/parse/types"
	"github.com/forbole/juno/v3/types/config"
	"github.com/spf13/cobra"

	"github.com/forbole/bdjuno/v3/database"
	"github.com/forbole/bdjuno/v3/modules/auth"
)

func allCmd(parseConfig *parsecmdtypes.Config) *cobra.Command {
	return &cobra.Command{
		Use: "all",
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

			// Build modules
			// stakingModule := staking.NewModule(sources.StakingSource, nil, parseCtx.EncodingConfig.Marshaler, db)
			authModule := auth.NewModule(sources.AuthSource, nil, parseCtx.EncodingConfig.Marshaler, db)

			accounts, err := authModule.GetAllBaseAccounts(0)
			if err != nil {
				return fmt.Errorf("error while getting account base accounts: %s", err)
			}

			err = db.SaveAccounts(accounts)
			if err != nil {
				return err
			}

			return nil
		},
	}
}
