package source

import (
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
)

type Source interface {
	GetContractInfo(height int64, contractAddr string) (*wasmtypes.QueryContractInfoResponse, error)
}
