package parseGenesis

import (
	"fmt"
	"strings"

	"github.com/forbole/bdjuno/v2/cmd/parse-genesis/utils"
	"github.com/forbole/juno/v2/cmd/parse"
	"github.com/spf13/cobra"
)

// NewFixCmd returns the Cobra command allowing to fix some BDJuno bugs without having to re-sync the whole database
func NewParseGenesisCmd(parseCfg *parse.Config) *cobra.Command {
	return &cobra.Command{
		Use:     "parse-genesis [optional: module names]",
		Short:   "Parse genesis file",
		Long:    "Parse genesis file, input desired module names as arguments to parse specific modules",
		Example: "bdjuno parse-genesis auth bank consensus gov history staking",
		PreRunE: parse.ReadConfig(parseCfg),
		RunE: func(cmd *cobra.Command, args []string) error {
			genesisDoc, genesisState, err := utils.GetGenesisDocAndState(parseCfg)
			if err != nil {
				return fmt.Errorf("error while getting genesis doc or state: %s", err)
			}

			registeredModules, err := utils.GetRegisteredModules(parseCfg)
			if err != nil {
				return fmt.Errorf("error while getting genesis registered modules: %s", err)
			}

			invalidMods, err := utils.ParseGenesis(registeredModules, genesisDoc, genesisState, args)
			if err != nil {
				return fmt.Errorf("error while parsing genesis: %s", err)
			}

			if len(invalidMods) != 0 {
				// Print out inlalid / unregistered module names
				return fmt.Errorf("not registered or invalid module name(s): %s", strings.Join(invalidMods, ", "))
			}

			return nil
		},
	}
}
