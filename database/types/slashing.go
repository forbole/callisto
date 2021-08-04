package types

import "time"

// ValidatorSigningInfoRow represents a single row of the validator_signing_info table
type ValidatorSigningInfoRow struct {
	ValidatorAddress    string    `db:"validator_address"`
	StartHeight         int64     `db:"start_height"`
	IndexOffset         int64     `db:"index_offset"`
	JailedUntil         time.Time `db:"jailed_until"`
	Tombstoned          bool      `db:"tombstoned"`
	MissedBlocksCounter int64     `db:"missed_blocks_counter"`
	Height              int64     `db:"height"`
}

// Equal tells whether v and w represent the same rows
func (v ValidatorSigningInfoRow) Equal(w ValidatorSigningInfoRow) bool {
	return v.ValidatorAddress == w.ValidatorAddress &&
		v.StartHeight == w.StartHeight &&
		v.IndexOffset == w.IndexOffset &&
		v.JailedUntil.Equal(w.JailedUntil) &&
		v.Tombstoned == w.Tombstoned &&
		v.MissedBlocksCounter == w.MissedBlocksCounter &&
		v.Height == w.Height
}

// NewValidatorSigningInfoRow allows to build a new ValidatorSigningInfoRow
func NewValidatorSigningInfoRow(
	validatorAddress string,
	startHeight int64,
	indexOffset int64,
	jailedUntil time.Time,
	tombstoned bool,
	missedBlocksCounter int64,
	height int64,
) ValidatorSigningInfoRow {
	return ValidatorSigningInfoRow{
		ValidatorAddress:    validatorAddress,
		StartHeight:         startHeight,
		IndexOffset:         indexOffset,
		JailedUntil:         jailedUntil,
		Tombstoned:          tombstoned,
		MissedBlocksCounter: missedBlocksCounter,
		Height:              height,
	}
}

// -------------------------------------------------------------------------------------------------------------------

// SlashingParamsRow represents a single row inside the slashing_params table
type SlashingParamsRow struct {
	OneRowID bool   `db:"one_row_id"`
	Params   string `db:"params"`
	Height   int64  `db:"height"`
}

// NewSlashingParamsRow allows to create a new instance SlashingParamsRow
func NewSlashingParamsRow(params string, height int64) SlashingParamsRow {
	return SlashingParamsRow{
		OneRowID: true,
		Params:   params,
		Height:   height,
	}
}
