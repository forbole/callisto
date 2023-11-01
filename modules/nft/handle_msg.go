package nft

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/forbole/bdjuno/v4/database"
	"github.com/forbole/bdjuno/v4/modules/utils"
	myTypes "github.com/forbole/bdjuno/v4/x/nft/types"
	juno "github.com/forbole/juno/v4/types"
	"github.com/rs/zerolog/log"
	"strconv"
)

// HandleMsg implements MessageModule
func (m *Module) HandleMsg(index int, msg sdk.Msg, tx *juno.Tx) error {
	if tx.Code != 0 {
		return nil
	}

	switch cosmosMsg := msg.(type) {
	case *myTypes.MsgIssueDenom:
		return m.handleMsgIssueDenom(tx, cosmosMsg)
	case *myTypes.MsgEditNFT:
		return m.handleMsgEditNFT(cosmosMsg)
	case *myTypes.MsgMintNFT:
		return m.handleMsgMintNFT(index, tx, cosmosMsg)
	default:
		return nil
	}
}

func (m *Module) handleMsgIssueDenom(tx *juno.Tx, msg *myTypes.MsgIssueDenom) error {
	log.Debug().
		Str("module", "nft").
		Str("denomId", msg.Id).
		Msg("handling message issue denom")

	return m.db.SaveDenom(tx.TxHash, msg.Id, msg.Name, msg.Schema, msg.Sender, msg.Uri)
}

func (m *Module) handleMsgEditNFT(msg *myTypes.MsgEditNFT) error {
	log.Debug().
		Str("module", "nft").
		Str("denomId", msg.DenomId).
		Str("tokenId", msg.Id).
		Msg("handling message edit nft")

	return m.db.UpdateNFT(msg.Id, msg.DenomId, msg.Name, msg.URI, msg.Description)
}

func (m *Module) handleMsgMintNFT(index int, tx *juno.Tx, msg *myTypes.MsgMintNFT) error {
	log.Debug().
		Str("module", "nft").
		Str("denomId", msg.DenomId).
		Str("name", msg.Name).
		Msg("handling message mint nft")

	tokenIDStr := utils.GetValueFromLogs(uint32(index), tx.Logs, myTypes.EventTypeMintNFT, myTypes.AttributeKeyTokenID)
	if tokenIDStr == "" {
		return fmt.Errorf("token id not found in tx %s", tx.TxHash)
	}

	tokenID, err := strconv.ParseUint(tokenIDStr, 10, 64)
	if err != nil {
		return err
	}

	//timestamp, err := generalUtils.ISO8601ToTimestamp(tx.Timestamp)
	//if err != nil {
	//	return err
	//}

	return m.db.ExecuteTx(func(dbTx *database.DbTx) error {
		//if err := dbTx.UpdateNFTHistory(tx.TxHash, tokenID, msg.DenomId, "0x0", msg.Sender, uint64(timestamp)); err != nil {
		//	return err
		//}

		return dbTx.SaveNFT(tx.TxHash, tokenID, msg.DenomId, msg.Name, msg.Description, msg.URI, msg.Tags, msg.Sender, msg.Recipient)
	})
}
