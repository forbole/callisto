package periodictask

import (
	"fmt"

	parsecmdtypes "github.com/forbole/juno/v3/cmd/parse/types"
	"github.com/forbole/juno/v3/types/config"
	"github.com/spf13/cobra"

	"github.com/forbole/bdjuno/v3/database"
	"github.com/forbole/bdjuno/v3/modules/pricefeed"
)

// pricefeedCmd returns the Cobra command allowing to refresh data that's obtained from pricefeed periodic tasks
func pricefeedCmd(parseConfig *parsecmdtypes.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "pricefeed",
		Short: "Run pricefeed periodic task",
		RunE: func(cmd *cobra.Command, args []string) error {
			parseCtx, err := parsecmdtypes.GetParserContext(config.Cfg, parseConfig)
			if err != nil {
				return err
			}

			// Get the database
			db := database.Cast(parseCtx.Database)

			// Build pricefeed module
			pricefeedModule := pricefeed.NewModule(config.Cfg, parseCtx.EncodingConfig.Marshaler, db)

			err = pricefeedModule.RunAdditionalOperations()
			if err != nil {
				return fmt.Errorf("error while storing tokens: %s", err)
			}

			err = pricefeedModule.UpdatePrice()
			if err != nil {
				return fmt.Errorf("error while updating price: %s", err)
			}

			err = pricefeedModule.UpdatePricesHistory()
			if err != nil {
				return fmt.Errorf("error while updating price history: %s", err)
			}

			return nil
		},
	}
}
