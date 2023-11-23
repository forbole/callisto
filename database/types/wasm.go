package types

import (
	"database/sql/driver"
	"fmt"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
)

// DbAccessConfig represents the information stored inside the database about a single access_config
type DbAccessConfig struct {
	Permission int      `db:"permission"`
	Address    []string `db:"address"`
}

// NewDbAccessConfig builds a DbAccessConfig starting from an CosmWasm type AccessConfig
func NewDbAccessConfig(accessCfg *wasmtypes.AccessConfig) DbAccessConfig {
	return DbAccessConfig{
		Permission: int(accessCfg.Permission),
		Address:    accessCfg.Addresses,
	}
}

// Value implements driver.Valuer
func (cfg *DbAccessConfig) Value() (driver.Value, error) {
	if cfg != nil {
		return fmt.Sprintf("(%d,%s)", cfg.Permission, cfg.Address), nil
	}

	return fmt.Sprintf("(%d,%s)", wasmtypes.AccessTypeUnspecified, ""), nil
}

// Equal tells whether a and b represent the same access_config
func (cfg *DbAccessConfig) Equal(b *DbAccessConfig) bool {
	return cfg.Permission == b.Permission
}

// ===================== Params =====================

// WasmParams represents the CosmWasm code in x/wasm module
type WasmParams struct {
	CodeUploadAccess             *DbAccessConfig `db:"code_upload_access"`
	InstantiateDefaultPermission int32           `db:"instantiate_default_permission"`
	Height                       int64           `db:"height"`
}

// NewWasmParams allows to build a new x/wasm params instance
func NewWasmParams(
	codeUploadAccess *DbAccessConfig, instantiateDefaultPermission int32, height int64,
) WasmParams {
	return WasmParams{
		CodeUploadAccess:             codeUploadAccess,
		InstantiateDefaultPermission: instantiateDefaultPermission,
		Height:                       height,
	}
}
