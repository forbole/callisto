package stakers

import (
	"fmt"

	stakerstypes "github.com/KYVENetwork/chain/x/stakers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/forbole/bdjuno/v4/types"
	juno "github.com/forbole/juno/v4/types"
)

// HandleMsg implements modules.MessageModule
func (m *Module) HandleMsg(index int, msg sdk.Msg, tx *juno.Tx) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	switch cosmosMsg := msg.(type) {
	case *stakerstypes.MsgUpdateParams:
		return m.handleMsgUpdateParams(tx, cosmosMsg)
	}

	return nil
}

// handleMsgUpdateParams allows to properly handle a MsgUpdateParams
func (m *Module) handleMsgUpdateParams(tx *juno.Tx, msg *stakerstypes.MsgUpdateParams) error {
	params, err := m.source.Params(tx.Height)
	if err != nil {
		return fmt.Errorf("error while getting stakers params: %s", err)
	}

	return m.db.SaveStakersParams(
		types.NewStakersParams(params, tx.Height))
}
