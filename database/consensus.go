package database

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/forbole/bdjuno/v2/types"
	"github.com/lib/pq"

	dbtypes "github.com/forbole/bdjuno/v2/database/types"
	"github.com/forbole/bdjuno/v2/utils"
	junotypes "github.com/forbole/juno/v2/types"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
)

// GetLastBlock returns the last block stored inside the database based on the heights
func (db *Db) GetLastBlock() (*dbtypes.BlockRow, error) {
	stmt := `SELECT * FROM block ORDER BY height DESC LIMIT 1`

	var blocks []dbtypes.BlockRow
	if err := db.Sqlx.Select(&blocks, stmt); err != nil {
		return nil, err
	}

	if len(blocks) == 0 {
		return nil, fmt.Errorf("cannot get block, no blocks saved")
	}

	return &blocks[0], nil
}

// GetLastBlockHeight returns the last block height stored inside the database
func (db *Db) GetLastBlockHeight() (int64, error) {
	block, err := db.GetLastBlock()
	if err != nil {
		return 0, err
	}
	if block == nil {
		return 0, fmt.Errorf("block table is empty")
	}
	return block.Height, nil
}

// -------------------------------------------------------------------------------------------------------------------

// getBlockHeightTime retrieves the block at the specific time
func (db *Db) getBlockHeightTime(pastTime time.Time) (dbtypes.BlockRow, error) {
	stmt := `SELECT * FROM block WHERE block.timestamp <= $1 ORDER BY block.timestamp DESC LIMIT 1;`

	var val []dbtypes.BlockRow
	if err := db.Sqlx.Select(&val, stmt, pastTime); err != nil {
		return dbtypes.BlockRow{}, err
	}

	if len(val) == 0 {
		return dbtypes.BlockRow{}, fmt.Errorf("cannot get block time, no blocks saved")
	}

	return val[0], nil
}

// GetBlockHeightTimeMinuteAgo return block height and time that a block proposals
// about a minute ago from input date
func (db *Db) GetBlockHeightTimeMinuteAgo(now time.Time) (dbtypes.BlockRow, error) {
	pastTime := now.Add(time.Minute * -1)
	return db.getBlockHeightTime(pastTime)
}

// GetBlockHeightTimeHourAgo return block height and time that a block proposals
// about a hour ago from input date
func (db *Db) GetBlockHeightTimeHourAgo(now time.Time) (dbtypes.BlockRow, error) {
	pastTime := now.Add(time.Hour * -1)
	return db.getBlockHeightTime(pastTime)
}

// GetBlockHeightTimeDayAgo return block height and time that a block proposals
// about a day (24hour) ago from input date
func (db *Db) GetBlockHeightTimeDayAgo(now time.Time) (dbtypes.BlockRow, error) {
	pastTime := now.Add(time.Hour * -24)
	return db.getBlockHeightTime(pastTime)
}

// -------------------------------------------------------------------------------------------------------------------

// SaveAverageBlockTimePerMin save the average block time in average_block_time_per_minute table
func (db *Db) SaveAverageBlockTimePerMin(averageTime float64, height int64) error {
	stmt := `
INSERT INTO average_block_time_per_minute(average_time, height) 
VALUES ($1, $2) 
ON CONFLICT (one_row_id) DO UPDATE 
    SET average_time = excluded.average_time, 
        height = excluded.height
WHERE average_block_time_per_minute.height <= excluded.height`

	_, err := db.Sqlx.Exec(stmt, averageTime, height)
	if err != nil {
		return fmt.Errorf("error while storing average block time per minute: %s", err)
	}

	return nil
}

// SaveAverageBlockTimePerHour save the average block time in average_block_time_per_hour table
func (db *Db) SaveAverageBlockTimePerHour(averageTime float64, height int64) error {
	stmt := `
INSERT INTO average_block_time_per_hour(average_time, height) 
VALUES ($1, $2) 
ON CONFLICT (one_row_id) DO UPDATE 
    SET average_time = excluded.average_time,
        height = excluded.height
WHERE average_block_time_per_hour.height <= excluded.height`

	_, err := db.Sqlx.Exec(stmt, averageTime, height)
	if err != nil {
		return fmt.Errorf("error while storing average block time per hour: %s", err)
	}

	return nil
}

// SaveAverageBlockTimePerDay save the average block time in average_block_time_per_day table
func (db *Db) SaveAverageBlockTimePerDay(averageTime float64, height int64) error {
	stmt := `
INSERT INTO average_block_time_per_day(average_time, height) 
VALUES ($1, $2)
ON CONFLICT (one_row_id) DO UPDATE 
    SET average_time = excluded.average_time,
        height = excluded.height
WHERE average_block_time_per_day.height <= excluded.height`

	_, err := db.Sqlx.Exec(stmt, averageTime, height)
	if err != nil {
		return fmt.Errorf("error while storing average block time per day: %s", err)
	}

	return nil
}

// SaveAverageBlockTimeGenesis save the average block time in average_block_time_from_genesis table
func (db *Db) SaveAverageBlockTimeGenesis(averageTime float64, height int64) error {
	stmt := `
INSERT INTO average_block_time_from_genesis(average_time ,height) 
VALUES ($1, $2) 
ON CONFLICT (one_row_id) DO UPDATE 
    SET average_time = excluded.average_time, 
        height = excluded.height
WHERE average_block_time_from_genesis.height <= excluded.height`

	_, err := db.Sqlx.Exec(stmt, averageTime, height)
	if err != nil {
		return fmt.Errorf("error while storing average block time since genesis: %s", err)
	}

	return nil
}

// -------------------------------------------------------------------------------------------------------------------

// SaveGenesis save the given genesis data
func (db *Db) SaveGenesis(genesis *types.Genesis) error {
	stmt := `
INSERT INTO genesis(time, chain_id, initial_height) 
VALUES ($1, $2, $3) ON CONFLICT (one_row_id) DO UPDATE 
    SET time = excluded.time,
        initial_height = excluded.initial_height,
        chain_id = excluded.chain_id`

	_, err := db.Sqlx.Exec(stmt, genesis.Time, genesis.ChainID, genesis.InitialHeight)
	if err != nil {
		return fmt.Errorf("error while storing genesis: %s", err)
	}

	return nil
}

// GetGenesis returns the genesis information stored inside the database
func (db *Db) GetGenesis() (*types.Genesis, error) {
	var rows []*dbtypes.GenesisRow
	err := db.Sqlx.Select(&rows, `SELECT * FROM genesis;`)
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return nil, fmt.Errorf("no rows inside the genesis table")
	}

	row := rows[0]
	return types.NewGenesis(row.ChainID, row.Time, row.InitialHeight), nil
}

// CheckIfBlockIsMissing checks if block is stored in database and returns boolean value
func (db *Db) CheckIfBlockIsMissing(height int64) bool {
	var block []dbtypes.BlockRow
	stmt := `SELECT * FROM block WHERE height = $1`

	err := db.Sqlx.Select(&block, stmt, height)
	if err != nil {
		return true
	}
	if len(block) != 0 {
		return false
	}

	return true
}

// UpdateBlocksInDatabase updates given block in database
func (db *Db) UpdateBlockInDatabase(block *tmctypes.ResultBlock, blockResults *tmctypes.ResultBlockResults) error {
	stmt := `
INSERT INTO block(height, hash, num_txs, total_gas, proposer_address, timestamp)
VALUES ($1, $2, $3, $4, $5, $6) 
ON CONFLICT DO NOTHING`

	proposerAddress := sdk.ConsAddress(block.Block.ProposerAddress).String()
	txResults := blockResults.TxsResults

	var totalGasUsed int64
	for _, tx := range txResults {
		totalGasUsed += tx.GetGasUsed()
	}

	_, err := db.Sqlx.Exec(stmt, block.Block.Height, strings.ToUpper(hex.EncodeToString(block.Block.Hash())), len(block.Block.Txs), totalGasUsed, proposerAddress, block.Block.Time)
	if err != nil {
		return fmt.Errorf("error while storing block %v for proposer %s, error:  %s", block.Block.Height, proposerAddress, err)
	}

	return nil
}

func (db *Db) UpdateTxInDatabase(i int, tx *junotypes.Tx) error {
	stmt := `
INSERT INTO transaction(hash, height, success, messages, memo, signatures, signer_infos, fee, gas_wanted, gas_used, raw_log, logs)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
ON CONFLICT DO NOTHING`

	message, err := codec.ProtoMarshalJSON(tx.Body.Messages[i], nil)
	if err != nil {
		return fmt.Errorf("error while marshaling message: %s", err)
	}

	fee, err := codec.ProtoMarshalJSON(tx.AuthInfo.Fee, nil)
	if err != nil {
		return fmt.Errorf("error while marshaling fees: %s", err)
	}

	logs, err := json.Marshal(&tx.Logs)
	if err != nil {
		return fmt.Errorf("error while marshaling logs: %s", err)
	}

	signers, err := codec.ProtoMarshalJSON(tx.AuthInfo.SignerInfos[0], nil)
	if err != nil {
		return fmt.Errorf("error while marshaling signers: %s", err)
	}

	var signatures []string
	for _, signature := range tx.Signatures {
		signature := signature
		eachSignature, err := json.Marshal(&signature)
		if err != nil {
			return fmt.Errorf("error while marshaling signatures: %s", err)
		}
		signatures = append(signatures, string(eachSignature))
	}

	_, err = db.Sqlx.Exec(stmt,
		tx.TxHash,
		tx.Height,
		tx.Successful(),
		message,
		tx.GetBody().Memo,
		pq.StringArray(signatures),
		signers,
		fee,
		tx.GasWanted,
		tx.GasUsed,
		tx.RawLog,
		logs)
	if err != nil {
		return fmt.Errorf("error while storing tx, error:  %s", err)
	}

	err = db.UpdateMsgsInDatabase(tx, tx.TxHash, i, tx.Body.Messages[i].TypeUrl, message)
	if err != nil {
		return fmt.Errorf("error while updating tx message in database: %s", err)
	}

	return nil
}

func (db *Db) UpdateMsgsInDatabase(tx *junotypes.Tx, txHash string, i int, typeURL string, message []byte) error {
	stmt := `
	INSERT INTO message(transaction_hash, index, type, value, involved_accounts_addresses)
	VALUES ($1, $2, $3, $4, $5)
	ON CONFLICT DO NOTHING`
	var eventTypes []string
	var involvedAccounts []string
	var attributeKeys []string

	for _, event := range tx.Logs {
		for _, eventType := range event.Events {
			eventTypess := eventType.Type
			eventTypes = append(eventTypes, eventTypess)
			attributes := eventType.Attributes
			for _, attribute := range attributes {
				attributeKeys = append(attributeKeys, attribute.Key)
			}
		}
	}
	attributeKeys = utils.RemoveDuplicateValues(attributeKeys)

	for _, eventType := range eventTypes {
		event, _ := tx.FindEventByType(i, eventType)
		for _, key := range attributeKeys {
			address, _ := tx.FindAttributeByKey(event, key)
			if len(address) >= 40 {
				involvedAccounts = append(involvedAccounts, address)
			}
		}
		involvedAccounts = utils.RemoveDuplicateValues(involvedAccounts)
	}

	_, err := db.Sqlx.Exec(stmt,
		txHash,
		i,
		typeURL,
		message,
		pq.StringArray(involvedAccounts))
	if err != nil {
		return fmt.Errorf("error while storing message in database, error:  %s", err)
	}

	return nil
}
