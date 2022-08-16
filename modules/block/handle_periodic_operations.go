package block

import (
	"fmt"

	"github.com/forbole/juno/v3/parser"

	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"

	parsecmdtypes "github.com/forbole/juno/v3/cmd/parse/types"
	"github.com/forbole/juno/v3/types/config"
)

func (m *Module) RegisterPeriodicOperations(scheduler *gocron.Scheduler) error {
	log.Debug().Str("module", "block refetch").Msg("setting up periodic tasks")

	if _, err := scheduler.Every(2).Minutes().Do(func() {
		m.checkMissingBlocks()
	}); err != nil {
		return fmt.Errorf("error while setting up block refetch periodic operation: %s", err)
	}

	return nil
}

// checkMissingBlocks checks for any missing blocks and refetches it
func (m *Module) checkMissingBlocks() error {
	log.Trace().Str("module", "blocks").Str("refetching", "blocks").
		Msg("refetching missing blocks")

	latestBlock, err := m.node.LatestHeight()
	if err != nil {
		return fmt.Errorf("error while getting latest block: %s", err)
	}

	blockCount, err := m.database.GetTotalBlocks()
	if err != nil {
		return fmt.Errorf("error while getting last block: %s", err)
	}

	// averageBlockPerDay, err := m.db.GetAverageBlockPerDay()
	// if err != nil {
	// 	return fmt.Errorf("error while getting average day block time: %s", err)
	// }

	var startHeight int64 = 1

	var endHeight int64 = 123911

	if blockCount != latestBlock {
		parseCtx, err := parsecmdtypes.GetParserContext(config.Cfg, m.config)
		if err != nil {
			return err
		}

		workerCtx := parser.NewContext(parseCtx.EncodingConfig, parseCtx.Node, parseCtx.Database, parseCtx.Logger, parseCtx.Modules)
		worker := parser.NewWorker(workerCtx, nil, 0)

		log.Info().Int64("start height", startHeight).Int64("end height", endHeight).
			Msg("getting missing blocks and transactions")
		for k := startHeight; k <= endHeight; k++ {
			err = worker.ProcessIfNotExists(k)
			if err != nil {
				return fmt.Errorf("error while re-fetching block %d: %s", k, err)
			}
		}

	}

	return nil

}
