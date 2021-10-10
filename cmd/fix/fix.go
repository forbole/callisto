package fix

import (
	"github.com/desmos-labs/juno/v2/cmd/parse"
	"github.com/spf13/cobra"

	"github.com/forbole/bdjuno/v2/cmd/fix/slashes"
)

// FixCmd returns the Cobra command allowing to fix some BDJuno bugs without having to re-sync the whole database
// nolint: golint
func FixCmd(parseCfg *parse.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "fix",
		Short:             "Apply some fixes without the need to re-syncing the whole database from scratch",
		PersistentPreRunE: parse.ReadConfig(parseCfg),
	}

	cmd.AddCommand(
		slashes.FixSlashesCmd(parseCfg),
	)

	return cmd
}
