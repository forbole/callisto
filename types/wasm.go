package types

import (
	"encoding/hex"
	"encoding/json"
	"time"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// WasmParams represents the CosmWasm code in x/wasm module
type WasmParams struct {
	Params wasmtypes.Params
	Height int64
}

// NewWasmParams allows to build a new x/wasm params instance
func NewWasmParams(
	params wasmtypes.Params, height int64,
) WasmParams {
	return WasmParams{
		Params: params,
		Height: height,
	}
}

// WasmCode represents the CosmWasm code in x/wasm module
type WasmCode struct {
	Sender                string
	WasmByteCode          []byte
	InstantiatePermission *wasmtypes.AccessConfig
	CodeID                uint64
	Height                int64
}

// NewWasmCode allows to build a new x/wasm code instance
func NewWasmCode(
	sender string, wasmByteCode []byte, iPermission *wasmtypes.AccessConfig, codeID uint64, height int64,
) WasmCode {
	return WasmCode{
		Sender:                sender,
		WasmByteCode:          wasmByteCode,
		InstantiatePermission: iPermission,
		CodeID:                codeID,
		Height:                height,
	}
}

// WasmContract represents the CosmWasm contract in x/wasm module
type WasmContract struct {
	Sender                string
	Admin                 string
	CodeID                uint64
	Label                 string
	RawContractMsg        wasmtypes.RawContractMessage
	Funds                 sdk.Coins
	ContractAddress       string
	Data                  string
	InstantiatedAt        time.Time
	Creator               string
	ContractInfoExtension string
	ContractStates        []byte
	Height                int64
}

// NewWasmCode allows to build a new x/wasm contract instance
func NewWasmContract(
	sender string, admin string, codeID uint64, label string, rawMsg wasmtypes.RawContractMessage, funds sdk.Coins, contractAddress string, data string,
	instantiatedAt time.Time, creator string, contractInfoExtension string, states []wasmtypes.Model, height int64,
) WasmContract {
	return WasmContract{
		Sender:                sender,
		Admin:                 admin,
		CodeID:                codeID,
		Label:                 label,
		RawContractMsg:        ConvertRawContractMessage(rawMsg),
		Funds:                 funds,
		ContractAddress:       contractAddress,
		Data:                  data,
		InstantiatedAt:        instantiatedAt,
		Creator:               creator,
		ContractInfoExtension: contractInfoExtension,
		ContractStates:        ConvertContractStates(states),
		Height:                height,
	}
}

// ConvertContractStates removes unaccepted hex characters for postgreSQL from the state key
func ConvertContractStates(states []wasmtypes.Model) []byte {
	var jsonStates = make(map[string]interface{})

	hexZero, _ := hex.DecodeString("00")
	for _, state := range states {
		key := state.Key
		// Remove initial 2 hex characters if the first is \x00
		if string(state.Key[:1]) == string(hexZero) {
			key = state.Key[2:]
		}

		// Remove \x00 hex characters in the middle
		for i := 0; i < len(key); i++ {
			if string(key[i]) == string(hexZero) {
				key = append(key[:i], key[i+1:]...)
				i--
			}
		}

		// Decode hex value
		keyBz, _ := hex.DecodeString(key.String())

		jsonStates[string(keyBz)] = string(state.Value)
	}
	jsonStatesBz, _ := json.Marshal(&jsonStates)

	return jsonStatesBz
}

// WasmExecuteContract represents the CosmWasm execute contract in x/wasm module
type WasmExecuteContract struct {
	Sender          string
	ContractAddress string
	RawContractMsg  []byte
	Funds           sdk.Coins
	Data            string
	ExecutedAt      time.Time
	Height          int64
}

// NewWasmExecuteContract allows to build a new x/wasm execute contract instance
func NewWasmExecuteContract(
	sender string, contractAddress string, rawMsg wasmtypes.RawContractMessage,
	funds sdk.Coins, data string, executedAt time.Time, height int64,
) WasmExecuteContract {

	return WasmExecuteContract{
		Sender:          sender,
		ContractAddress: contractAddress,
		RawContractMsg:  ConvertRawContractMessage(rawMsg),
		Funds:           funds,
		Data:            data,
		ExecutedAt:      executedAt,
		Height:          height,
	}
}

func ConvertRawContractMessage(rawMsg wasmtypes.RawContractMessage) []byte {
	rawContractMsg, _ := rawMsg.MarshalJSON()
	if len(rawMsg) == 0 {
		rawContractMsg, _ = json.Marshal("")
	}

	return rawContractMsg
}
