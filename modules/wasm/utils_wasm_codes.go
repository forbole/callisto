package wasm

import (
	"fmt"

	"github.com/forbole/callisto/v4/types"
)

func (m *Module) GetWasmCodes(height int64) ([]types.WasmCode, error) {
	// Get contract code list with infos
	codesInfosRes, err := m.source.GetCodesInfos(height)
	if err != nil {
		return nil, fmt.Errorf("error while getting info of codes: %s", err)
	}

	// Traverse code infos to get contract code binary
	var wasmCodes = make([]types.WasmCode, len(codesInfosRes))
	for i, c := range codesInfosRes {
		binary, err := m.source.GetCodeBinary(c.CodeID, height)
		if err != nil {
			return nil, fmt.Errorf("error while getting code binary: %s", err)
		}
		wasmCodes[i] = types.NewWasmCode(c.Creator, binary, &c.InstantiatePermission, c.CodeID, height)
	}

	return wasmCodes, nil
}
