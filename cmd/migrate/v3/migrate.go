package v3

import (
	"fmt"
	"os"

	"github.com/forbole/callisto/v4/modules/actions"

	parsecmdtypes "github.com/forbole/juno/v5/cmd/parse/types"

	"gopkg.in/yaml.v3"

	junov4 "github.com/forbole/juno/v5/cmd/migrate/v4"
	"github.com/forbole/juno/v5/types/config"
)

// RunMigration runs the migrations from v2 to v3
func RunMigration(parseConfig *parsecmdtypes.Config) error {
	// Run Juno migration
	err := junov4.RunMigration(parseConfig)
	if err != nil {
		return err
	}

	// Migrate the config
	cfg, err := migrateConfig()
	if err != nil {
		return fmt.Errorf("error while migrating config: %s", err)
	}

	// Refresh the global configuration
	err = parsecmdtypes.UpdatedGlobalCfg(parseConfig)
	if err != nil {
		return err
	}

	bz, err := yaml.Marshal(&cfg)
	if err != nil {
		return fmt.Errorf("error while serializing config: %s", err)
	}

	err = os.WriteFile(config.GetConfigFilePath(), bz, 0600)
	if err != nil {
		return fmt.Errorf("error while writing v3 config: %s", err)
	}

	return nil
}

func migrateConfig() (Config, error) {
	cfg, err := GetConfig()
	if err != nil {
		return Config{}, fmt.Errorf("error while reading v2 config: %s", err)
	}

	// Enable the actions module if not enabled
	if !cfg.Chain.IsModuleEnabled(actions.ModuleName) {
		cfg.Chain.Modules = append(cfg.Chain.Modules, actions.ModuleName)
	}

	if cfg.Actions == nil {
		cfg.Actions = actions.NewConfig("127.0.0.1", 3000, nil)
	}

	return cfg, nil
}
