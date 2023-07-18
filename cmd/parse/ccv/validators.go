package ccv

import (
	"fmt"

	"github.com/forbole/bdjuno/v4/database"
	consumer "github.com/forbole/bdjuno/v4/modules/ccv/consumer"
	parsecmdtypes "github.com/forbole/juno/v5/cmd/parse/types"
	"github.com/forbole/juno/v5/types/config"
	"github.com/spf13/cobra"
)

// ccvValidatorsCmd returns a Cobra command that allows to refresh ccv validators info in database
func ccvValidatorsCmd(parseConfig *parsecmdtypes.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "validators",
		Short: "Fix the ccv validators details stored in database",
		RunE: func(cmd *cobra.Command, args []string) error {
			parseCtx, err := parsecmdtypes.GetParserContext(config.Cfg, parseConfig)
			if err != nil {
				return err
			}

			// Get the database
			db := database.Cast(parseCtx.Database)

			// Build the ccv consumer module
			ccvConsumerModule := consumer.NewModule(parseCtx.EncodingConfig.Codec, db)

			// Get latest height
			height, err := parseCtx.Node.LatestHeight()
			if err != nil {
				return fmt.Errorf("error while getting latest block height: %s", err)
			}

			err = ccvConsumerModule.UpdateCcvValidators(height)
			if err != nil {
				return fmt.Errorf("error while updating all ccv validators at height: %s", err)
			}

			return nil
		},
	}
}
