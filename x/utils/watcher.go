package utils

// DONTCOVER

import (
	"io/ioutil"
	"path"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/desmos-labs/juno/config"
	"github.com/desmos-labs/juno/db"
	"github.com/desmos-labs/juno/parse"
	"github.com/desmos-labs/juno/parse/client"
	"github.com/forbole/bdjuno/database"
	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"
)

// WatchMethod allows to watch for a method that returns an error.
// It executes the given method in a goroutine, logging any error that might raise.
func WatchMethod(method func() error) {
	go func() {
		err := method()
		if err != nil {
			log.Error().Err(err).Send()
		}
	}()
}

// WatchModules check goroutine and get the running enabled modules every 30 second
// The list of module depends on bdjuno/x/... so that user only need to update schema of database
// when adding a new module
func WatchModules(scheduler *gocron.Scheduler) parse.AdditionalOperation {

	return func(_ config.Config, _ *codec.Codec, _ client.ClientProxy, db db.Database) error {
		bdDatabase, ok := db.(database.BigDipperDb)
		if !ok {
			log.Fatal().Str("module", "util").Msg("given database instance is not a BigDipperDb")
		}

		if _, err := scheduler.Every(30).Second().StartImmediately().Do(func() {
			WatchMethod(func() error { return watchModules(bdDatabase) })
		}); err != nil {
			return err
		}

		return nil
	}
}

func watchModules(bdDatabase database.BigDipperDb) error {
	_, b, _, _ := runtime.Caller(0)
	root := filepath.Dir(path.Join(path.Dir(b)))
	files, err := ioutil.ReadDir(root)
	if err != nil {
		log.Fatal()
	}
	modules := make(map[string]bool)
	var profilingBuffer strings.Builder
	pprof.Lookup("goroutine").WriteTo(&profilingBuffer, 1)
	s := profilingBuffer.String()
	for _, name := range files {
		modules[name.Name()] = strings.Contains(s, name.Name())
	}

	_, ok := modules[".DS_Store"]
	if ok {
		delete(modules, ".DS_Store") // delete system directory
	}

	_, ok = modules["utils"]
	if ok {
		delete(modules, "utils") // delete system directory
	}

	return bdDatabase.InsertEnableModules(modules)
}
