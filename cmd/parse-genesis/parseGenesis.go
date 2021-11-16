package parseGenesis

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

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
		Example: "bdjuno parse-genesis auth bank consensus gov history staking",
		PreRunE: parse.ReadConfig(parseCfg),
		PostRun: func(cmd *cobra.Command, args []string) {
			fmt.Println("Genesis file has been parsed")
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("no module name specified")
			}

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

			var invalidMods []string
			var registeredModuleName string
			for _, argModuleName := range args {
				// Traverse arguments
				for _, module := range registeredModules {
					registeredModuleName = module.Name()
					genesisModule, ok := module.(modules.GenesisModule)

					if ok && argModuleName == registeredModuleName {
						err = genesisModule.HandleGenesis(genesisDoc, genesisState)
						if err != nil {
							return fmt.Errorf("error while handling genesis of %s module: %s", registeredModuleName, err)
						}
						break
					}

				}
				if argModuleName != registeredModuleName {
					invalidMods = append(invalidMods, argModuleName)
				}
			}

			if len(invalidMods) != 0 {
				return fmt.Errorf("not registered or invalid module name(s): %s", strings.Join(invalidMods, ", "))
			}

			return nil
		},
	}
}
