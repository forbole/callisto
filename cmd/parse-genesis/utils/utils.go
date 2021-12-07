package utils

import (
	"encoding/json"
	"fmt"

	"github.com/forbole/juno/v2/modules"
	tmtypes "github.com/tendermint/tendermint/types"
)

// ParseGenesis parses the module that implements HandleGenesis method, and parses only certain modules if specified with arguments
func ParseGenesis(
	registeredMods []modules.Module, genesisDoc *tmtypes.GenesisDoc,
	genesisState map[string]json.RawMessage, arguments []string,
) ([]string, error) {
	inputArgsLen := len(arguments)

	for _, module := range registeredMods {
		genesisModule, implemented := module.(modules.GenesisModule)
		toParse := false

		for i, argModuleName := range arguments {
			// Find the registered module name that matches any provided argument, and parse the module
			if module.Name() == argModuleName {
				toParse = true
				// Remove argument from the list if found
				arguments[i] = arguments[len(arguments)-1]
				arguments = arguments[:len(arguments)-1]
			}
		}

		if inputArgsLen == 0 {
			// If no module was specified in the argument, parse all genesis modules
			toParse = true
		}

		if implemented && toParse {
			// Parse the genesis module if argument module name matches registered module name
			fmt.Printf("Parsing genesis: %s module \n", module.Name())
			err := genesisModule.HandleGenesis(genesisDoc, genesisState)
			if err != nil {
				return []string{}, err
			}
		}
	}

	// Return the rest of arguments (invalid modules)
	return arguments, nil
}
