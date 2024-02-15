package wasm

import (
	"fmt"

	modulestypes "github.com/forbole/callisto/v4/modules/types"
	"github.com/forbole/callisto/v4/modules/wasm"

	parsecmdtypes "github.com/forbole/juno/v5/cmd/parse/types"
	"github.com/forbole/juno/v5/types/config"
	"github.com/spf13/cobra"

	"github.com/forbole/callisto/v4/database"
)

// codeCmd returns the Cobra command allowing to fix all things related to x/wasm contract
func codeCmd(parseConfig *parsecmdtypes.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "codes",
		Short: "Parse the x/wasm codes and store them in the database",
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

			height, err := db.GetLastBlockHeight()
			if err != nil {
				return err
			}

			wasmModule := wasm.NewModule(sources.WasmSource, parseCtx.EncodingConfig.Codec, db)

			wasmCodes, err := wasmModule.GetWasmCodes(height)
			if err != nil {
				return fmt.Errorf("error while getting wasm codes: %s", err)
			}

			err = db.SaveWasmCodes(wasmCodes)
			if err != nil {
				return fmt.Errorf("error while saving wasm codes: %s", err)
			}

			return nil
		},
	}
}
