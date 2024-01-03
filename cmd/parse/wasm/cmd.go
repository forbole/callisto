package wasm

import (
	parsecmdtypes "github.com/forbole/juno/v5/cmd/parse/types"
	"github.com/spf13/cobra"
)

// NewWasmCmd returns the Cobra command allowing to fix various things related to the x/wasm module
func NewWasmCmd(parseConfig *parsecmdtypes.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "wasm",
		Short: "Fix things related to the x/wasm module",
	}

	cmd.AddCommand(
		codeCmd(parseConfig),
	)

	return cmd
}
