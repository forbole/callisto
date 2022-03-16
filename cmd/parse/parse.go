package parse

import (
	parse "github.com/forbole/juno/v3/cmd/parse/types"
	"github.com/spf13/cobra"

	fixblocks "github.com/forbole/juno/v3/cmd/parse/blocks"

	fixauth "github.com/forbole/bdjuno/v2/cmd/parse/auth"
	fixgov "github.com/forbole/bdjuno/v2/cmd/parse/gov"
	fixstaking "github.com/forbole/bdjuno/v2/cmd/parse/staking"
	parsegenesis "github.com/forbole/juno/v3/cmd/parse/genesis"
)

// NewParseCmd returns the Cobra command allowing to parse some chain data without having to re-sync the whole database
func NewParseCmd(parseCfg *parse.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "parse",
		Short:             "Parse some data without the need to re-syncing the whole database from scratch",
		PersistentPreRunE: runPersistentPreRuns(parse.ReadConfigPreRunE(parseCfg)),
	}

	cmd.AddCommand(
		fixauth.NewAuthCmd(parseCfg),
		fixblocks.NewBlocksCmd(parseCfg),
		fixgov.NewGovCmd(parseCfg),
		fixstaking.NewStakingCmd(parseCfg),
		parsegenesis.NewGenesisCmd(parseCfg),
	)

	return cmd
}

func runPersistentPreRuns(preRun func(_ *cobra.Command, _ []string) error) func(_ *cobra.Command, _ []string) error {
	return func(cmd *cobra.Command, args []string) error {
		if root := cmd.Root(); root != nil {
			if root.PersistentPreRunE != nil {
				err := root.PersistentPreRunE(root, args)
				if err != nil {
					return err
				}
			}
		}

		return preRun(cmd, args)
	}
}
