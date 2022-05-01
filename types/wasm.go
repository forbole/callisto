package types

import (
	"time"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// WasmCode represents the CosmWasm code in x/wasm module
type WasmCode struct {
	Sender                string
	WasmByteCode          []byte
	InstantiatePermission *wasmtypes.AccessConfig
	CodeID                int64
	Height                int64
}

// NewWasmCode allows to build a new x/wasm code instance from wasmtypes.MsgStoreCode
func NewWasmCode(msg *wasmtypes.MsgStoreCode, codeID int64, height int64) WasmCode {
	return WasmCode{
		Sender:                msg.Sender,
		WasmByteCode:          msg.WASMByteCode,
		InstantiatePermission: msg.InstantiatePermission,
		CodeID:                codeID,
		Height:                height,
	}
}

// WasmContract represents the CosmWasm contract in x/wasm module
type WasmContract struct {
	Sender                string
	Admin                 string
	CodeID                uint64
	Label                 string
	RawContractMsg        []byte
	Funds                 sdk.Coins
	ContractAddress       string
	Data                  []byte
	InstantiatedAt        time.Time
	ContractInfoExtension wasmtypes.ContractInfoExtension
	Height                int64
}

// NewWasmCode allows to build a new x/wasm contract instance from wasmtypes.MsgStoreCode
func NewWasmContract(
	msg *wasmtypes.MsgInstantiateContract, contractAddress string, data []byte,
	instantiatedAt time.Time, contractInfoExtension wasmtypes.ContractInfoExtension, height int64,
) WasmContract {
	return WasmContract{
		Sender:                msg.Sender,
		Admin:                 msg.Admin,
		CodeID:                msg.CodeID,
		Label:                 msg.Label,
		RawContractMsg:        msg.Msg,
		Funds:                 msg.Funds,
		ContractAddress:       contractAddress,
		Data:                  data,
		InstantiatedAt:        instantiatedAt,
		ContractInfoExtension: contractInfoExtension,
		Height:                height,
	}
}

// WasmExecuteContract represents the CosmWasm execute contract in x/wasm module
type WasmExecuteContract struct {
	Sender          string
	ContractAddress string
	RawContractMsg  []byte
	Funds           sdk.Coins
	Data            []byte
	ExecutedAt      time.Time
	Height          int64
}

// NewWasmExecuteContract allows to build a new x/wasm execute contract instance from wasmtypes.MsgExecuteContract
func NewWasmExecuteContract(
	msg *wasmtypes.MsgExecuteContract, data []byte,
	executedAt time.Time, height int64,
) WasmExecuteContract {
	return WasmExecuteContract{
		Sender:          msg.Sender,
		ContractAddress: msg.Contract,
		RawContractMsg:  msg.Msg,
		Funds:           msg.Funds,
		Data:            data,
		ExecutedAt:      executedAt,
		Height:          height,
	}
}
