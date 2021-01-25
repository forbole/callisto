package database

import (
	"fmt"

	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"

	"github.com/forbole/bdjuno/x/slashing/types"
)

// SaveSlashingParams saves the slashing params for the given height
func (db *BigDipperDb) SaveSlashingParams(params slashingtypes.Params, height int64) error {
	stmt := `
INSERT INTO slashing_params 
    (signed_block_window, min_signed_per_window, downtime_jail_duration, slash_fraction_double_sign, slash_fraction_downtime, height) 
VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := db.Sql.Exec(stmt,
		params.SignedBlocksWindow, params.MinSignedPerWindow.String(), params.DowntimeJailDuration,
		params.SlashFractionDoubleSign.String(), params.SlashFractionDowntime.String(), height)
	return err
}

// SaveValidatorsSigningInfos saves the given infos inside the database
func (db *BigDipperDb) SaveValidatorsSigningInfos(infos []types.ValidatorSigningInfo) error {
	stmt := `
INSERT INTO validator_signing_info 
    (validator_address, start_height, index_offset, jailed_until, tombstoned, missed_blocks_counter, height)
VALUES `
	var args []interface{}

	for i, info := range infos {
		ii := i * 7

		stmt += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d),", ii+1, ii+2, ii+3, ii+4, ii+5, ii+6, ii+7)
		args = append(args,
			info.ValidatorAddress, info.StartHeight, info.IndexOffset, info.JailedUntil, info.Tombstoned,
			info.MissedBlocksCounter, info.Height,
		)
	}

	stmt = stmt[:len(stmt)-1] // Remove trailing ","
	stmt += ` ON CONFLICT DO NOTHING`
	_, err := db.Sql.Exec(stmt, args...)
	return err
}
