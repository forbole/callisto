package database

import (
	"time"

	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/rs/zerolog/log"
)

// SaveStakingPool allows to save for the given height the given staking pool
func (db BigDipperDb) SaveStakingPool(height int64, pool staking.Pool) error {
	log.Debug().Int64("height", height).Msg("updating staking pool")

	statement := `INSERT INTO staking_pool (timestamp, height, bonded_tokens, not_bonded_tokens) 
				  VALUES ($1, $2, $3, $4)`
	_, err := db.Sql.Exec(statement,
		time.Now().UTC(), height, pool.BondedTokens.Int64(), pool.NotBondedTokens.Int64())
	return err
}
