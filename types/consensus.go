package types

import "time"

// Genesis contains the useful information about the genesis
type Genesis struct {
	ChainID       string
	Time          time.Time
	InitialHeight int64
}

// NewGenesis allows to build a new Genesis instance
func NewGenesis(chainID string, startTime time.Time, initialHeight int64) *Genesis {
	return &Genesis{
		ChainID:       chainID,
		Time:          startTime,
		InitialHeight: initialHeight,
	}
}

// Equal returns true iff g and other contain the same data
func (g *Genesis) Equal(other *Genesis) bool {
	return g.ChainID == other.ChainID &&
		g.Time.Equal(other.Time) &&
		g.InitialHeight == other.InitialHeight
}

// ------------------------------------------------------------------------------------------------------------------

// ConsensusEvent represents a consensus event
type ConsensusEvent struct {
	Height int64  `json:"height"`
	Round  int32  `json:"round"`
	Step   string `json:"step"`
}

// NewConsensusEvent allows to easily build a new ConsensusEvent object
func NewConsensusEvent(height int64, round int32, step string) *ConsensusEvent {
	return &ConsensusEvent{
		Height: height,
		Round:  round,
		Step:   step,
	}
}

// Equal tells whether c and other contain the same data
func (c ConsensusEvent) Equal(other ConsensusEvent) bool {
	return c.Height == other.Height &&
		c.Round == other.Round &&
		c.Step == other.Step
}
