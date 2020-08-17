package utils

// DONTCOVER

import (
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/desmos-labs/juno/config"
	"github.com/desmos-labs/juno/db"
	"github.com/desmos-labs/juno/parse"
	"github.com/desmos-labs/juno/parse/client"
	"github.com/forbole/bdjuno/database"
	"github.com/rs/zerolog/log"
)

// WatchModules is a one shot operation that save the operating modules into the database
func WatchModules() parse.AdditionalOperation {
	return func(_ config.Config, _ *codec.Codec, _ client.ClientProxy, db db.Database) error {
		bdDatabase, ok := db.(database.BigDipperDb)
		if !ok {
			log.Fatal().Str("module", "util").Msg("given database instance is not a BigDipperDb")
		}

		WatchMethod(func() error { return watchModules(bdDatabase) })

		return nil
	}
}

func watchModules(bdDatabase database.BigDipperDb) error {
	modules := make(map[string]bool)
	modules["staking"] = true
	modules["auth"] = true
	modules["supply"] = true
	modules["distribution"] = true
	modules["pricefeed"] = true
	modules["bank"] = true
	modules["consensus"] = true
	modules["mint"] = true
	now, err := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	if err != nil {
		return err
	}
	return bdDatabase.InsertEnableModules(modules, now)
}
