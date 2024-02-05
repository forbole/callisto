package wasm

import (
	"encoding/json"
	"fmt"

	"github.com/forbole/bdjuno/v4/types"
	tmtypes "github.com/tendermint/tendermint/types"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/rs/zerolog/log"
)

// HandleGenesis implements GenesisModule
func (m *Module) HandleGenesis(doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage) error {
	log.Debug().Str("module", "wasm").Msg("parsing genesis")

	// Read the genesis state
	var genState wasmtypes.GenesisState
	err := m.cdc.UnmarshalJSON(appState[wasmtypes.ModuleName], &genState)
	if err != nil {
		return fmt.Errorf("error while unmarshaling wasm genesis state: %s", err)
	}

	err = m.SaveGenesisParams(genState.Params, doc.InitialHeight)
	if err != nil {
		return fmt.Errorf("error while saving genesis wasm params: %s", err)
	}

	err = m.SaveGenesisCodes(genState.Codes, doc.InitialHeight)
	if err != nil {
		return fmt.Errorf("error while saving genesis wasm codes: %s", err)
	}

	err = m.SaveGenesisContracts(genState.Contracts, doc)
	if err != nil {
		return fmt.Errorf("error while saving genesis wasm contracts: %s", err)
	}

	err = m.SaveGenesisMsgs(genState.GenMsgs, doc)
	if err != nil {
		return fmt.Errorf("error while saving genesis wasm execute contracts: %s", err)
	}

	return nil
}

func (m *Module) SaveGenesisParams(params wasmtypes.Params, initHeight int64) error {
	err := m.db.SaveWasmParams(
		types.NewWasmParams(&params.CodeUploadAccess, int32(params.InstantiateDefaultPermission), initHeight),
	)
	if err != nil {
		return fmt.Errorf("error while saving genesis wasm params: %s", err)
	}
	return nil
}

func (m *Module) SaveGenesisCodes(codes []wasmtypes.Code, initHeight int64) error {
	log.Debug().Str("module", "wasm").Str("operation", "genesis codes").
		Int("code counts", len(codes)).Msg("parsing genesis")

	var wasmCodes = []types.WasmCode{}
	for _, code := range codes {
		if code.CodeID != 0 {
			instantiateConfig := code.CodeInfo.InstantiateConfig
			wasmCodes = append(wasmCodes, types.NewWasmCode(
				"", code.CodeBytes, &instantiateConfig, code.CodeID, initHeight,
			))
		}
	}

	if len(wasmCodes) == 0 {
		return nil
	}

	err := m.db.SaveWasmCodes(wasmCodes)
	if err != nil {
		return fmt.Errorf("error while saving genesis wasm codes: %s", err)
	}

	return nil
}

func (m *Module) SaveGenesisContracts(contracts []wasmtypes.Contract, doc *tmtypes.GenesisDoc) error {
	log.Debug().Str("module", "wasm").Str("operation", "genesis contracts").
		Int("contract counts", len(contracts)).Msg("parsing genesis")

	for _, contract := range contracts {

		// Unpack contract info extension
		var contractInfoExt string
		if contract.ContractInfo.Extension != nil {
			var extentionI wasmtypes.ContractInfoExtension
			err := m.cdc.UnpackAny(contract.ContractInfo.Extension, &extentionI)
			if err != nil {
				return fmt.Errorf("error while unpacking genesis contract info extension: %s", err)
			}
			contractInfoExt = extentionI.String()
		}

		// Get contract states
		contractStates, err := m.source.GetContractStates(doc.InitialHeight, contract.ContractAddress)
		if err != nil {
			return fmt.Errorf("error while getting genesis contract states: %s", err)
		}

		contract := types.NewWasmContract(
			"", contract.ContractInfo.Admin, contract.ContractInfo.CodeID, contract.ContractInfo.Label, nil, nil,
			contract.ContractAddress, "", doc.GenesisTime, contract.ContractInfo.Creator, contractInfoExt, contractStates, doc.InitialHeight,
		)

		err = m.db.SaveWasmContracts([]types.WasmContract{contract})
		if err != nil {
			return fmt.Errorf("error while saving genesis wasm contracts: %s", err)
		}
	}

	return nil
}

func (m *Module) SaveGenesisMsgs(msgs []wasmtypes.GenesisState_GenMsgs, doc *tmtypes.GenesisDoc) error {
	log.Debug().Str("module", "wasm").Str("operation", "genesis messages").
		Int("message counts", len(msgs)).Msg("parsing genesis")

	var genesisExecuteContracts = []types.WasmExecuteContract{}
	for _, msg := range msgs {
		// Handle genesis execute contract messages
		if msgExecuteContract, ok := msg.Sum.(*wasmtypes.GenesisState_GenMsgs_ExecuteContract); ok {
			execution := msgExecuteContract.ExecuteContract
			executeContract := types.NewWasmExecuteContract(
				execution.Sender,
				execution.Contract,
				execution.Msg,
				execution.Funds,
				"",
				doc.GenesisTime,
				doc.InitialHeight,
			)
			genesisExecuteContracts = append(genesisExecuteContracts, executeContract)
		}
	}

	if len(genesisExecuteContracts) == 0 {
		return nil
	}

	return m.db.SaveWasmExecuteContracts(genesisExecuteContracts)
}
