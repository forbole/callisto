package nft

import (
	"fmt"
	"strconv"

	nftTypes "github.com/Library-Genesis/libgen/x/nft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/forbole/bdjuno/v4/database"
	utils "github.com/forbole/bdjuno/v4/modules/utils"
	generalUtils "github.com/forbole/bdjuno/v4/utils"
	juno "github.com/forbole/juno/v4/types"
	"github.com/rs/zerolog/log"
)

// HandleMsg implements MessageModule
func (m *Module) HandleMsg(index int, msg sdk.Msg, tx *juno.Tx) error {
	if tx.Code != 0 {
		return nil
	}

	switch cosmosMsg := msg.(type) {
	case *nftTypes.MsgIssueDenom:
		return m.handleMsgIssueDenom(tx, cosmosMsg)
	case *nftTypes.MsgMintNFT:
		return m.handleMsgMintNFT(index, tx, cosmosMsg)
	case *nftTypes.MsgEditNFT:
		return m.handleMsgEditNFT(cosmosMsg)
	case *nftTypes.MsgBurnNFT:
		return m.handleMsgBurnNFT(index, tx, cosmosMsg)
	default:
		return nil
	}
}

func (m *Module) handleMsgIssueDenom(tx *juno.Tx, msg *nftTypes.MsgIssueDenom) error {
	log.Debug().Str("module", "nft").Str("denomId", msg.Id).Msg("handling message issue denom")

	return m.db.SaveDenom(tx.TxHash, msg.Id, msg.Name, msg.Schema, msg.Sender, msg.Uri)
}

func (m *Module) handleMsgMintNFT(index int, tx *juno.Tx, msg *nftTypes.MsgMintNFT) error {
	log.Debug().Str("module", "nft").Str("denomId", msg.DenomId).Str("name", msg.Name).Msg("handling message mint nft")

	tokenIDStr := utils.GetValueFromLogs(uint32(index), tx.Logs, nftTypes.EventTypeMintNFT, nftTypes.AttributeKeyTokenID)
	if tokenIDStr == "" {
		return fmt.Errorf("token id not found in tx %s", tx.TxHash)
	}

	tokenID, err := strconv.ParseUint(tokenIDStr, 10, 64)
	if err != nil {
		return err
	}

	timestamp, err := generalUtils.ISO8601ToTimestamp(tx.Timestamp)
	if err != nil {
		return err
	}

	return m.db.ExecuteTx(func(dbTx *database.DbTx) error {
		if err := dbTx.UpdateNFTHistory(tx.TxHash, tokenID, msg.DenomId, "0x0", msg.Sender, uint64(timestamp)); err != nil {
			return err
		}

		return dbTx.SaveNFT(tx.TxHash, tokenID, msg.DenomId, msg.Name, msg.Description, msg.URI, msg.Recipient, msg.Sender, msg.Recipient)
	})
}

func (m *Module) handleMsgEditNFT(msg *nftTypes.MsgEditNFT) error {
	log.Debug().Str("module", "nft").Str("denomId", msg.DenomId).Str("tokenId", msg.Id).Msg("handling message edit nft")

	dataJSON, dataText := utils.GetData(msg.Data)

	return m.db.UpdateNFT(msg.Id, msg.DenomId, msg.Name, msg.URI, utils.SanitizeUTF8(dataJSON), dataText)
}

func (m *Module) handleMsgBurnNFT(index int, tx *juno.Tx, msg *nftTypes.MsgBurnNFT) error {
	log.Debug().Str("module", "nft").Str("denomId", msg.DenomId).Str("tokenId", msg.Id).Msg("handling message burn nft")

	tokenIDStr := utils.GetValueFromLogs(uint32(index), tx.Logs, nftTypes.EventTypeBurnNFT, nftTypes.AttributeKeyTokenID)
	if tokenIDStr == "" {
		return fmt.Errorf("token id not found in tx %s", tx.TxHash)
	}

	tokenID, err := strconv.ParseUint(tokenIDStr, 10, 64)
	if err != nil {
		return err
	}

	timestamp, err := generalUtils.ISO8601ToTimestamp(tx.Timestamp)
	if err != nil {
		return err
	}

	return m.db.ExecuteTx(func(dbTx *database.DbTx) error {
		if error := dbTx.UpdateNFTHistory(tx.TxHash, tokenID, msg.DenomId, msg.Sender, "0x0", uint64(timestamp)); error != nil {
			return error
		}

		return dbTx.BurnNFT(msg.Id, msg.DenomId)
	})
}
