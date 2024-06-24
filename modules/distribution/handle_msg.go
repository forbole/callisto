package distribution

import (
	"fmt"

	juno "github.com/forbole/juno/v6/types"
	"github.com/rs/zerolog/log"
)

var msgFilter = map[string]bool{
	"/cosmos.distribution.v1beta1.MsgFundCommunityPool": true,
}

// HandleMsgExec implements modules.AuthzMessageModule
func (m *Module) HandleMsgExec(index int, _ int, executedMsg juno.Message, tx *juno.Transaction) error {
	return m.HandleMsg(index, executedMsg, tx)
}

// HandleMsg implements modules.MessageModule
func (m *Module) HandleMsg(_ int, msg juno.Message, tx *juno.Transaction) error {
	if _, ok := msgFilter[msg.GetType()]; !ok {
		return nil
	}

	log.Debug().Str("module", "distribution").Str("hash", tx.TxHash).Uint64("height", tx.Height).Msg(fmt.Sprintf("handling distribution message %s", msg.GetType()))

	if msg.GetType() == "/cosmos.distribution.v1beta1.MsgFundCommunityPool" {
		return m.updateCommunityPool(int64(tx.Height))
	}
	return nil
}
