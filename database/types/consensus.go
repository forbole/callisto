package types

import (
	"database/sql"
	"time"
)

type GenesisRow struct {
	OneRowID      bool      `db:"one_row_id"`
	ChainID       string    `db:"chain_id"`
	Time          time.Time `db:"time"`
	InitialHeight int64     `db:"initial_height"`
}

func NewGenesisRow(chainID string, time time.Time, initialHeight int64) GenesisRow {
	return GenesisRow{
		OneRowID:      true,
		ChainID:       chainID,
		Time:          time,
		InitialHeight: initialHeight,
	}
}

func (r GenesisRow) Equal(s GenesisRow) bool {
	return r.Time.Equal(s.Time) &&
		r.ChainID == s.ChainID &&
		r.InitialHeight == s.InitialHeight
}

// -------------------------------------------------------------------------------------------------------------------

// ConsensusRow represents a single row inside the consensus table
type ConsensusRow struct {
	Step     string `db:"step"`
	Height   int64  `db:"height"`
	Round    int32  `db:"round"`
	OneRowID bool   `db:"one_row_id"`
}

// NewConsensusRow allows to build a new ConsensusRow instance
func NewConsensusRow(height int64, round int32, step string) ConsensusRow {
	return ConsensusRow{
		OneRowID: true,
		Height:   height,
		Round:    round,
		Step:     step,
	}
}

// Equal tells whether r and s contain the same data
func (r ConsensusRow) Equal(s ConsensusRow) bool {
	return r.Height == s.Height &&
		r.Round == s.Round &&
		r.Step == s.Step
}

// AverageTimeRow is the average block time each minute/hour/day
type AverageTimeRow struct {
	OneRowID    bool    `db:"one_row_id"`
	AverageTime float64 `db:"average_time"`
	Height      int64   `db:"height"`
}

func NewAverageTimeRow(averageTime float64, height int64) AverageTimeRow {
	return AverageTimeRow{
		OneRowID:    true,
		AverageTime: averageTime,
		Height:      height,
	}
}

// Equal return true if two AverageTimeRow are true
func (r AverageTimeRow) Equal(s AverageTimeRow) bool {
	return r.AverageTime == s.AverageTime &&
		r.Height == s.Height
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

// -------------------------------------------------------------------------------------------------------------------

// BlockHeightAndTimestamp represents last block height and timestamp stored inside the database
type BlockHeightAndTimestamp struct {
	Height         int64     `db:"height"`
	BlockTimestamp time.Time `db:"timestamp"`
}
