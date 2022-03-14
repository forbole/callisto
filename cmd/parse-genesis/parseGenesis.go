package parsegenesis

import (
	"fmt"

	"github.com/forbole/juno/v3/cmd/parse"
	"github.com/forbole/juno/v3/modules"
	"github.com/forbole/juno/v3/types/config"
	junoutils "github.com/forbole/juno/v3/types/utils"
	"github.com/spf13/cobra"
)

// NewParseGenesisCmd returns the Cobra command allowing to parse the genesis file
func NewParseGenesisCmd(parseCfg *parse.Config) *cobra.Command {
	return &cobra.Command{
		Use:     "parse-genesis [[module names]]",
		Short:   "Parse genesis file. To parse specific modules, input module names as arguments",
		Example: "bdjuno parse-genesis auth bank consensus gov staking",
		PreRunE: parse.ReadConfig(parseCfg),
		RunE: func(cmd *cobra.Command, args []string) error {
			parseCtx, err := parse.GetParsingContext(parseCfg)
			if err != nil {
				return err
			}

			// Get the modules to parse
			var modulesToParse []modules.Module
			for _, moduleName := range args {
				module, found := getModule(moduleName, parseCtx)
				if !found {
					return fmt.Errorf("module %s is not registered", moduleName)
				}

				modulesToParse = append(modulesToParse, module)
			}

			// Default to all the modules
			if len(modulesToParse) == 0 {
				modulesToParse = parseCtx.Modules
			}

			// Get the genesis doc and state
			genesisDoc, genesisState, err := junoutils.GetGenesisDocAndState(config.Cfg.Parser.GenesisFilePath, parseCtx.Node)
			if err != nil {
				return fmt.Errorf("error while getting genesis doc and state: %s", err)
			}

			// For each module, parse the genesis
			for _, module := range modulesToParse {
				if genesisModule, ok := module.(modules.GenesisModule); ok {
					err = genesisModule.HandleGenesis(genesisDoc, genesisState)
					if err != nil {
						return fmt.Errorf("error while parsing genesis of %s module: %s", module.Name(), err)
					}
				}
			}

			return nil
		},
	}
}

// doesModuleExist tells whether a module with the given name exist inside the specified context ot not
func getModule(module string, parseCtx *parse.Context) (modules.Module, bool) {
	for _, mod := range parseCtx.Modules {
		if module == mod.Name() {
			return mod, true
		}
	}
	return nil, false
}
