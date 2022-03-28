/*
 * Copyright 2022 Business Process Technologies. All rights reserved.
 */

package wallets

import (
	"encoding/json"
  
	"github.com/rs/zerolog/log"
	tmtypes "github.com/tendermint/tendermint/types"
  walletstypes "git.ooo.ua/vipcoin/chain/x/wallets/types"
)

// HandleGenesis implements modules.GenesisModule
func (m *Module) HandleGenesis(doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	log.Debug().Str("module", "wallets").Msg("parsing genesis")

	// Unmarshal the bank state
	var walletsState walletstypes.GenesisState
	if err := m.cdc.UnmarshalJSON(appState[walletstypes.ModuleName], &walletsState); err != nil {
		return err
	}

	return m.walletsRepo.SaveWallets(walletsState.Wallets...)
}
