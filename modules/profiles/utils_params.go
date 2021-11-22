package profiles

import (
	"github.com/forbole/bdjuno/v2/types"
	"github.com/rs/zerolog/log"
)

// updateParamd fetches from the REST APIs the latest value for the
// profile params, and saves it inside the database.
func (m *Module) UpdateParams() error {
	log.Debug().Str("module", "profiles").Str("operation", "profiles").
		Msg("getting profiles params")

	height, err := m.db.GetLastBlockHeight()
	if err != nil {
		return err
	}

	// Get the params
	params, err := m.source.GetParams(height)
	if err != nil {
		return err
	}

	return m.db.SaveProfilesParams(types.NewProfilesParams(params, height))
}
