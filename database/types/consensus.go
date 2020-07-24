package types

// ConsensusRow represents a single row inside the consensus table
type ConsensusRow struct {
	Height int64  `db:"height"`
	Round  int    `db:"round"`
	Step   string `db:"step"`
}

// Equal tells whether r and s contain the same data
func (r ConsensusRow) Equal(s ConsensusRow) bool {
	return r.Height == s.Height &&
		r.Round == s.Round &&
		r.Step == s.Step
}
