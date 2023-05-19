package types

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/simapp/params"
	"github.com/forbole/juno/v4/node/remote"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	ccvconsumertypes "github.com/cosmos/interchain-security/x/ccv/consumer/types"
	ccvconsumersource "github.com/forbole/bdjuno/v4/modules/ccv/consumer/source"
	remoteccvconsumersource "github.com/forbole/bdjuno/v4/modules/ccv/consumer/source/remote"
	remotewasmsource "github.com/forbole/bdjuno/v4/modules/wasm/source/remote"
	"github.com/forbole/juno/v4/node/local"

	banksource "github.com/forbole/bdjuno/v4/modules/bank/source"
	remotebanksource "github.com/forbole/bdjuno/v4/modules/bank/source/remote"
	slashingsource "github.com/forbole/bdjuno/v4/modules/slashing/source"
	remoteslashingsource "github.com/forbole/bdjuno/v4/modules/slashing/source/remote"
	wasmsource "github.com/forbole/bdjuno/v4/modules/wasm/source"
	nodeconfig "github.com/forbole/juno/v4/node/config"
)

type Sources struct {
	BankSource        banksource.Source
	CcvConsumerSource ccvconsumersource.Source
	SlashingSource    slashingsource.Source
	WasmSource        wasmsource.Source
}

func BuildSources(nodeCfg nodeconfig.Config, encodingConfig *params.EncodingConfig) (*Sources, error) {
	switch cfg := nodeCfg.Details.(type) {
	case *remote.Details:
		return buildRemoteSources(cfg)
	case *local.Details:
		return nil, fmt.Errorf("local source is not supported: %T", cfg)
	default:
		return nil, fmt.Errorf("invalid configuration type: %T", cfg)
	}
}

func buildRemoteSources(cfg *remote.Details) (*Sources, error) {
	source, err := remote.NewSource(cfg.GRPC)
	if err != nil {
		return nil, fmt.Errorf("error while creating remote source: %s", err)
	}

	return &Sources{
		BankSource:        remotebanksource.NewSource(source, banktypes.NewQueryClient(source.GrpcConn)),
		CcvConsumerSource: remoteccvconsumersource.NewSource(source, ccvconsumertypes.NewQueryClient(source.GrpcConn)),
		SlashingSource:    remoteslashingsource.NewSource(source, slashingtypes.NewQueryClient(source.GrpcConn)),
		WasmSource:        remotewasmsource.NewSource(source, wasmtypes.NewQueryClient(source.GrpcConn)),
	}, nil
}
