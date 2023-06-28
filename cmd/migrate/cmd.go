package migrate

import (
	"fmt"
	"os"

	parsecmdtypes "github.com/forbole/juno/v5/cmd/parse/types"
	"github.com/spf13/cobra"

	v3 "github.com/forbole/bdjuno/v4/cmd/migrate/v3"
)

type Migrator func(parseCfg *parsecmdtypes.Config) error

var (
	migrations = map[string]Migrator{
		"v3": v3.RunMigration,
	}
)

func getVersions() []string {
	var versions []string
	for key := range migrations {
		versions = append(versions, key)
	}
	return versions
}

// NewMigrateCmd returns the Cobra command allowing to migrate config and tables to v3 version
func NewMigrateCmd(appName string, parseConfig *parsecmdtypes.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "migrate [to-version]",
		Short: "Perform the migrations from the current version to the specified one",
		Long: `Migrates all the necessary things (config file, database, etc) from the current version to the new one.
Note that migrations must be performed in order: to migrate from vX to vX+2 you need to do vX -> vX+1 and then vX+1 -> vX+2. 
`,
		Example: fmt.Sprintf("%s migrate v3", appName),
		Args:    cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SetOut(os.Stdout)
			if len(args) == 0 {
				cmd.Println("Please specify a version to migrate to. Available versions:")
				for _, version := range getVersions() {
					cmd.Println("-", version)
				}
				return nil
			}

			version := args[0]
			migrator, ok := migrations[version]
			if !ok {
				return fmt.Errorf("migration for version %s not found", version)
			}

			return migrator(parseConfig)
		},
	}
}
