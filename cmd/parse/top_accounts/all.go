package top_accounts

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/forbole/bdjuno/v3/modules/bank"
	"github.com/forbole/bdjuno/v3/modules/distribution"
	"github.com/forbole/bdjuno/v3/modules/staking"
	topaccounts "github.com/forbole/bdjuno/v3/modules/top_accounts"
	modulestypes "github.com/forbole/bdjuno/v3/modules/types"
	"github.com/forbole/bdjuno/v3/types"
	"github.com/rs/zerolog/log"

	parsecmdtypes "github.com/forbole/juno/v3/cmd/parse/types"
	"github.com/forbole/juno/v3/parser"
	"github.com/forbole/juno/v3/types/config"
	"github.com/spf13/cobra"

	"github.com/forbole/bdjuno/v3/database"
	"github.com/forbole/bdjuno/v3/modules/auth"
)

var (
	waitGroup sync.WaitGroup
)

const (
	flagWorker = "worker"
)

func allCmd(parseConfig *parsecmdtypes.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use: "all",
		RunE: func(cmd *cobra.Command, args []string) error {
			parseCtx, err := parsecmdtypes.GetParserContext(config.Cfg, parseConfig)
			if err != nil {
				return err
			}

			sources, err := modulestypes.BuildSources(config.Cfg.Node, parseCtx.EncodingConfig)
			if err != nil {
				return err
			}

			// Get the database
			db := database.Cast(parseCtx.Database)

			// Build modules
			authModule := auth.NewModule(sources.AuthSource, nil, parseCtx.EncodingConfig.Marshaler, db)
			bankModule := bank.NewModule(nil, sources.BankSource, parseCtx.EncodingConfig.Marshaler, db)
			distriModule := distribution.NewModule(sources.DistrSource, parseCtx.EncodingConfig.Marshaler, db)
			stakingModule := staking.NewModule(sources.StakingSource, parseCtx.EncodingConfig.Marshaler, db)
			topaccountsModule := topaccounts.NewModule(bankModule, distriModule, stakingModule, nil, parseCtx.EncodingConfig.Marshaler, db)

			// Get workers
			exportQueue := NewQueue(5)
			workerCount, _ := cmd.Flags().GetInt64(flagWorker)
			workers := make([]Worker, workerCount)
			for i := range workers {
				workers[i] = NewWorker(exportQueue, topaccountsModule)
			}

			waitGroup.Add(1)

			// Get all base accounts, height set to 0 for querying the latest data on chain
			accounts, err := authModule.GetAllBaseAccounts(0)
			if err != nil {
				return fmt.Errorf("error while getting base accounts: %s", err)
			}

			log.Debug().Int("total", len(accounts)).Msg("saving accounts...")
			// Store accounts
			err = db.SaveAccounts(accounts)
			if err != nil {
				return err
			}

			for i, w := range workers {
				log.Debug().Int("number", i+1).Msg("starting worker...")
				go w.start()
			}

			trapSignal(parseCtx)

			go enqueueAddresses(exportQueue, accounts)

			waitGroup.Wait()
			return nil
		},
	}

	cmd.Flags().Int64(flagWorker, 1, "worker count")

	return cmd
}

func enqueueAddresses(exportQueue AddressQueue, accounts []types.Account) {
	for _, account := range accounts {
		exportQueue <- account.Address
	}
}

// trapSignal will listen for any OS signal and invoke Done on the main
// WaitGroup allowing the main process to gracefully exit.
func trapSignal(ctx *parser.Context) {
	var sigCh = make(chan os.Signal, 1)

	signal.Notify(sigCh, syscall.SIGTERM)
	signal.Notify(sigCh, syscall.SIGINT)

	go func() {
		sig := <-sigCh
		log.Info().Str("signal", sig.String()).Msg("caught signal; shutting down...")
		defer ctx.Node.Stop()
		defer ctx.Database.Close()
		defer waitGroup.Done()
	}()
}
