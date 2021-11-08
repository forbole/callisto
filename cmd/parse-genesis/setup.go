package parseGenesis

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	parsecmd "github.com/forbole/juno/v2/cmd/parse"
	"github.com/forbole/juno/v2/modules"
	modsregistrar "github.com/forbole/juno/v2/modules/registrar"

	"github.com/forbole/juno/v2/database"
	"github.com/forbole/juno/v2/types/config"
)

// GetParsingContext setups all the things that should be later passed to StartParsing in order
// to parse the chain data properly.
func GetGenesisModules(parseConfig *parsecmd.Config) ([]modules.Module, error) {
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
