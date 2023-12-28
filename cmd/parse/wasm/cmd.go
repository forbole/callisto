package wasm

import (
	parsecmdtypes "github.com/forbole/juno/v5/cmd/parse/types"
	"github.com/spf13/cobra"
)

// NewWasmCmd returns the Cobra command that allows to fix all the things related to the x/wasm module
func NewWasmCmd(parseConfig *parsecmdtypes.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "wasm",
		Short: "Fix things related to the x/wasm module",
	}

	cmd.AddCommand(
		contractsCmd(parseConfig),
	)

	return cmd
}
