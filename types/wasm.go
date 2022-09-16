package types

import (
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
)

// WasmParams represents the CosmWasm code in x/wasm module
type WasmParams struct {
	CodeUploadAccess             *wasmtypes.AccessConfig
	InstantiateDefaultPermission int32
	MaxWasmCodeSize              uint64
	Height                       int64
}

// NewWasmParams allows to build a new x/wasm params instance
func NewWasmParams(
	codeUploadAccess *wasmtypes.AccessConfig, instantiateDefaultPermission int32, maxWasmCodeSize uint64, height int64,
) WasmParams {
	return WasmParams{
		CodeUploadAccess:             codeUploadAccess,
		InstantiateDefaultPermission: instantiateDefaultPermission,
		MaxWasmCodeSize:              maxWasmCodeSize,
		Height:                       height,
	}
}
