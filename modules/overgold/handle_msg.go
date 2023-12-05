package overgold

import (
	"errors"
	"fmt"
	"os"
	"time"

	"git.ooo.ua/vipcoin/lib/errs"
	"git.ooo.ua/vipcoin/lib/filter"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/forbole/juno/v5/modules"
	"github.com/forbole/juno/v5/types"

	dbtypes "github.com/forbole/bdjuno/v4/database/types"
)

// scheduler runs the scheduler
func (m *Module) scheduler() {
	for {
		lastBlock, err := m.lastBlockRepo.Get()
		if err != nil {
			m.logger.Error("Fail lastBlockRepo.Get", "module", m.Name(), "error", err)
			continue
		}

		lastBlock++

		if err = m.parseBlock(lastBlock); err != nil {
			time.Sleep(intervalLastBlock)

			if errors.As(err, &errs.NotFound{}) {
				continue
			}

			m.logger.Error("Fail parseBlock", "module", m.Name(), "error", err)
			continue
		}

		if err = m.lastBlockRepo.Update(lastBlock); err != nil {
			m.logger.Error("Fail lastBlockRepo.Update", "module", m.Name(), "error", err)
			os.Exit(1)
		}
	}
}

// parseBlock parse block
func (m *Module) parseBlock(lastBlock uint64) error {
	block, err := m.db.GetBlock(filter.NewFilter().SetArgument(dbtypes.FieldHeight, lastBlock))
	if err != nil {
		if errors.As(err, &errs.NotFound{}) {
			if block, _, err = m.parseMissingBlocksAndTransactions(int64(lastBlock)); err != nil {
				m.logger.Error("Fail parseMissingBlocksAndTransactions", "module", m.Name(), "error", err)
				return errs.Internal{Cause: "Fail parseMissingBlocksAndTransactions, error: " + err.Error()}
			}
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

// parseTx parse txs from block
func (m *Module) parseTx(block dbtypes.BlockRow) error {
	txs, err := m.db.GetTransactions(filter.NewFilter().SetArgument(dbtypes.FieldHeight, block.Height))
	if err != nil {
		if errors.As(err, &errs.NotFound{}) {
			if block, _, err = m.parseMissingBlocksAndTransactions(block.Height); err != nil {
				m.logger.Error("Fail parseMissingBlocksAndTransactions", "module", m.Name(), "error", err)
				return errs.Internal{Cause: "Fail parseMissingBlocksAndTransactions, error: " + err.Error()}
			}
			return err
		}

		return errs.Internal{Cause: err.Error()}
	}

	if err = block.CheckTxNumCount(int64(len(txs))); err != nil {
		if _, txs, err = m.parseMissingBlocksAndTransactions(block.Height); err != nil {
			m.logger.Error("Fail parseMissingBlocksAndTransactions", "module", m.Name(), "error", err)
			return errs.Internal{Cause: "Fail parseMissingBlocksAndTransactions, error: " + err.Error()}
		}

		if err = block.CheckTxNumCount(int64(len(txs))); err != nil {
			return err
		}
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
func (m *Module) parseMessages(tx *types.Tx) error {
	for i, msg := range tx.Body.Messages {
		var stdMsg sdk.Msg
		if err := m.cdc.UnpackAny(msg, &stdMsg); err != nil {
			return fmt.Errorf("error while an unpacking message: %s", err)
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
