package v5

import (
	"fmt"

	v5db "github.com/forbole/callisto/v4/database/migrate/v5"
	parse "github.com/forbole/juno/v5/cmd/parse/types"
	"github.com/forbole/juno/v5/database"
	"github.com/forbole/juno/v5/database/postgresql"
	"github.com/forbole/juno/v5/types/config"
)

// RunMigration runs the migrations to v5
func RunMigration(parseConfig *parse.Config) error {
	cfg, err := GetConfig()
	if err != nil {
		return fmt.Errorf("error while reading config: %s", err)
	}

	// Migrate the database
	err = migrateDb(cfg, parseConfig)
	if err != nil {
		return fmt.Errorf("error while migrating database: %s", err)
	}

	return nil
}

func migrateDb(cfg config.Config, parseConfig *parse.Config) error {
	// Build the codec
	encodingConfig := parseConfig.GetEncodingConfigBuilder()()

	// Get the db
	databaseCtx := database.NewContext(cfg.Database, encodingConfig, parseConfig.GetLogger())
	db, err := postgresql.Builder(databaseCtx)
	if err != nil {
		return fmt.Errorf("error while building the db: %s", err)
	}

	// Build the migrator and perform the migrations
	migrator := v5db.NewMigrator(db.(*postgresql.Database))
	return migrator.Migrate()
}
