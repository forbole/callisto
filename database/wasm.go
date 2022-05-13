package database

import (
	"fmt"

	dbtypes "github.com/forbole/bdjuno/v3/database/types"
	dbutils "github.com/forbole/bdjuno/v3/database/utils"
	"github.com/forbole/bdjuno/v3/types"
	"github.com/lib/pq"
)

// SaveWasmParams allows to store the wasm params
func (db *Db) SaveWasmParams(params types.WasmParams) error {
	stmt := `
INSERT INTO wasm_params(code_upload_access, instantiate_default_permission, max_wasm_code_size, height) 
VALUES ($1, $2, $3, $4) 
ON CONFLICT (one_row_id) DO UPDATE 
	SET code_upload_access = excluded.code_upload_access, 
		instantiate_default_permission = excluded.instantiate_default_permission, 
		max_wasm_code_size = excluded.max_wasm_code_size 
WHERE wasm_params.height <= excluded.height
`
	accessConfig := dbtypes.NewDbAccessConfig(params.CodeUploadAccess)
	cfgValue, _ := accessConfig.Value()

	_, err := db.Sql.Exec(stmt,
		cfgValue, params.InstantiateDefaultPermission, params.MaxWasmCodeSize, params.Height,
	)
	if err != nil {
		return fmt.Errorf("error while saving wasm params: %s", err)
	}

	return nil
}

// SaveWasmCode allows to store a single wasm code
func (db *Db) SaveWasmCode(wasmCode types.WasmCode) error {
	return db.SaveWasmCodes([]types.WasmCode{wasmCode})
}

// SaveWasmCodes allows to store the wasm code slice
func (db *Db) SaveWasmCodes(wasmCodes []types.WasmCode) error {
	stmt := `
INSERT INTO wasm_code(sender, byte_code, instantiate_permission, code_id, height) 
VALUES `

	// TO-DO: check if string(wasmCode.WasmByteCode) saved as string in DB

	var args []interface{}
	for i, code := range wasmCodes {
		ii := i * 5

		accessConfig := dbtypes.NewDbAccessConfig(code.InstantiatePermission)
		cfgValue, _ := accessConfig.Value()

		stmt += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d),", ii+1, ii+2, ii+3, ii+4, ii+5)
		args = append(args, code.Sender, code.WasmByteCode, cfgValue, code.CodeID, code.Height)
	}

	stmt = stmt[:len(stmt)-1] // Remove trailing ","

	stmt += `
	ON CONFLICT (code_id) DO UPDATE 
		SET sender = excluded.sender,
			byte_code = excluded.byte_code,
			instantiate_permission = excluded.instantiate_permission,
			height = excluded.height
	WHERE wasm_code.height <= excluded.height`

	_, err := db.Sql.Exec(stmt, args...)
	if err != nil {
		return fmt.Errorf("error while saving wasm code: %s", err)
	}

	return nil
}

// SaveWasmContracts allows to store the wasm contract slice
func (db *Db) SaveWasmContracts(contracts []types.WasmContract) error {
	paramsNumber := 12
	slices := dbutils.SplitWasmContracts(contracts, paramsNumber)

	for _, contracts := range slices {
		if len(contracts) == 0 {
			continue
		}

		err := db.saveWasmContracts(paramsNumber, contracts)
		if err != nil {
			return fmt.Errorf("error while storing contracts: %s", err)
		}
	}

	return nil
}

func (db *Db) saveWasmContracts(paramsNumber int, wasmContracts []types.WasmContract) error {

	stmt := `
INSERT INTO wasm_contract 
(sender, creator, admin, code_id, label, raw_contract_message, funds, contract_address, data, instantiated_at, contract_info_extension, height) 
VALUES `

	var args []interface{}
	for i, contract := range wasmContracts {
		ii := i * paramsNumber
		stmt += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d),",
			ii+1, ii+2, ii+3, ii+4, ii+5, ii+6, ii+7, ii+8, ii+9, ii+10, ii+11, ii+12)
		args = append(args,
			contract.Sender, contract.Creator, contract.Admin, contract.CodeID, contract.Label, string(contract.RawContractMsg),
			pq.Array(dbtypes.NewDbCoins(contract.Funds)), contract.ContractAddress, contract.Data,
			contract.InstantiatedAt, contract.ContractInfoExtension, contract.Height,
		)
	}

	stmt = stmt[:len(stmt)-1] // Remove trailing ","
	stmt += `
	ON CONFLICT (contract_address) DO UPDATE 
		SET sender = excluded.sender,
			creator = excluded.creator,
			admin = excluded.admin,
			code_id = excluded.code_id,
			label = excluded.label,
			raw_contract_message = excluded.raw_contract_message,
			funds = excluded.funds,
			data = excluded.data,
			instantiated_at = excluded.instantiated_at,
			contract_info_extension = excluded.contract_info_extension,
			height = excluded.height
	WHERE wasm_contract.height <= excluded.height`

	// TO-DO: check if the below is stored as Json in DB:
	// - Data
	// - ContractInfoExtension
	// - RawContractMsg

	_, err := db.Sql.Exec(stmt, args...)
	if err != nil {
		return fmt.Errorf("error while saving wasm contracts: %s", err)
	}

	return nil
}

// SaveWasmExecuteContract allows to store the wasm contract
func (db *Db) SaveWasmExecuteContract(wasmExecuteContract types.WasmExecuteContract) error {
	return db.SaveWasmExecuteContracts([]types.WasmExecuteContract{wasmExecuteContract})
}

// SaveWasmContracts allows to store the wasm contract slice
func (db *Db) SaveWasmExecuteContracts(executeContracts []types.WasmExecuteContract) error {
	paramsNumber := 7
	slices := dbutils.SplitWasmExecuteContracts(executeContracts, paramsNumber)

	for _, contracts := range slices {
		if len(contracts) == 0 {
			continue
		}

		err := db.saveWasmExecuteContracts(paramsNumber, executeContracts)
		if err != nil {
			return fmt.Errorf("error while storing contracts: %s", err)
		}
	}

	return nil
}

func (db *Db) saveWasmExecuteContracts(paramNumber int, executeContracts []types.WasmExecuteContract) error {
	stmt := `
INSERT INTO wasm_execute_contract 
(sender, contract_address, raw_contract_message, funds, data, executed_at, height) 
VALUES ($1, $2, $3, $4, $5, $6, $7) 
ON CONFLICT DO NOTHING`

	var args []interface{}
	for i, executeContract := range executeContracts {
		ii := i * paramNumber
		stmt += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d),",
			ii+1, ii+2, ii+3, ii+4, ii+5, ii+6, ii+7)
		args = append(args,
			executeContract.Sender, executeContract.ContractAddress, string(executeContract.RawContractMsg),
			pq.Array(dbtypes.NewDbCoins(executeContract.Funds)), executeContract.Data, executeContract.ExecutedAt, executeContract.Height)
	}

	stmt = stmt[:len(stmt)-1] // Remove trailing ","

	fmt.Println(stmt)
	fmt.Println(args)

	_, err := db.Sql.Exec(stmt, args...)
	if err != nil {
		return fmt.Errorf("error while saving wasm execute contracts: %s", err)
	}

	return nil
}

func (db *Db) UpdateContractWithMsgMigrateContract(
	sender string, contractAddress string, codeID uint64, rawContractMsg []byte, data string,
) error {

	stmt := `UPDATE wasm_contract SET 
sender = $1, code_id = $2, raw_contract_message = $3, data = $4 
WHERE contract_address = $5 `

	// TO-DO: check if the below is stored as Json in DB:
	// - rawContractMsg
	// - Data

	_, err := db.Sql.Exec(stmt,
		sender, codeID, string(rawContractMsg), data,
		contractAddress,
	)
	if err != nil {
		return fmt.Errorf("error while updating wasm contract from contract migration: %s", err)

	}
	return nil
}

func (db *Db) UpdateContractAdmin(sender string, contractAddress string, newAdmin string) error {

	stmt := `UPDATE wasm_contract SET 
sender = $1, admin = $2 WHERE contract_address = $2 `

	_, err := db.Sql.Exec(stmt, sender, newAdmin, contractAddress)
	if err != nil {
		return fmt.Errorf("error while updating wsm contract admin: %s", err)
	}
	return nil
}
