package types

import (
	"time"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/forbole/callisto/v4/types"
)

// ===================== Params =====================

// WasmParamsRow represents a single row inside the wasm_params table
type WasmParamsRow struct {
	OneRowID bool   `db:"one_row_id"`
	Params   string `db:"params"`
	Height   int64  `db:"height"`
}

// ===================== Code =====================

// WasmCodeRow represents a single row inside the "wasm_code" table
type WasmCodeRow struct {
	Sender                string `db:"sender"`
	WasmByteCode          string `db:"byte_code"`
	InstantiatePermission string `db:"instantiate_permission"`
	CodeID                uint64 `db:"code_id"`
	Height                int64  `db:"height"`
}

// NewWasmCodeRow allows to easily create a new NewWasmCodeRow
func NewWasmCodeRow(
	sender string,
	wasmByteCode string,
	instantiatePermission string,
	codeID uint64,
	height int64,
) WasmCodeRow {
	return WasmCodeRow{
		Sender:                sender,
		WasmByteCode:          wasmByteCode,
		InstantiatePermission: instantiatePermission,
		CodeID:                codeID,
		Height:                height,
	}
}

// ===================== Contract =====================

// WasmContractRow represents a single row inside the "wasm_contract" table
type WasmContractRow struct {
	Sender                string    `db:"sender"`
	Creator               string    `db:"creator"`
	Admin                 string    `db:"admin"`
	CodeID                uint64    `db:"code_id"`
	Label                 string    `db:"label"`
	RawContractMessage    string    `db:"raw_contract_message"`
	Funds                 *DbCoins  `db:"funds"`
	ContractAddress       string    `db:"contract_address"`
	Data                  string    `db:"data"`
	InstantiatedAt        time.Time `db:"instantiated_at"`
	ContractInfoExtension string    `db:"contract_info_extension"`
	ContractStates        string    `db:"contract_states"`
	Height                int64     `db:"height"`
}

// NewWasmContractRow allows to easily create a new WasmContractRow
func NewWasmContractRow(
	sender string,
	admin string,
	codeID uint64,
	label string,
	rawMsg wasmtypes.RawContractMessage,
	funds *DbCoins,
	contractAddress string,
	data string,
	instantiatedAt time.Time,
	creator string,
	contractInfoExtension string,
	height int64,
) WasmContractRow {
	rawContractMsg := types.ConvertRawContractMessage(rawMsg)

	return WasmContractRow{
		Sender:                sender,
		Admin:                 admin,
		CodeID:                codeID,
		Label:                 label,
		RawContractMessage:    string(rawContractMsg),
		Funds:                 funds,
		ContractAddress:       contractAddress,
		Data:                  data,
		InstantiatedAt:        instantiatedAt,
		Creator:               creator,
		ContractInfoExtension: contractInfoExtension,
		Height:                height,
	}
}

// Equals return true if one WasmContractRow representing the same row as the original one
func (a WasmContractRow) Equals(b WasmContractRow) bool {
	return a.Sender == b.Sender &&
		a.Creator == b.Creator &&
		a.Admin == b.Admin &&
		a.CodeID == b.CodeID &&
		a.Label == b.Label &&
		a.RawContractMessage == b.RawContractMessage &&
		a.Funds.Equal(b.Funds) &&
		a.ContractAddress == b.ContractAddress &&
		a.Data == b.Data &&
		a.InstantiatedAt.Equal(b.InstantiatedAt) &&
		a.ContractInfoExtension == b.ContractInfoExtension &&
		a.Height == b.Height
}

// ===================== Execute Contract =====================

// WasmExecuteContractRow represents a single row inside the "wasm_execute_contract" table
type WasmExecuteContractRow struct {
	Sender             string    `db:"sender"`
	ContractAddress    string    `db:"contract_address"`
	RawContractMessage string    `db:"raw_contract_message"`
	Funds              *DbCoins  `db:"funds"`
	Data               string    `db:"data"`
	ExecutedAt         time.Time `db:"executed_at"`
	Height             int64     `db:"height"`
}

// NewWasmExecuteContractRow allows to easily create a new WasmExecuteContractRow
func NewWasmExecuteContractRow(
	sender string,
	contractAddress string,
	rawMsg wasmtypes.RawContractMessage,
	funds *DbCoins,
	data string,
	executedAt time.Time,
	height int64,
) WasmExecuteContractRow {
	rawContractMsg := types.ConvertRawContractMessage(rawMsg)

	return WasmExecuteContractRow{
		Sender:             sender,
		RawContractMessage: string(rawContractMsg),
		Funds:              funds,
		ContractAddress:    contractAddress,
		Data:               data,
		ExecutedAt:         executedAt,
		Height:             height,
	}
}

// Equals return true if one WasmExecuteContractRow representing the same row as the original one
func (a WasmExecuteContractRow) Equals(b WasmExecuteContractRow) bool {
	return a.Sender == b.Sender &&
		a.ContractAddress == b.ContractAddress &&
		a.RawContractMessage == b.RawContractMessage &&
		a.Funds.Equal(b.Funds) &&
		a.Data == b.Data &&
		a.ExecutedAt.Equal(b.ExecutedAt) &&
		a.Height == b.Height
}
