package database

import (
	constypes "github.com/forbole/bdjuno/x/consensus/types"
	"github.com/rs/zerolog/log"
)

// SaveConsensus allows to properly store the given consensus event into the database.
// Note that only one consensus event is allowed inside the database at any time.
func (db BigDipperDb) SaveConsensus(event constypes.ConsensusEvent) error {
	log.Debug().
		Str("module", "consensus").
		Int64("height", event.Height).
		Int("round", event.Round).
		Str("step", event.Step).
		Msg("saving consensus")

	// Delete all the existing events
	stmt := `DELETE FROM consensus WHERE true`
	if _, err := db.Sql.Exec(stmt); err != nil {
		return err
	}

	stmt = `INSERT INTO consensus (height, round, step) VALUES ($1, $2, $3)`
	_, err := db.Sql.Exec(stmt, event.Height, event.Round, event.Step)
	return err
}
