package issuer

import (
	"github.com/forbole/bdjuno/database"

	issuertypes "github.com/e-money/em-ledger/x/issuer/types"
	issuerutils "github.com/forbole/bdjuno/modules/issuer/utils"

	sdk "github.com/cosmos/cosmos-sdk/types"
	juno "github.com/desmos-labs/juno/types"
)

// HandleMsg allows to handle MsgSetInflation
func HandleMsg(
	tx *juno.Tx, index int, msg sdk.Msg, issuerClient issuertypes.QueryClient, db *database.Db,
) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	if cosmosMsg, ok := msg.(*issuertypes.MsgSetInflation); ok {
		return handleMsgSetInflation(tx, cosmosMsg, db)
	}

	return nil
}

// ---------------------------------------------------------------------------------------------------------------------

// handleMsgSetInflation stores the inflation data inside the database
func handleMsgSetInflation(
	tx *juno.Tx, msg *issuertypes.MsgSetInflation, db *database.Db,
) error {
	err := issuerutils.StoreEmoneyInflationFromMessage(tx.Height, msg, db)
	if err != nil {
		return err
	}
	return nil
}
