package blocks

import (
	"fmt"

	"github.com/forbole/juno/v2/cmd/parse"
	"github.com/spf13/cobra"

	"github.com/forbole/bdjuno/v2/database"
	"github.com/forbole/bdjuno/v2/modules/consensus"
	"github.com/forbole/bdjuno/v2/utils"
)

// blocksCmd returns a Cobra command that allows to fix missing blocks in database
func blocksCmd(parseConfig *parse.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "refetch",
		Short: "Fix missing blocks in database from the latest known height",
		RunE: func(cmd *cobra.Command, args []string) error {
			parseCtx, err := parse.GetParsingContext(parseConfig)
			if err != nil {
				return err
			}

			// Get the database
			db := database.Cast(parseCtx.Database)

			// Build the consensus module
			consensusModule := consensus.NewModule(db)

			// // Get latest height
			// height, err := parseCtx.Node.LatestHeight()
			// if err != nil {
			// 	return fmt.Errorf("error while getting chain latest block height: %s", err)
			// }

			var k int64 = 1
			for ; k <= 7607; k++ {
				err = refreshBlock(parseCtx, k, consensusModule)
				if err != nil {
					return err
				}
			}

			return nil
		},
	}
}

func refreshBlock(parseCtx *parse.Context, blockHeight int64, consensusModule *consensus.Module) error {
	// Get the block details
	blockDetails, err := utils.QueryBlock(parseCtx.Node, blockHeight)
	if err != nil {
		return err
	}

	err = consensusModule.UpdateBlock(blockDetails)

	if err != nil {
		return fmt.Errorf("error while updating block %v: %s", blockHeight, err)
	}

	return nil
}
