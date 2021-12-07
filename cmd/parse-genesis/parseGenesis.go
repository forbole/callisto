package parsegenesis

import (
	"fmt"

	"github.com/forbole/juno/v2/cmd/parse"
	"github.com/forbole/juno/v2/modules"
	"github.com/forbole/juno/v2/types/config"
	junoutils "github.com/forbole/juno/v2/types/utils"
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
			parseCtx, err := parse.GetParsingContext(parseCfg)
			if err != nil {
				return err
			}

			genesisDoc, genesisState, err := junoutils.GetGenesisDocAndState(config.Cfg.Parser.GenesisFilePath, parseCtx.Node)
			if err != nil {
				return fmt.Errorf("error while getting genesis doc or state: %s", err)
			}

			var modulesToParse []modules.Module
			for _, argsModule := range args {
				var found bool
				for _, module := range parseCtx.Modules {
					if argsModule == module.Name() {
						found = true
						modulesToParse = append(modulesToParse, module)
					}
				}
				if !found {
					return fmt.Errorf("module is not registered: %s", argsModule)
				}
			}

			if len(args) == 0 {
				modulesToParse = parseCtx.Modules
			}

			for _, module := range modulesToParse {
				if genesisModule, ok := module.(modules.GenesisModule); ok {
					err := genesisModule.HandleGenesis(genesisDoc, genesisState)
					if err != nil {
						return fmt.Errorf("error while parsing genesis of %s module: %s", module.Name(), err)
					}
				}
			}

			return nil
		},
	}
}
