package vipcoin

import (
	"errors"
	"fmt"
	"os"
	"time"

	"git.ooo.ua/vipcoin/lib/errs"
	"git.ooo.ua/vipcoin/lib/filter"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/forbole/juno/v2/modules"
	"github.com/forbole/juno/v2/types"

	dbtypes "github.com/forbole/bdjuno/v2/database/types"
)

func (m *module) scheduler() {
	for {
		lastBlock, err := m.lastBlockRepo.Get()
		if err != nil {
			m.logger.Error("Fail lastBlockRepo.Get", "module", "overgold", "error", err)
			continue
		}

		lastBlock++

		if err := m.parseBlock(lastBlock); err != nil {
			time.Sleep(time.Second)

			if errors.As(err, &errs.NotFound{}) {
				continue
			}

			m.logger.Error("Fail parseBlock", "module", "overgold", "error", err)
			continue
		}

		if err = m.lastBlockRepo.Update(lastBlock); err != nil {
			m.logger.Error("Fail lastBlockRepo.Update", "module", "overgold", "error", err)
			os.Exit(1)
		}
	}
}

func (m *module) parseBlock(lastBlock uint64) error {
	block, err := m.db.GetBlock(filter.NewFilter().SetArgument(dbtypes.FieldHeight, lastBlock))
	if err != nil {
		if errors.As(err, &errs.NotFound{}) {
			return err
		}

		return errs.Internal{Cause: err.Error()}
	}

	m.logger.Debug("parse block", "height", block.Height)

	if block.TxNum == 0 {
		return nil
	}

	return m.parseTx(block)
}

func (m *module) parseTx(block dbtypes.BlockRow) error {
	txs, err := m.db.GetTransactions(
		filter.NewFilter().
			SetCondition(filter.ConditionAND).
			SetArgument(dbtypes.FieldHeight, block.Height),
	)
	if err != nil {
		if errors.As(err, &errs.NotFound{}) {
			return err
		}

		return errs.Internal{Cause: err.Error()}
	}

	if block.TxNum != int64(len(txs)) {
		return &errs.Conflict{
			Cause: fmt.Errorf("mismatch txs in block: height - %d, expected tx num - %d, exist - %d ",
				block.Height,
				block.TxNum,
				int64(len(txs))).Error()}
	}

	for _, tx := range txs {
		if !tx.Successful() {
			continue
		}

		if err = m.parseMessages(tx); err != nil {
			return errs.Internal{Cause: err.Error()}
		}
	}

	return nil
}

// parseMessages - parse messages from transaction
func (m *module) parseMessages(tx *types.Tx) error {
	for i, msg := range tx.Body.Messages {
		var stdMsg sdk.Msg
		if err := m.cdc.UnpackAny(msg, &stdMsg); err != nil {
			return fmt.Errorf("error while unpacking message: %s", err)
		}

		for _, module := range m.overgoldModules {
			if messageModule, ok := module.(modules.MessageModule); ok {
				if err := messageModule.HandleMsg(i, stdMsg, tx); err != nil {
					m.logger.MsgError(module, tx, stdMsg, err)
					return err
				}
			}
		}
	}

	return nil
}
