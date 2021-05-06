package main

import (
	desmosapp "github.com/desmos-labs/desmos/app"
	"github.com/desmos-labs/juno/cmd"
	"github.com/desmos-labs/juno/cmd/parse"
	"github.com/forbole/bdjuno/types"

	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/modules"
)

func main() {
	// Setup the config
	config := parse.NewConfig("bdjuno").
		WithConfigParser(types.ParseConfig).
		WithRegistrar(modules.NewRegistrar()).
		WithDBBuilder(database.Builder).
		WithEncodingConfigBuilder(desmosapp.MakeTestEncodingConfig)

	// Run the command
	executor := cmd.BuildDefaultExecutor(config)
	err := executor.Execute()
	if err != nil {
		panic(err)
	}
}
