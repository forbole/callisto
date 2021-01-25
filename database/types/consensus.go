package types

import (
	"time"
)

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

// BlockTimeRow is the average block time each minute/hour/day
type BlockTimeRow struct {
	AverageTime float64 `db:"average_time"`
	Height      int64   `db:"height"`
}

func NewBlockTimeRow(averageTime float64, height int64) BlockTimeRow {
	return BlockTimeRow{
		AverageTime: averageTime,
		Height:      height,
	}
}

// Equal return true if two BlockTimeRow are true
func (r BlockTimeRow) Equal(s BlockTimeRow) bool {
	return r.AverageTime == s.AverageTime &&
		r.Height == s.Height
}

//Container to return block needed in certain height
type BlockRow struct {
	Height          int64     `db:"height"`
	Hash            string    `db:"hash"`
	TxNum           int64     `db:"num_txs"`
	TotalGas        int64     `db:"total_gas"`
	ProposerAddress string    `db:"proposer_address"`
	PreCommitsNum   int64     `db:"pre_commits"`
	Timestamp       time.Time `db:"timestamp"`
}
