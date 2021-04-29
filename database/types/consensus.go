package types

import (
	"database/sql"
	"time"
)

type GenesisRow struct {
	ChainID string    `db:"chain_id"`
	Time    time.Time `db:"time"`
}

func NewGenesisRow(chainID string, time time.Time) GenesisRow {
	return GenesisRow{
		ChainID: chainID,
		Time:    time,
	}
}

func (r GenesisRow) Equal(s GenesisRow) bool {
	return r.Time.Equal(s.Time) &&
		r.ChainID == s.ChainID
}

// -------------------------------------------------------------------------------------------------------------------

// ConsensusRow represents a single row inside the consensus table
type ConsensusRow struct {
	Height int64  `db:"height"`
	Round  int32  `db:"round"`
	Step   string `db:"step"`
}

// Equal tells whether r and s contain the same data
func (r ConsensusRow) Equal(s ConsensusRow) bool {
	return r.Height == s.Height &&
		r.Round == s.Round &&
		r.Step == s.Step
}

// -------------------------------------------------------------------------------------------------------------------

// AverageTimeRow represents a row inside the average_block_time_per_minute/hour/day table
type AverageTimeRow struct {
	AverageTime float64 `db:"average_time"`
}

func NewBlockTimeRow(averageTime float64) AverageTimeRow {
	return AverageTimeRow{
		AverageTime: averageTime,
	}
}

// Equal return true if two AverageTimeRow are true
func (r AverageTimeRow) Equal(s AverageTimeRow) bool {
	return r.AverageTime == s.AverageTime
}

// -------------------------------------------------------------------------------------------------------------------

// BlockRow represents a single block row stored inside the database
type BlockRow struct {
	Height          int64          `db:"height"`
	Hash            string         `db:"hash"`
	TxNum           int64          `db:"num_txs"`
	TotalGas        int64          `db:"total_gas"`
	ProposerAddress sql.NullString `db:"proposer_address"`
	PreCommitsNum   int64          `db:"pre_commits"`
	Timestamp       time.Time      `db:"timestamp"`
}
