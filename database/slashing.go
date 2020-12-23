package database

import (
	"github.com/forbole/bdjunog/x/slashing/types"
	dbtypes "github.com/forbole/bdjuno/database/types"

)

func (db BigDipperDb) SaveValidatorSigningInfo(t types.ValidatorSigningInfo) error {
	stmt = "INSERT INTO validator_signing_info(validator_address,start_height,index_offset,jailed_until,tombstoned,missed_blocks_counter,height,timestamp) VALUES ($1,$2,$3,$4,$5,$6,$7,$8) ON CONFLICT DO NOTHING"
	_, err := db.Sql.Exec(stmt, t.ValidatorAddress.String(),
		t.StartHeight,
		t.IndexOffset,
		t.JailedUntil,
		t.Tombstoned,
		t.MissedBlocksCounter,
		t.Height,
		t.Timestamp)
	if err != nil {
		return err
	}
}
