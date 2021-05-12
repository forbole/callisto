package types

// DoubleSignEvidence represent a double sign evidence on each tendermint block
type DoubleSignEvidence struct {
	VoteA  DoubleSignVote
	VoteB  DoubleSignVote
	Height int64
}

// NewDoubleSignEvidence return a new DoubleSignEvidence object
func NewDoubleSignEvidence(height int64, voteA DoubleSignVote, voteB DoubleSignVote) DoubleSignEvidence {
	return DoubleSignEvidence{
		VoteA:  voteA,
		VoteB:  voteB,
		Height: height,
	}
}

// DoubleSignVote represents a double vote which is included inside a DoubleSignEvidence
type DoubleSignVote struct {
	BlockID          string
	ValidatorAddress string
	Signature        string
	Type             int
	Height           int64
	Round            int32
	ValidatorIndex   int32
}

// NewDoubleSignVote allows to create a new DoubleSignVote instance
func NewDoubleSignVote(
	roundType int,
	height int64,
	round int32,
	blockID string,
	validatorAddress string,
	validatorIndex int32,
	signature string,
) DoubleSignVote {
	return DoubleSignVote{
		Type:             roundType,
		Height:           height,
		Round:            round,
		BlockID:          blockID,
		ValidatorAddress: validatorAddress,
		ValidatorIndex:   validatorIndex,
		Signature:        signature,
	}
}
