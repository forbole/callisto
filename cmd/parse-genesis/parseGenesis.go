package parsegenesis

import (
	"fmt"
	"strings"

	"github.com/forbole/bdjuno/v2/cmd/parse-genesis/utils"
	"github.com/forbole/juno/v2/cmd/parse"
	juno "github.com/forbole/juno/v2/types/utils"
	"github.com/spf13/cobra"
)

// NewParseGenesisCmd returns the Cobra command allowing to parse the genesis file
func NewParseGenesisCmd(parseCfg *parse.Config) *cobra.Command {
	return &cobra.Command{
		Use:     "parse-genesis [optional: module names]",
		Short:   "Parse genesis file. To parse specific modules, input module names as arguments",
		Example: "bdjuno parse-genesis auth bank consensus gov history staking",
		PreRunE: parse.ReadConfig(parseCfg),
		RunE: func(cmd *cobra.Command, args []string) error {
			genesisDoc, genesisState, err := juno.GetGenesisDocAndState()
			if err != nil {
				return fmt.Errorf("error while getting genesis doc or state: %s", err)
			}

			parseCtx, err := parse.GetParsingContext(parseCfg)
			if err != nil {
				return err
			}

			invalidInputs, err := utils.ParseGenesis(parseCtx.Modules, genesisDoc, genesisState, args)
			if err != nil {
				return fmt.Errorf("error while parsing genesis: %s", err)
			}

			if len(invalidInputs) != 0 {
				// Print out invalid / unregistered module names
				return fmt.Errorf("not registered or invalid module name(s): %s", strings.Join(invalidInputs, ", "))
			}

			return nil
		},
	}
}
