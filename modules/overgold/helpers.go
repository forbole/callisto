package overgold

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"git.ooo.ua/vipcoin/lib/errs"
	tmctypes "github.com/cometbft/cometbft/rpc/core/types"
	tmtypes "github.com/cometbft/cometbft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/forbole/juno/v5/modules"
	"github.com/forbole/juno/v5/types"
	txtypes "github.com/forbole/juno/v5/types"

	dbhelpers "github.com/forbole/bdjuno/v4/database/types"
	dbtypes "github.com/forbole/bdjuno/v4/database/types"
)

const (
	intervalLastBlock = time.Second

	module = "overgold"
)

// parseMissingBlocksAndTransactions - parse missing blocks and transactions
func (m *Module) parseMissingBlocksAndTransactions(height int64) (dbtypes.BlockRow, []*txtypes.Tx, error) {
	block, err := m.node.Block(height)
	if err != nil { // skip height > current
		if !strings.Contains(err.Error(), "must be less than or equal to the current blockchain height") {
			return dbtypes.BlockRow{}, []*txtypes.Tx{}, fmt.Errorf("failed to get block from node: %s", err)
		}
	}

	events, err := m.node.BlockResults(height)
	if err != nil {
		return dbtypes.BlockRow{}, []*txtypes.Tx{}, fmt.Errorf("failed to get block results from node: %s", err)
	}

	txs, err := m.node.Txs(block)
	if err != nil {
		return dbtypes.BlockRow{}, []*txtypes.Tx{}, fmt.Errorf("failed to get transactions for block: %s", err)
	}

	vals, err := m.node.Validators(height)
	if err != nil {
		return dbtypes.BlockRow{}, []*txtypes.Tx{}, fmt.Errorf("failed to get validators for block: %s", err)
	}

	return m.ExportBlock(block, events, txs, vals)
}

// ExportBlock accepts a finalized block and a corresponding set of transactions
// and persists them to the database along with attributable metadata. An error
// is returned if the writing fails.
func (m *Module) ExportBlock(
	b *tmctypes.ResultBlock, r *tmctypes.ResultBlockResults, txs []*types.Tx, vals *tmctypes.ResultValidators,
) (dbtypes.BlockRow, []*txtypes.Tx, error) {
	// Save all validators
	err := m.SaveValidators(vals.Validators)
	if err != nil {
		return dbtypes.BlockRow{}, []*txtypes.Tx{}, err
	}

	// Make sure the proposer exists
	proposerAddr := sdk.ConsAddress(b.Block.ProposerAddress).String()
	val := findValidatorByAddr(proposerAddr, vals)
	if val == nil {
		return dbtypes.BlockRow{}, []*txtypes.Tx{},
			fmt.Errorf("failed to find validator by proposer address %s: %s", proposerAddr, err)
	}

	block := types.NewBlockFromTmBlock(b, sumGasTxs(txs))

	// Save the block
	err = m.db.SaveBlock(block)
	if err != nil {
		return dbtypes.BlockRow{}, []*txtypes.Tx{}, fmt.Errorf("failed to persist block: %s", err)
	}

	// Save the commits
	err = m.ExportCommit(b.Block.LastCommit, vals)
	if err != nil {
		return dbtypes.BlockRow{}, []*txtypes.Tx{}, err
	}

	// Call the block handlers
	for _, mod := range m.overgoldModules {
		if blockModule, ok := mod.(modules.BlockModule); ok {
			err = blockModule.HandleBlock(b, r, txs, vals)
			if err != nil {
				m.logger.BlockError(mod, b, err)
			}
		}
	}

	tsx, err := m.ExportTxs(txs)
	if err != nil {
		return dbtypes.BlockRow{}, []*txtypes.Tx{}, err
	}

	// Export the transactions
	return dbtypes.BlockRow{
		Height:          block.Height,
		Hash:            block.Hash,
		TxNum:           int64(block.TxNum),
		TotalGas:        int64(block.TotalGas),
		ProposerAddress: dbhelpers.ToNullString(block.ProposerAddress),
		Timestamp:       block.Timestamp,
	}, tsx, nil
}

// SaveValidators persists a list of Tendermint validators with an address and a
// consensus public key. An error is returned if the public key cannot be Bech32
// encoded or if the DB write fails.
func (m *Module) SaveValidators(vals []*tmtypes.Validator) error {
	var validators = make([]*types.Validator, len(vals))
	for index, val := range vals {
		consAddr := sdk.ConsAddress(val.Address).String()

		consPubKey, err := types.ConvertValidatorPubKeyToBech32String(val.PubKey)
		if err != nil {
			return fmt.Errorf("failed to convert validator public key for validators %s: %s", consAddr, err)
		}

		validators[index] = types.NewValidator(consAddr, consPubKey)
	}

	err := m.db.SaveValidators(validators)
	if err != nil {
		return fmt.Errorf("error while saving validators: %s", err)
	}

	return nil
}

// findValidatorByAddr finds a validator by a consensus address given a set of
// Tendermint validators for a particular block. If no validator is found, nil
// is returned.
func findValidatorByAddr(consAddr string, vals *tmctypes.ResultValidators) *tmtypes.Validator {
	for _, val := range vals.Validators {
		if consAddr == sdk.ConsAddress(val.Address).String() {
			return val
		}
	}

	return nil
}

// sumGasTxs returns the total gas consumed by a set of transactions.
func sumGasTxs(txs []*types.Tx) uint64 {
	var totalGas uint64

	for _, tx := range txs {
		totalGas += uint64(tx.GasUsed)
	}

	return totalGas
}

// ExportCommit accepts a block commitment and a corresponding set of
// validators for the commitment and persists them to the database. An error is
// returned if any writing fails or if there is any missing-aggregated data.
func (m *Module) ExportCommit(commit *tmtypes.Commit, vals *tmctypes.ResultValidators) error {
	var signatures []*types.CommitSig
	for _, commitSig := range commit.Signatures {
		// Avoid empty commits
		if commitSig.Signature == nil {
			continue
		}

		valAddr := sdk.ConsAddress(commitSig.ValidatorAddress)
		val := findValidatorByAddr(valAddr.String(), vals)
		if val == nil {
			return fmt.Errorf("failed to find validator by commit validator address %s", valAddr.String())
		}

		signatures = append(signatures, types.NewCommitSig(
			types.ConvertValidatorAddressToBech32String(commitSig.ValidatorAddress),
			val.VotingPower,
			val.ProposerPriority,
			commit.Height,
			commitSig.Timestamp,
		))
	}

	err := m.db.SaveCommitSignatures(signatures)
	if err != nil {
		return fmt.Errorf("error while saving commit signatures: %s", err)
	}

	return nil
}

// ExportTxs accepts a slice of transactions and persists then inside the database.
// An error is returned if the write fails.
func (m *Module) ExportTxs(txs []*types.Tx) ([]*txtypes.Tx, error) {
	// Handle all the transactions inside the block
	for _, tx := range txs {
		// Save the transaction itself
		err := m.db.SaveTx(tx)
		if err != nil {
			return []*txtypes.Tx{}, fmt.Errorf("failed to handle transaction with hash %s: %s", tx.TxHash, err)
		}

		// Call the tx handlers
		for _, module := range m.overgoldModules {
			if transactionModule, ok := module.(modules.TransactionModule); ok {
				err = transactionModule.HandleTx(tx)
				if err != nil {
					m.logger.TxError(module, tx, err)
				}
			}
		}

		// Handle all the messages contained inside the transaction
		for _, msg := range tx.Body.Messages {
			var stdMsg sdk.Msg
			err = m.cdc.UnpackAny(msg, &stdMsg)
			if err != nil {
				return []*txtypes.Tx{}, fmt.Errorf("error while unpacking message: %s", err)
			}

			// Call the handlers
			for i, mod := range m.overgoldModules {
				if messageModule, ok := mod.(modules.MessageModule); ok {
					err = messageModule.HandleMsg(i, stdMsg, tx)
					if err != nil {
						if errors.As(err, &errs.NotFound{}) {
							continue
						}

						m.logger.MsgError(mod, tx, stdMsg, err)
					}
				}
			}
		}
	}

	return txs, nil
}
