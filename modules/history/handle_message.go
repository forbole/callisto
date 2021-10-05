package history

import (
	"fmt"
	"time"

	juno "github.com/desmos-labs/juno/v2/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/forbole/bdjuno/v2/modules/utils"
)

// HandleMsg implements modules.MessageModule
func (m *Module) HandleMsg(_ int, msg sdk.Msg, tx *juno.Tx) error {
	timestamp, err := time.Parse(time.RFC3339, tx.Timestamp)
	if err != nil {
		return fmt.Errorf("error while parsing time: %s", err)
	}

	addresses, err := m.getAddresses(m.cdc, msg)
	if err != nil {
		return fmt.Errorf("error while getting accounts after message of type %s", msg.Type())
	}

	for _, address := range utils.FilterNonAccountAddresses(addresses) {
		err = m.UpdateAccountBalanceHistoryWithTime(address, timestamp)
		if err != nil {
			return fmt.Errorf("error while updating account balance history: %s", err)
		}
	}

	return nil
}
