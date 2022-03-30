/*
 * Copyright 2022 Business Process Technologies. All rights reserved.
 */

package banking

import (
	"encoding/json"

	bankingtypes "git.ooo.ua/vipcoin/chain/x/banking/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

// HandleGenesis implements GenesisModule
func (m *Module) HandleGenesis(doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	var bankingState bankingtypes.GenesisState
	if err := m.cdc.UnmarshalJSON(appState[bankingtypes.ModuleName], &bankingState); err != nil {
		return err
	}

	return nil
}
