/*
 * Copyright 2022 Business Process Technologies. All rights reserved.
 */

package accounts

import (
	"encoding/json"

	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/rs/zerolog/log"
)

// HandleGenesis implements GenesisModule
func (m *Module) HandleGenesis(doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	log.Debug().Str("module", "accounts").Msg("parsing genesis")

	return nil
}
