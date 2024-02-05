package pricefeed

import (
	"fmt"

	parsecmdtypes "github.com/forbole/juno/v5/cmd/parse/types"
	"github.com/forbole/juno/v5/types/config"
	"github.com/spf13/cobra"

	"github.com/forbole/callisto/v4/database"
	"github.com/forbole/callisto/v4/modules/pricefeed"
)

// priceHistoryCmd returns the Cobra command allowing to store token price history
func priceHistoryCmd(parseConfig *parsecmdtypes.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "history",
		Short: "Store token price history",
		RunE: func(cmd *cobra.Command, args []string) error {
			parseCtx, err := parsecmdtypes.GetParserContext(config.Cfg, parseConfig)
			if err != nil {
				return err
			}

			// Get the database
			db := database.Cast(parseCtx.Database)

			// Build pricefeed module
			pricefeedModule := pricefeed.NewModule(config.Cfg, parseCtx.EncodingConfig.Codec, db)

			err = pricefeedModule.RunAdditionalOperations()
			if err != nil {
				return fmt.Errorf("error while storing tokens: %s", err)
			}

			err = pricefeedModule.UpdatePricesHistory()
			if err != nil {
				return fmt.Errorf("error while updating price history: %s", err)
			}

			return nil
		},
	}
}
