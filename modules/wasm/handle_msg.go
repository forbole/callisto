package wasm

import (
	"fmt"
	"strconv"
	"time"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/forbole/bdjuno/v3/types"
	juno "github.com/forbole/juno/v3/types"
)

// HandleMsg implements modules.MessageModule
func (m *Module) HandleMsg(index int, msg sdk.Msg, tx *juno.Tx) error {
	if len(tx.Logs) == 0 {
		return nil
	}

	switch cosmosMsg := msg.(type) {
	case *wasmtypes.MsgStoreCode:
		return m.HandleMsgStoreCode(index, tx, cosmosMsg)
	case *wasmtypes.MsgInstantiateContract:
		return m.HandleMsgInstantiateContract(index, tx, cosmosMsg)
	case *wasmtypes.MsgExecuteContract:
		return m.HandleMsgExecuteContract(index, tx, cosmosMsg)

	}

	return nil
}

// HandleMsgStoreCode allows to properly handle a MsgStoreCode
func (m *Module) HandleMsgStoreCode(index int, tx *juno.Tx, msg *wasmtypes.MsgStoreCode) error {
	// Get store code event
	event, err := tx.FindEventByType(index, wasmtypes.EventTypeStoreCode)
	if err != nil {
		return fmt.Errorf("error while searching for EventTypeInstantiate: %s", err)
	}

	// Get code ID from store code event
	codeIDKey, err := tx.FindAttributeByKey(event, wasmtypes.AttributeKeyCodeID)
	if err != nil {
		return fmt.Errorf("error while searching for AttributeKeyContractAddr: %s", err)
	}

	codeID, err := strconv.ParseInt(codeIDKey, 10, 64)
	if err != nil {
		return fmt.Errorf("error while parsing code id to int64: %s", err)
	}

	return m.db.SaveWasmCode(
		types.NewWasmCode(msg, codeID, tx.Height),
	)
}

// HandleMsgInstantiateContract allows to properly handle a MsgInstantiateContract
func (m *Module) HandleMsgInstantiateContract(index int, tx *juno.Tx, msg *wasmtypes.MsgInstantiateContract) error {
	// Get instantiate contract event
	event, err := tx.FindEventByType(index, wasmtypes.EventTypeInstantiate)
	if err != nil {
		return fmt.Errorf("error while searching for EventTypeInstantiate: %s", err)
	}

	// Get contract address
	contractAddress, err := tx.FindAttributeByKey(event, wasmtypes.AttributeKeyContractAddr)
	if err != nil {
		return fmt.Errorf("error while searching for AttributeKeyContractAddr: %s", err)
	}

	// Get result data
	resultData, err := tx.FindAttributeByKey(event, wasmtypes.AttributeKeyResultDataHex)
	if err != nil {
		return fmt.Errorf("error while searching for AttributeKeyResultDataHex: %s", err)
	}

	// Get the contract info
	contractInfo, err := m.source.GetContractInfo(tx.Height, contractAddress)
	if err != nil {
		return fmt.Errorf("error while getting proposal: %s", err)
	}

	timestamp, err := time.Parse(time.RFC3339, tx.Timestamp)
	if err != nil {
		return fmt.Errorf("error while parsing time: %s", err)
	}

	return m.db.SaveWasmContract(
		types.NewWasmContract(msg, contractAddress, []byte(resultData), timestamp, contractInfo.Extension, tx.Height),
	)
}

// HandleMsgExecuteContract allows to properly handle a MsgExecuteContract
func (m *Module) HandleMsgExecuteContract(index int, tx *juno.Tx, msg *wasmtypes.MsgExecuteContract) error {
	// Get Execute Contract event
	event, err := tx.FindEventByType(index, wasmtypes.EventTypeExecute)
	if err != nil {
		return fmt.Errorf("error while searching for EventTypeExecute: %s", err)
	}

	// Get result data
	resultData, err := tx.FindAttributeByKey(event, wasmtypes.AttributeKeyResultDataHex)
	if err != nil {
		return fmt.Errorf("error while searching for AttributeKeyResultDataHex: %s", err)
	}

	timestamp, err := time.Parse(time.RFC3339, tx.Timestamp)
	if err != nil {
		return fmt.Errorf("error while parsing time: %s", err)
	}

	//TO-DO save execute contract
	return m.db.SaveWasmExecuteContract(
		types.NewWasmExecuteContract(msg, []byte(resultData), timestamp, tx.Height),
	)
}
