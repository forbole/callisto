package types

import (
	"time"
)

// ValidatorSigningInfo contains the signing info of a validator at a given height
type ValidatorSigningInfo struct {
	ValidatorAddress    string
	StartHeight         int64
	IndexOffset         int64
	JailedUntil         time.Time
	Tombstoned          bool
	MissedBlocksCounter int64
	Height              int64
}

// Equal tells whether v and w represent the same rows
func (v ValidatorSigningInfo) Equal(w ValidatorSigningInfo) bool {
	return v.ValidatorAddress == w.ValidatorAddress &&
		v.StartHeight == w.StartHeight &&
		v.IndexOffset == w.IndexOffset &&
		v.JailedUntil.Equal(w.JailedUntil) &&
		v.Tombstoned == w.Tombstoned &&
		v.MissedBlocksCounter == w.MissedBlocksCounter &&
		v.Height == w.Height
}

// ValidatorSigningInfo allows to build a new ValidatorSigningInfo
func NewValidatorSigningInfo(
	validatorAddress string,
	startHeight int64,
	indexOffset int64,
	jailedUntil time.Time,
	tombstoned bool,
	missedBlocksCounter int64,
	height int64,
) ValidatorSigningInfo {
	return ValidatorSigningInfo{
		ValidatorAddress:    validatorAddress,
		StartHeight:         startHeight,
		IndexOffset:         indexOffset,
		JailedUntil:         jailedUntil,
		Tombstoned:          tombstoned,
		MissedBlocksCounter: missedBlocksCounter,
		Height:              height,
	}
}
