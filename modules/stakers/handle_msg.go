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

	switch msg.(type) {
	case *stakerstypes.MsgUpdateParams:
		return m.handleMsgUpdateParams(tx)
	case *stakerstypes.MsgUpdateCommission:
		return m.handleMsgUpdateCommission(tx)
	case *stakerstypes.MsgUpdateMetadata:
		return m.handleMsgUpdateMetadata(tx)

	}

	return nil
}

// handleMsgUpdateParams allows to properly handle a MsgUpdateParams
func (m *Module) handleMsgUpdateParams(tx *juno.Tx) error {
	params, err := m.source.Params(tx.Height)
	if err != nil {
		return fmt.Errorf("error while getting stakers params: %s", err)
	}

	return m.db.SaveStakersParams(
		types.NewStakersParams(params, tx.Height))
}

// handleMsgUpdateCommission allows to properly handle a MsgUpdateCommission
func (m *Module) handleMsgUpdateCommission(tx *juno.Tx) error {
	// refresh commission for all protocol validators
	err := m.UpdateProtocolValidatorsCommission(tx.Height)
	if err != nil {
		return fmt.Errorf("error while updating protocol validators commission: %s", err)
	}

	return nil
}

// handleMsgUpdateMetadata allows to properly handle a MsgUpdateMetadata
func (m *Module) handleMsgUpdateMetadata(tx *juno.Tx) error {
	// refresh description for all protocol validators
	err := m.UpdateProtocolValidatorsDescription(tx.Height)
	if err != nil {
		return fmt.Errorf("error while updating protocol validators description: %s", err)
	}

	return nil
}
