package wasm

import (
	"fmt"

	modulestypes "github.com/forbole/bdjuno/v4/modules/types"
	"github.com/forbole/bdjuno/v4/modules/wasm"

	parsecmdtypes "github.com/forbole/juno/v5/cmd/parse/types"
	"github.com/forbole/juno/v5/types/config"
	"github.com/spf13/cobra"

	"github.com/forbole/bdjuno/v4/database"
)

// contractsCmd returns a Cobra command that allows to fix the contracts info
func contractsCmd(parseConfig *parsecmdtypes.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "contracts",
		Short: "Query all available contracts",
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

			// Build the wasm module
			wasmModule := wasm.NewModule(sources.WasmSource, parseCtx.EncodingConfig.Marshaler, db)

			// Get latest height
			height, err := parseCtx.Node.LatestHeight()
			if err != nil {
				return fmt.Errorf("error while getting latest block height: %s", err)
			}

			err = wasmModule.StoreContracts(height)
			if err != nil {
				return fmt.Errorf("error while stroing all x/wasm contracts infos: %s", err)
			}

			return nil
		},
	}
}
