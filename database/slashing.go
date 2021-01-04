package database

import (
	"fmt"

	"github.com/forbole/bdjuno/x/slashing/types"
)

// SaveValidatorsSigningInfos saves the given infos inside the database
func (db *BigDipperDb) SaveValidatorsSigningInfos(infos []types.ValidatorSigningInfo) error {
	stmt := `
INSERT INTO validator_signing_info 
    (validator_address, start_height, index_offset, jailed_until, tombstoned, missed_blocks_counter, height, timestamp) 
VALUES `
	var args []interface{}

	for i, info := range infos {
		ii := i * 8

		stmt += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d),",
			ii+1, ii+2, ii+3, ii+4, ii+5, ii+6, ii+7, ii+8)
		args = append(args,
			info.ValidatorAddress.String(), info.StartHeight, info.IndexOffset, info.JailedUntil, info.Tombstoned,
			info.MissedBlocksCounter, info.Height, info.Timestamp,
		)
	}

	stmt = stmt[:len(stmt)-1] // Remove trailing ","
	stmt += ` ON CONFLICT DO NOTHING`
	_, err := db.Sql.Exec(stmt, args...)
	return err
}
