package bank

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
	juno "github.com/forbole/juno/v5/types"
)

// HandleMsg implements MessageModule
func (m *Module) HandleMsg(index int, msg sdk.Msg, tx *juno.Tx) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	switch bankMsg := msg.(type) {
	case *bank.MsgSend:
		return m.handleMsgSend(tx, index, bankMsg)
	case *bank.MsgMultiSend:
		return m.handleMsgMultiSend(tx, index, bankMsg)
	default:
		return fmt.Errorf("unrecognized %s message type: %T", m.Name(), bankMsg)
	}
}
