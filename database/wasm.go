package database

import (
	"fmt"

	dbtypes "github.com/forbole/bdjuno/v3/database/types"
	"github.com/forbole/bdjuno/v3/types"
	"github.com/lib/pq"
)

// SaveWasmCode allows to store the wasm code from MsgStoreCode
func (db *Db) SaveWasmCode(wasmCode types.WasmCode) error {

	stmt := `
INSERT INTO wasm_code(sender, byte_code, instantiate_permission, code_id, height) 
VALUES ($1, $2, $3, $4, $5) 
ON CONFLICT DO NOTHING`

	// TO-DO: check if string(wasmCode.WasmByteCode) saved as string in DB

	_, err := db.Sql.Exec(stmt,
		wasmCode.Sender, string(wasmCode.WasmByteCode),
		pq.Array(dbtypes.NewDbAccessConfig(wasmCode.InstantiatePermission)),
		wasmCode.CodeID, wasmCode.Height,
	)
	if err != nil {
		return fmt.Errorf("error while saving wasm code: %s", err)
	}

	return nil
}

// SaveWasmContract allows to store the wasm contract from MsgInstantiateContract
func (db *Db) SaveWasmContract(wasmContract types.WasmContract) error {

	stmt := `
INSERT INTO wasm_contract 
(sender, admin, code_id, label, raw_contract_message, funds, contract_address, data, instantiated_at, contract_info_extension, height) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) 
ON CONFLICT DO NOTHING`

	ExtensionBz, err := db.EncodingConfig.Marshaler.MarshalJSON(wasmContract.ContractInfoExtension)
	if err != nil {
		return fmt.Errorf("error while marshaling contract info extension: %s", err)
	}

	// TO-DO: check if the below is stored as Json in DB:
	// - Data
	// - ContractInfoExtension

	_, err = db.Sql.Exec(stmt,
		wasmContract.Sender, wasmContract.Admin, wasmContract.CodeID, wasmContract.Label, wasmContract.RawContractMsg,
		pq.Array(dbtypes.NewDbCoins(wasmContract.Funds)), wasmContract.ContractAddress, string(wasmContract.Data),
		wasmContract.InstantiatedAt, string(ExtensionBz), wasmContract.Height,
	)

	if err != nil {
		return fmt.Errorf("error while saving wasm contract: %s", err)
	}

	return nil
}

// SaveWasmExecuteContract allows to store the wasm contract from MsgExecuteeContract
func (db *Db) SaveWasmExecuteContract(executeContract types.WasmExecuteContract) error {

	stmt := `
INSERT INTO wasm_execute_contract 
(sender, contract_address, raw_contract_message, funds, data, executed_at, height) 
VALUES ($1, $2, $3, $4, $5, $6, $7) 
ON CONFLICT DO NOTHING`

	// TO-DO: check if the below is stored as Json in DB:
	// - Data

	_, err := db.Sql.Exec(stmt,
		executeContract.Sender, executeContract.ContractAddress, executeContract.RawContractMsg,
		pq.Array(dbtypes.NewDbCoins(executeContract.Funds)), string(executeContract.Data),
		executeContract.ExecutedAt, executeContract.Height,
	)

	if err != nil {
		return fmt.Errorf("error while saving wasm contract: %s", err)
	}

	return nil
}
