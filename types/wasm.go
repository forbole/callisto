package types

import (
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
)

// WasmCode represents the Code in x/wasm module
type WasmCode struct {
	Sender                string
	WasmByteCode          []byte
	InstantiatePermission *wasmtypes.AccessConfig
	CodeID                int64
	Height                int64
}

// NewWasmCode allows to build a new x/wasm Code instance
func NewWasmCode(msg *wasmtypes.MsgStoreCode, codeID int64, height int64) WasmCode {
	return WasmCode{
		Sender:                msg.Sender,
		WasmByteCode:          msg.WASMByteCode,
		InstantiatePermission: msg.InstantiatePermission,
		CodeID:                codeID,
		Height:                height,
	}
}
