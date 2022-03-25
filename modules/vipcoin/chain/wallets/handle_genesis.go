/*
 * Copyright 2022 Business Process Technologies. All rights reserved.
 */

package wallets

import (
	"encoding/json"

	"github.com/rs/zerolog/log"
	tmtypes "github.com/tendermint/tendermint/types"
)

// HandleGenesis implements modules.GenesisModule
func (m *Module) HandleGenesis(doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	log.Debug().Str("module", "wallets").Msg("parsing genesis")

	return nil
}
