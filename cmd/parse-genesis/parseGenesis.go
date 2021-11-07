package parseGenesis

import (
	"encoding/json"
	"fmt"

	"github.com/forbole/juno/v2/modules"

	"github.com/forbole/juno/v2/cmd/parse"
	"github.com/spf13/cobra"
)

// NewParseGenesisCmd returns the command to be run for parsing the genesis
func NewParseGenesisCmd(parseCfg *parse.Config) *cobra.Command {
	return &cobra.Command{
		Use:     "parse-genesis",
		Short:   "Parse the genesis file",
		PreRunE: parse.ReadConfig(parseCfg),
		RunE: func(cmd *cobra.Command, args []string) error {
			context, err := parse.GetParsingContext(parseCfg)
			if err != nil {
				return err
			}
			genesis, err := context.Node.Genesis()
			if err != nil {
				return err
			}

			var appState map[string]json.RawMessage
			if err := json.Unmarshal(genesis.Genesis.AppState, &appState); err != nil {
				return fmt.Errorf("error unmarshalling genesis state: %s", err)
			}

			for _, module := range context.Modules {
				if genesisModule, ok := module.(modules.GenesisModule); ok {
					err = genesisModule.HandleGenesis(genesis.Genesis, appState)
					if err != nil {
						return fmt.Errorf("error while handling genesis of %s module: %s", module.Name(), err)
					}
				}
			}
			return nil
		},
	}
}
