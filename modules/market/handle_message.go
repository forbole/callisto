package market

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	junotypes "github.com/forbole/juno/v3/types"
	markettypes "github.com/ovrclk/akash/x/market/types/v1beta2"
)

// HandleMsg implements modules.MessageModule
func (m *Module) HandleMsg(_ int, msg sdk.Msg, tx *junotypes.Tx) error {
	switch cosmosMsg := msg.(type) {
	case *markettypes.MsgCreateLease:

	case *markettypes.MsgCloseLease:
	}

	return nil
}
