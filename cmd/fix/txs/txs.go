package txs

import (
	"fmt"

	"github.com/forbole/juno/v2/cmd/parse"
	"github.com/forbole/juno/v2/types/config"
	"github.com/spf13/cobra"

	"github.com/forbole/bdjuno/v2/database"
	"github.com/forbole/bdjuno/v2/modules"
)

// slashesCmd returns a Cobra command that allows to fix the delegations for all the slashed validators.
func transactionsCmd(parseConfig *parse.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "fix",
		Short: "Fix missing transactions in database from the latest known height",
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

			// Get latest height
			height, err := parseCtx.Node.LatestHeight()
			if err != nil {
				return fmt.Errorf("error while getting latest block height: %s", err)
			}

			return nil
		},
	}
}
