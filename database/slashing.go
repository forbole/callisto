package database

import (
	"encoding/json"
	"fmt"

	"github.com/forbole/bdjuno/v4/types"
)

// SaveValidatorsSigningInfos saves the given infos inside the database
func (db *Db) SaveValidatorsSigningInfos(infos []types.ValidatorSigningInfo) error {
	if len(infos) == 0 {
		return nil
	}

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
	stmt += `
ON CONFLICT (validator_address) DO UPDATE 
	SET validator_address = excluded.validator_address,
		start_height = excluded.start_height,
		index_offset = excluded.index_offset,
		jailed_until = excluded.jailed_until,
		tombstoned = excluded.tombstoned,
		missed_blocks_counter = excluded.missed_blocks_counter,
		height = excluded.height
WHERE validator_signing_info.height <= excluded.height`

	_, err := db.SQL.Exec(stmt, args...)
	if err != nil {
		return fmt.Errorf("error while storing validators signing infos: %s", err)
	}

	return nil
}

// SaveSlashingParams saves the slashing params for the given height
func (db *Db) SaveSlashingParams(params *types.SlashingParams) error {
	paramsBz, err := json.Marshal(&params.Params)
	if err != nil {
		return err
	}

	stmt := `
INSERT INTO slashing_params (params, height) 
VALUES ($1, $2)
ON CONFLICT (one_row_id) DO UPDATE 
    SET params = excluded.params, 
        height = excluded.height
WHERE slashing_params.height <= excluded.height`

	_, err = db.SQL.Exec(stmt, string(paramsBz), params.Height)
	if err != nil {
		return fmt.Errorf("error while storing slashing params: %s", err)
	}

	return nil
}
