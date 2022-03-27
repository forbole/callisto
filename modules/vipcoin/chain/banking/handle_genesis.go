/*
 * Copyright 2022 Business Process Technologies. All rights reserved.
 */

package banking

import (
	"encoding/json"

	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/rs/zerolog/log"
)

// HandleGenesis implements GenesisModule
func (m *Module) HandleGenesis(doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	log.Debug().Str("module", "banking").Msg("parsing genesis")

	return nil
}
