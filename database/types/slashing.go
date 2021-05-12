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

// ValidatorSigningInfoRow allows to build a new ValidatorSigningInfoRow
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

type SlashingParamsRow struct {
	OneRowID                bool   `db:"one_row_id"`
	SignedBlockWindow       int64  `db:"signed_block_window"`
	MinSignedPerWindow      string `db:"min_signed_per_window"`
	DowntimeJailDuration    int64  `db:"downtime_jail_duration"`
	SlashFractionDoubleSign string `db:"slash_fraction_double_sign"`
	SlashFractionDowntime   string `db:"slash_fraction_downtime"`
	Height                  int64  `db:"height"`
}

func NewSlashingParamsRow(
	signedBlockPerWindow int64, minSignedPerWindow string, downtimeJailDuration int64,
	slashFractionDoubleSign string, slashFractionDownTime string, height int64,
) SlashingParamsRow {
	return SlashingParamsRow{
		OneRowID:                true,
		SignedBlockWindow:       signedBlockPerWindow,
		MinSignedPerWindow:      minSignedPerWindow,
		DowntimeJailDuration:    downtimeJailDuration,
		SlashFractionDoubleSign: slashFractionDoubleSign,
		SlashFractionDowntime:   slashFractionDownTime,
		Height:                  height,
	}
}

func (p SlashingParamsRow) Equal(q SlashingParamsRow) bool {
	return p.SignedBlockWindow == q.SignedBlockWindow &&
		p.MinSignedPerWindow == q.MinSignedPerWindow &&
		p.DowntimeJailDuration == q.DowntimeJailDuration &&
		p.SlashFractionDowntime == q.SlashFractionDowntime &&
		p.SlashFractionDoubleSign == q.SlashFractionDoubleSign &&
		p.Height == q.Height
}
