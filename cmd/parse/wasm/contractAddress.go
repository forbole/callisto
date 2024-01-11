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

// contractAddressCmd returns a Cobra command that allows to fix the contract info
// with given contract address
func contractAddressCmd(parseConfig *parsecmdtypes.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "contract-address [address]",
		Short: "Query available contract details woth the given contract address",
		RunE: func(cmd *cobra.Command, args []string) error {
			contractAddress := args[0]

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
			wasmModule := wasm.NewModule(sources.WasmSource, parseCtx.EncodingConfig.Codec, db)

			// Get latest height
			height, err := parseCtx.Node.LatestHeight()
			if err != nil {
				return fmt.Errorf("error while getting latest block height: %s", err)
			}

			err = wasmModule.ParseContractDetails(contractAddress, height-50)
			if err != nil {
				return fmt.Errorf("error while stroing x/wasm contract info for the contract address %s: %s", contractAddress, err)
			}

			return nil
		},
	}
}
