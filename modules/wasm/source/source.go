package source

import (
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
)

type Source interface {
	GetCodesInfos(height int64) ([]wasmtypes.CodeInfoResponse, error)
	GetCodeBinary(codeID uint64, height int64) ([]byte, error)
	GetContractInfo(height int64, contractAddr string) (*wasmtypes.QueryContractInfoResponse, error)
	GetContractStates(height int64, contractAddress string) ([]wasmtypes.Model, error)
}
