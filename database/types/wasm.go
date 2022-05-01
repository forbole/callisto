package types

import (
	"time"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
)

// ===================== Store Code =====================

// DbAccessConfig represents the information stored inside the database about a single access_config
type DbAccessConfig struct {
	Permission int
	Address    string
}

// NewDbAccessConfig builds a DbAccessConfig starting from an CosmWasm type AccessConfig
func NewDbAccessConfig(accessCfg *wasmtypes.AccessConfig) DbAccessConfig {
	return DbAccessConfig{
		Permission: int(accessCfg.Permission),
		Address:    accessCfg.Address,
	}
}

// Equal tells whether a and b represent the same access_config
func (a *DbAccessConfig) Equal(b *DbAccessConfig) bool {
	return a.Address == b.Address && a.Permission == b.Permission
}

// WasmCodeRow represents a single row inside the "wasm_code" table
type WasmCodeRow struct {
	Sender                string          `db:"sender"`
	WasmByteCode          string          `db:"wasm_byte_code"`
	InstantiatePermission *DbAccessConfig `db:"instantiate_permission"`
	CodeID                int64           `db:"code_id"`
	Height                int64           `db:"height"`
}

// NewWasmCodeRow allows to easily create a new NewWasmCodeRow
func NewWasmCodeRow(
	sender string,
	wasmByteCode string,
	instantiatePermission *DbAccessConfig,
	codeID int64,
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

// Equals return true if one WasmCodeRow representing the same row as the original one
func (a WasmCodeRow) Equals(b WasmCodeRow) bool {
	return a.Sender == b.Sender &&
		a.WasmByteCode == b.WasmByteCode &&
		a.InstantiatePermission.Equal(b.InstantiatePermission) &&
		a.CodeID == b.CodeID &&
		a.Height == b.Height
}

// ===================== Wasm Contract =====================

// WasmContractRow represents a single row inside the "wasm_contract" table
type WasmContractRow struct {
	Sender                string    `db:"sender"`
	Admin                 string    `db:"admin"`
	CodeID                int64     `db:"code_id"`
	Label                 string    `db:"label"`
	RawContractMessage    string    `db:"raw_contract_message"`
	Funds                 *DbCoins  `db:"funds"`
	ContractAddress       string    `db:"contract_address"`
	Data                  string    `db:"data"`
	InstantiatedAt        time.Time `db:"instantiated_at"`
	ContractInfoExtension string    `db:"contract_info_extension"`
	Height                int64     `db:"height"`
}

// NewWasmContractRow allows to easily create a new NewWasmContractRow
func NewWasmContractRow(
	sender string,
	admin string,
	codeID int64,
	label string,
	rawContractMessage string,
	funds *DbCoins,
	contractAddress string,
	data string,
	instantiatedAt time.Time,
	contractInfoExtension string,
	height int64,
) WasmContractRow {
	return WasmContractRow{
		Sender:                sender,
		Admin:                 admin,
		CodeID:                codeID,
		Label:                 label,
		RawContractMessage:    rawContractMessage,
		Funds:                 funds,
		ContractAddress:       contractAddress,
		Data:                  data,
		InstantiatedAt:        instantiatedAt,
		ContractInfoExtension: contractInfoExtension,
		Height:                height,
	}
}

// Equals return true if one WasmContractRow representing the same row as the original one
func (a WasmContractRow) Equals(b WasmContractRow) bool {
	return a.Sender == b.Sender &&
		a.Admin == b.Admin &&
		a.CodeID == b.CodeID &&
		a.Label == b.Label &&
		a.RawContractMessage == b.RawContractMessage &&
		a.Funds.Equal(a.Funds) &&
		a.ContractAddress == b.ContractAddress &&
		a.Data == b.Data &&
		a.InstantiatedAt == b.InstantiatedAt &&
		a.ContractInfoExtension == b.ContractInfoExtension &&
		a.Height == b.Height
}
