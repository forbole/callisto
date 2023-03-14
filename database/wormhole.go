package database

import (
	"encoding/json"
	"fmt"

	"github.com/forbole/bdjuno/v3/types"
)

// SaveWormholeConfig allows to store the given config inside the database
func (db *Db) SaveWormholeConfig(config *types.WormholeConfig) error {
	configBz, err := json.Marshal(&config.Config)
	if err != nil {
		return fmt.Errorf("error while marshaling wormhole config: %s", err)
	}

	stmt := `
INSERT INTO wormhole_config (config, height) 
VALUES ($1, $2)
ON CONFLICT (one_row_id) DO UPDATE 
    SET config = excluded.config,
        height = excluded.height
WHERE wormhole_config.height <= excluded.height`

	_, err = db.SQL.Exec(stmt, string(configBz), config.Height)
	if err != nil {
		return fmt.Errorf("error while storing wormhole config: %s", err)
	}

	return nil
}
