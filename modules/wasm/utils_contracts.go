package wasm

import (
	"fmt"
	"strings"
	"time"

	"github.com/forbole/bdjuno/v4/types"
	"github.com/rs/zerolog/log"
)

// StoreContracts gets the available contracts and stores them inside the database
func (m *Module) StoreContracts(height int64) error {
	log.Debug().Str("module", "wasm").Int64("height", height).
		Msg("storing x/wasm contracts")

	codes, err := m.getWasmCodes(height)
	if err != nil {
		return fmt.Errorf("error while handling contracts codes: %s", err)
	}

	contracts, err := m.getContractsByCode(codes, height)
	if err != nil {
		return fmt.Errorf("error while handling contracts codes: %s", err)
	}

	var wasmContracts []types.WasmContract
	for _, contract := range contracts {
		contractStates, err := m.source.GetContractStates(height, contract)
		if err != nil {
			return fmt.Errorf("error while getting contracts states: %s", err)
		}

		contractInfo, err := m.source.GetContractInfo(height, contract)
		if err != nil {
			return fmt.Errorf("error while getting contracts info: %s", err)
		}

		wasmContracts = append(wasmContracts, types.NewWasmContract(
			"", contractInfo.ContractInfo.Admin, contractInfo.ContractInfo.CodeID, contractInfo.ContractInfo.Label, nil, nil,
			contract, "", time.Now(), contractInfo.ContractInfo.Creator, contractInfo.ContractInfo.Extension.GoString(), contractStates, height,
		))

	}

	err = m.db.SaveWasmContracts(wasmContracts)
	if err != nil {
		return fmt.Errorf("error while saving wasm contracts: %s", err)
	}

	return nil
}

func (m *Module) getWasmCodes(height int64) ([]types.WasmCode, error) {
	var wasmCodes = []types.WasmCode{}
	codes, err := m.source.GetCodes(height)
	if err != nil {
		return nil, fmt.Errorf("error while getting contracts codes: %s", err)
	}

	for _, c := range codes {
		instantiatePermission := c.InstantiatePermission
		wasmCodes = append(wasmCodes, types.NewWasmCode(
			"", c.DataHash, &instantiatePermission, c.CodeID, height,
		))
	}

	if len(wasmCodes) == 0 {
		return nil, nil
	}

	err = m.db.SaveWasmCodes(wasmCodes)
	if err != nil {
		return nil, fmt.Errorf("error while saving wasm codes: %s", err)
	}

	return wasmCodes, nil
}

func (m *Module) getContractsByCode(codes []types.WasmCode, height int64) ([]string, error) {
	var contracts []string
	for _, code := range codes {
		contract, err := m.source.GetContractsByCode(height, code.CodeID)
		if err != nil {
			return nil, fmt.Errorf("error while getting contracts codes: %s", err)
		}
		for _, d := range contract {
			values := strings.Split(d, " ")
			contracts = append(contracts, values...)
		}
	}

	return contracts, nil
}
