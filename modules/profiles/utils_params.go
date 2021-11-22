package profiles

import (
	profilestypes "github.com/desmos-labs/desmos/x/profiles/types"
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

// SaveGenesisParams saves the profiles parameters at genesis
func (m *Module) SaveGenesisParams(params profilestypes.Params, height int64) error {
	return m.db.SaveProfilesParams(types.NewProfilesParams(
		params,
		height,
	))
}
