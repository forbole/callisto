package history

import (
	"fmt"

	"github.com/gogo/protobuf/proto"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/desmos-labs/juno/modules/messages"
	"github.com/gogo/protobuf/proto"

	"github.com/forbole/bdjuno/database"
	historyutils "github.com/forbole/bdjuno/modules/history/utils"
	"github.com/forbole/bdjuno/modules/utils"
)

// HandleMsg handles any message updating the involved accounts
func HandleMsg(msg sdk.Msg, getAddresses messages.MessageAddressesParser, cdc codec.Codec, db *database.Db) error {
	addresses, err := getAddresses(cdc, msg)
	if err != nil {
		return fmt.Errorf("error while getting accounts after message of type %s", proto.MessageName(msg))
	}

	for _, address := range utils.FilterNonAccountAddresses(addresses) {
		err = historyutils.UpdateAccountBalanceHistory(address, db)
		if err != nil {
			return err
		}
	}

	return nil
}
