package authority

import (
	authoritytypes "github.com/e-money/em-ledger/x/authority/types"
	"github.com/forbole/bdjuno/database"
	authorityutils "github.com/forbole/bdjuno/modules/authority/utils"

	sdk "github.com/cosmos/cosmos-sdk/types"
	juno "github.com/desmos-labs/juno/types"
)

// HandleMsg allows to handle the different msg under authority modules(for now it handles only MsgSetGasPrices)
func HandleMsg(
	tx *juno.Tx, msg sdk.Msg, authorityClient authoritytypes.QueryClient, db *database.Db,
) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	if msg, ok := msg.(*authoritytypes.MsgSetGasPrices); ok {
		return handleMsgSetGasprices(tx, msg, authorityClient, db)
	}

	return nil
}

// handleMsgSetGasprices handles storing the gas prices inside the database
func handleMsgSetGasprices(
	tx *juno.Tx, msg *authoritytypes.MsgSetGasPrices,
	authorityClient authoritytypes.QueryClient, db *database.Db,
) error {
	err := authorityutils.StoreSetGasPricesFromMessage(tx.Height, tx, msg, authorityClient, db)
	if err != nil {
		return err
	}

	return err
}
