package distribution

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	juno "github.com/desmos-labs/juno/types"

	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/modules/distribution/utils"
)

// HandleMsg allows to handle the different utils related to the distribution module
func HandleMsg(tx *juno.Tx, msg sdk.Msg, client distrtypes.QueryClient, db *database.Db) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	if _, ok := msg.(*distrtypes.MsgFundCommunityPool); ok {
		return utils.UpdateCommunityPool(tx.Height, client, db)
	}

	return nil
}
