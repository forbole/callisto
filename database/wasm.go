package database

import (
	"encoding/json"
	"fmt"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	dbtypes "github.com/forbole/callisto/v4/database/types"
	dbutils "github.com/forbole/callisto/v4/database/utils"
	"github.com/forbole/callisto/v4/types"
	"github.com/lib/pq"
)

// SaveWasmParams allows to store the wasm params
func (db *Db) SaveWasmParams(params types.WasmParams) error {
	paramsBz, err := json.Marshal(&params.Params)
	if err != nil {
		return fmt.Errorf("error while marshaling wasm params: %s", err)
	}

	stmt := `
INSERT INTO wasm_params(params, height) 
VALUES ($1, $2) 
ON CONFLICT (one_row_id) DO UPDATE 
	SET params = excluded.params 
WHERE wasm_params.height <= excluded.height
`

	_, err = db.SQL.Exec(stmt, string(paramsBz), params.Height)

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

	if len(wasmCodes) == 0 {
		return fmt.Errorf("wasm codes list is empty")
	}

	var args []interface{}
	var accounts = make([]types.Account, len(wasmCodes))
	for i, code := range wasmCodes {
		ii := i * 5

		var permissionBz []byte
		var err error
		if code.InstantiatePermission != nil {
			permissionBz, err = json.Marshal(code.InstantiatePermission)
			if err != nil {
				return fmt.Errorf("error while marshaling wasm instantiate permission: %s", err)
			}
		}

		stmt += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d),", ii+1, ii+2, ii+3, ii+4, ii+5)
		args = append(args, code.Sender, code.WasmByteCode, string(permissionBz), code.CodeID, code.Height)
		accounts[i] = types.NewAccount(code.Sender)
	}

	stmt = stmt[:len(stmt)-1] // Remove trailing ","

	stmt += `
	ON CONFLICT (code_id) DO UPDATE 
		SET sender = excluded.sender,
			byte_code = excluded.byte_code,
			instantiate_permission = excluded.instantiate_permission,
			height = excluded.height
	WHERE wasm_code.height <= excluded.height`

	err := db.SaveAccounts(accounts)
	if err != nil {
		return fmt.Errorf("error while saving accounts: %s", err)
	}

	_, err = db.SQL.Exec(stmt, args...)
	if err != nil {
		return fmt.Errorf("error while saving wasm code: %s", err)
	}

	return nil
}

// SaveWasmContracts allows to store the wasm contract slice
func (db *Db) SaveWasmContracts(contracts []types.WasmContract) error {
	paramsNumber := 13
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
(sender, creator, admin, code_id, label, raw_contract_message, funds, contract_address, 
data, instantiated_at, contract_info_extension, contract_states, height) 
VALUES `

	var args []interface{}
	var accounts = make([]types.Account, len(wasmContracts))
	for i, contract := range wasmContracts {
		ii := i * paramsNumber
		stmt += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d),",
			ii+1, ii+2, ii+3, ii+4, ii+5, ii+6, ii+7, ii+8, ii+9, ii+10, ii+11, ii+12, ii+13)
		args = append(args,
			contract.Sender,
			contract.Creator,
			contract.Admin,
			contract.CodeID,
			contract.Label,
			string(contract.RawContractMsg),
			pq.Array(dbtypes.NewDbCoins(contract.Funds)),
			contract.ContractAddress,
			contract.Data,
			contract.InstantiatedAt,
			contract.ContractInfoExtension,
			string(contract.ContractStates),
			contract.Height,
		)

		accounts[i] = types.NewAccount(contract.Creator)
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
			contract_states = excluded.contract_states,
			height = excluded.height
	WHERE wasm_contract.height <= excluded.height`

	err := db.SaveAccounts(accounts)
	if err != nil {
		return fmt.Errorf("error while saving accounts: %s", err)
	}

	_, err = db.SQL.Exec(stmt, args...)
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
VALUES `

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

	stmt += ` ON CONFLICT DO NOTHING`

	_, err := db.SQL.Exec(stmt, args...)
	if err != nil {
		return fmt.Errorf("error while saving wasm execute contracts: %s", err)
	}

	return nil
}

func (db *Db) UpdateContractWithMsgMigrateContract(
	sender string, contractAddress string, codeID uint64, rawContractMsg wasmtypes.RawContractMessage, data string,
) error {
	stmt := `UPDATE wasm_contract SET 
sender = $1, code_id = $2, raw_contract_message = $3, data = $4 
WHERE contract_address = $5 `

	converted := types.ConvertRawContractMessage(rawContractMsg)

	_, err := db.SQL.Exec(stmt,
		sender,
		codeID,
		string(converted),
		data,
		contractAddress,
	)
	if err != nil {
		return fmt.Errorf("error while updating wasm contract from contract migration: %s", err)

	}
	return nil
}

func (db *Db) UpdateContractAdmin(sender string, contractAddress string, newAdmin string) error {

	stmt := `UPDATE wasm_contract SET 
sender = $1, admin = $2 WHERE contract_address = $3 `

	_, err := db.SQL.Exec(stmt, sender, newAdmin, contractAddress)
	if err != nil {
		return fmt.Errorf("error while updating wsm contract admin: %s", err)
	}
	return nil
}
