package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	sdk "github.com/cosmos/cosmos-sdk/types"
	parsecmd "github.com/forbole/juno/v2/cmd/parse"
	"github.com/forbole/juno/v2/database"
	"github.com/forbole/juno/v2/modules"
	modsregistrar "github.com/forbole/juno/v2/modules/registrar"
	"github.com/forbole/juno/v2/types/config"
	tmtypes "github.com/tendermint/tendermint/types"
)

// GetRegisteredModules returns the registered modules from the config file
func GetRegisteredModules(parseConfig *parsecmd.Config) ([]modules.Module, error) {
	// Get the global config
	cfg := config.Cfg

	// Build the codec
	encodingConfig := parseConfig.GetEncodingConfigBuilder()()

	// Setup the SDK configuration
	sdkConfig := sdk.GetConfig()
	parseConfig.GetSetupConfig()(cfg, sdkConfig)
	sdkConfig.Seal()

	// Get the db
	databaseCtx := database.NewContext(cfg.Database, &encodingConfig, parseConfig.GetLogger())
	db, err := parseConfig.GetDBBuilder()(databaseCtx)
	if err != nil {
		return nil, err
	}

	// Get the modules
	context := modsregistrar.NewContext(cfg, sdkConfig, &encodingConfig, db, nil, parseConfig.GetLogger())
	mods := parseConfig.GetRegistrar().BuildModules(context)
	registeredModules := modsregistrar.GetModules(mods, cfg.Chain.Modules, parseConfig.GetLogger())

	return registeredModules, nil
}

// GetGenesisDocAndState generates and returns the genesis doc and genesis state with genesis.json file
func GetGenesisDocAndState() (*tmtypes.GenesisDoc, map[string]json.RawMessage, error) {
	var genesisState map[string]json.RawMessage

	genesisFile, err := ioutil.ReadFile(config.GetGenesisFilePath())
	if err != nil {
		return &tmtypes.GenesisDoc{}, genesisState, fmt.Errorf("error while reading genesis file: %s", err)
	}

	genesisDoc, err := tmtypes.GenesisDocFromJSON(genesisFile)
	if err != nil {
		return &tmtypes.GenesisDoc{}, genesisState, fmt.Errorf("error while generating genesis doc from genesis.json: %s", err)
	}

	err = json.Unmarshal(genesisDoc.AppState, &genesisState)
	if err != nil {
		return &tmtypes.GenesisDoc{}, genesisState, fmt.Errorf("error while unmarshalling genesis state: %s", err)
	}

	return genesisDoc, genesisState, nil
}

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
