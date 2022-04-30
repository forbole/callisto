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

	// TO-DO: check if string(wasmCode.WasmByteCode) works

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
