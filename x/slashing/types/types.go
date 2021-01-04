package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ValidatorSigningInfo struct {
	ValidatorAddress    sdk.ConsAddress
	StartHeight         int64
	IndexOffset         int64
	JailedUntil         time.Time
	Tombstoned          bool
	MissedBlocksCounter int64
	Height              int64
	Timestamp           time.Time
}

// Equal tells whether v and w represent the same rows
func (v ValidatorSigningInfo) Equal(w ValidatorSigningInfo) bool {
	return v.ValidatorAddress.Equals(w.ValidatorAddress) &&
		v.StartHeight == w.StartHeight &&
		v.IndexOffset == w.IndexOffset &&
		v.JailedUntil.Equal(w.JailedUntil) &&
		v.Tombstoned == w.Tombstoned &&
		v.MissedBlocksCounter == w.MissedBlocksCounter &&
		v.Height == w.Height &&
		v.Timestamp.Equal(w.Timestamp)
}

// ValidatorSigningInfo allows to build a new ValidatorSigningInfo
func NewValidatorSigningInfo(
	validatorAddress sdk.ConsAddress,
	startHeight int64,
	indexOffset int64,
	jailedUntil time.Time,
	tombstoned bool,
	missedBlocksCounter int64,
	height int64,
	timestamp time.Time) ValidatorSigningInfo {
	return ValidatorSigningInfo{
		ValidatorAddress:    validatorAddress,
		StartHeight:         startHeight,
		IndexOffset:         indexOffset,
		JailedUntil:         jailedUntil,
		Tombstoned:          tombstoned,
		MissedBlocksCounter: missedBlocksCounter,
		Height:              height,
		Timestamp:           timestamp,
	}
}
