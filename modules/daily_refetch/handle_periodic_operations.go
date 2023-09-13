package daily_refetch

import (
	"fmt"
	"time"

	"github.com/forbole/juno/v5/parser"
	"github.com/forbole/juno/v5/types/config"

	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"

	parsecmdtypes "github.com/forbole/juno/v5/cmd/parse/types"
)

func (m *Module) RegisterPeriodicOperations(scheduler *gocron.Scheduler) error {
	log.Debug().Str("module", "daily refetch").Msg("setting up periodic tasks")

	// Setup a cron job to run every midnight
	if _, err := scheduler.Every(1).Day().At("00:00").Do(func() {
		m.refetchMissingBlocks()
	}); err != nil {
		return fmt.Errorf("error while setting up daily refetch periodic operation: %s", err)
	}

	return nil
}

// refetchMissingBlocks checks for missing blocks from one day ago and refetches them
func (m *Module) refetchMissingBlocks() error {
	log.Trace().Str("module", "daily refetch").Str("refetching", "blocks").
		Msg("refetching missing blocks")

	latestBlock, err := m.node.LatestHeight()
	if err != nil {
		return fmt.Errorf("error while getting latest block: %s", err)
	}

	blockHeightDayAgo, err := m.database.GetBlockHeightTimeDayAgo(time.Now())
	if err != nil {
		return fmt.Errorf("error while getting block height from a day ago: %s", err)
	}
	var startHeight = blockHeightDayAgo.Height

	missingBlocks := m.database.GetMissingBlocks(startHeight, latestBlock)

	// return if no blocks are missing
	if len(missingBlocks) == 0 {
		return nil
	}

	parseCtx, err := parsecmdtypes.GetParserContext(config.Cfg, parsecmdtypes.NewConfig())
	if err != nil {
		return err
	}

	workerCtx := parser.NewContext(parseCtx.EncodingConfig, parseCtx.Node, parseCtx.Database, parseCtx.Logger, parseCtx.Modules)
	worker := parser.NewWorker(workerCtx, nil, 0)

	log.Info().Int64("start height", startHeight).Int64("end height", latestBlock).
		Msg("getting missing blocks and transactions from a day ago")
	for _, block := range missingBlocks {
		err = worker.Process(block)
		if err != nil {
			return fmt.Errorf("error while re-fetching block %d: %s", block, err)
		}
	}

	return nil

}
