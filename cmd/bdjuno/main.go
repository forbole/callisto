package main

import (
	desmosapp "github.com/desmos-labs/desmos/app"
	"github.com/desmos-labs/juno/cmd"
	initcmd "github.com/desmos-labs/juno/cmd/init"
	parsecmd "github.com/desmos-labs/juno/cmd/parse"

	"github.com/forbole/bdjuno/types/config"

	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/modules"
)

func main() {
	// Setup the config
	initCfg := initcmd.NewConfig().
		WithConfigFlagSetup(config.SetupConfigFlags).
		WithConfigCreator(config.CreateConfig)

	parseCfg := parsecmd.NewConfig().
		WithConfigParser(config.ParseConfig).
		WithRegistrar(modules.NewRegistrar()).
		WithDBBuilder(database.Builder).
		WithEncodingConfigBuilder(desmosapp.MakeTestEncodingConfig)

	cfg := cmd.NewConfig("bdjuno").
		WithInitConfig(initCfg).
		WithParseConfig(parseCfg)

	// Run the command
	executor := cmd.BuildDefaultExecutor(cfg)
	err := executor.Execute()
	if err != nil {
		panic(err)
	}
}
