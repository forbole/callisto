package types

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
