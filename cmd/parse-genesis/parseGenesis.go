package parseGenesis

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/forbole/juno/v2/cmd/parse"
	"github.com/forbole/juno/v2/modules"
	"github.com/forbole/juno/v2/types/config"
	"github.com/spf13/cobra"
	tmtypes "github.com/tendermint/tendermint/types"
)

var (
	HomePath = ""
)

// NewParseGenesisCmd returns the command to be run for parsing the genesis
func NewParseGenesisCmd(parseCfg *parse.Config) *cobra.Command {
	return &cobra.Command{
		Use:     "parse-genesis [module names]",
		Short:   "Parse the genesis file",
		PreRunE: parse.ReadConfig(parseCfg),
		RunE: func(cmd *cobra.Command, args []string) error {

			genesisFile, err := ioutil.ReadFile(config.GetGenesisFilePath())
			if err != nil {
				return fmt.Errorf("error while reading genesis file: %s", err)
			}

			genesisDoc, _ := tmtypes.GenesisDocFromJSON(genesisFile)

			var genesisState map[string]json.RawMessage
			err = json.Unmarshal(genesisDoc.AppState, &genesisState)
			if err != nil {
				return fmt.Errorf("error while unmarshalling genesis state: %s", err)
			}

			registeredModules, err := GetRegisteredModules(parseCfg)
			if err != nil {
				return fmt.Errorf("error while getting genesis registered modules: %s", err)
			}

			argsModules := make(map[string]bool)
			for _, arg := range args {
				argsModules[arg] = true
			}

			for _, module := range registeredModules {
				genesisModule, ok := module.(modules.GenesisModule)

				if ok && argsModules[module.Name()] {
					err = genesisModule.HandleGenesis(genesisDoc, genesisState)
					if err != nil {
						return fmt.Errorf("error while handling genesis of %s module: %s", module.Name(), err)
					}
				}
			}
			return nil
		},
	}
}
