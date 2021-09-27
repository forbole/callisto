package modules

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/simapp/params"
	"github.com/desmos-labs/juno/node/remote"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/desmos-labs/juno/node/local"

	jmodules "github.com/desmos-labs/juno/modules"
	"github.com/desmos-labs/juno/modules/messages"
	"github.com/desmos-labs/juno/modules/registrar"

	"github.com/forbole/bdjuno/utils"

	"github.com/forbole/bdjuno/modules/history"

	nodeconfig "github.com/desmos-labs/juno/node/config"
	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/modules/auth"
	"github.com/forbole/bdjuno/modules/bank"
	banksource "github.com/forbole/bdjuno/modules/bank/source"
	localbanksource "github.com/forbole/bdjuno/modules/bank/source/local"
	remotebanksource "github.com/forbole/bdjuno/modules/bank/source/remote"
	"github.com/forbole/bdjuno/modules/consensus"
	"github.com/forbole/bdjuno/modules/distribution"
	"github.com/forbole/bdjuno/modules/gov"
	"github.com/forbole/bdjuno/modules/mint"
	"github.com/forbole/bdjuno/modules/modules"
	"github.com/forbole/bdjuno/modules/pricefeed"
	"github.com/forbole/bdjuno/modules/slashing"
	"github.com/forbole/bdjuno/modules/staking"
)

// UniqueAddressesParser returns a wrapper around the given parser that removes all duplicated addresses
func UniqueAddressesParser(parser messages.MessageAddressesParser) messages.MessageAddressesParser {
	return func(cdc codec.Marshaler, msg sdk.Msg) ([]string, error) {
		addresses, err := parser(cdc, msg)
		if err != nil {
			return nil, err
		}

		return utils.RemoveDuplicateValues(addresses), nil
	}
}

// --------------------------------------------------------------------------------------------------------------------

var (
	_ registrar.Registrar = &Registrar{}
)

// Registrar represents the modules.Registrar that allows to register all modules that are supported by BigDipper
type Registrar struct {
	parser messages.MessageAddressesParser
}

// NewRegistrar allows to build a new Registrar instance
func NewRegistrar(parser messages.MessageAddressesParser) *Registrar {
	return &Registrar{
		parser: UniqueAddressesParser(parser),
	}
}

// BuildModules implements modules.Registrar
func (r *Registrar) BuildModules(ctx registrar.Context) jmodules.Modules {
	cdc := ctx.EncodingConfig.Marshaler
	db := database.Cast(ctx.Database)

	sources, err := buildSources(ctx.JunoConfig.Node, ctx.EncodingConfig)
	if err != nil {
		panic(err)
	}

	return []jmodules.Module{
		messages.NewModule(r.parser, cdc, ctx.Database),
		auth.NewModule(r.parser, cdc, db),
		bank.NewModule(r.parser, sources.BankSource, cdc, db),
		consensus.NewModule(db),

		distribution.NewModule(bdjunoCfg, bankClient, distrClient, db),
		gov.NewModule(bankClient, govClient, stakingClient, encodingConfig, db),
		mint.NewModule(mintClient, db),
		modules.NewModule(ctx.ParsingConfig, db),
		pricefeed.NewModule(bdjunoCfg, encodingConfig, db),
		slashing.NewModule(slashingClient, db),
		staking.NewModule(ctx.ParsingConfig, bankClient, stakingClient, distrClient, encodingConfig, db),
		history.NewModule(r.parser, encodingConfig, db),
	}
}

type Sources struct {
	BankSource banksource.Source
}

func buildSources(nodeCfg nodeconfig.Config, encodingConfig *params.EncodingConfig) (*Sources, error) {
	switch cfg := nodeCfg.Details.(type) {
	case *remote.Details:
		return buildRemoteSources(cfg)
	case *local.Details:
		return buildLocalSources(cfg, encodingConfig)

	default:
		return nil, fmt.Errorf("invalid configuration type: %T", cfg)
	}
}

func buildLocalSources(cfg *local.Details, encodingConfig *params.EncodingConfig) (*Sources, error) {
	source, err := local.NewKeeper(cfg.Home, encodingConfig)
	if err != nil {
		return nil, err
	}

	ak := authkeeper.NewAccountKeeper(source.Codec, source.RegisterKey(authtypes.StoreKey), source.RegisterSubspace(authtypes.ModuleName), authtypes.ProtoBaseAccount, nil)
	bk := bankkeeper.NewBaseKeeper(source.Codec, source.RegisterKey(banktypes.StoreKey), ak, source.RegisterSubspace(banktypes.ModuleName), nil)

	return &Sources{
		BankSource: localbanksource.NewSource(source, bk),
	}, nil
}

func buildRemoteSources(cfg *remote.Details) (*Sources, error) {
	source, err := remote.NewKeeper(cfg.GRPC)
	if err != nil {
		return nil, fmt.Errorf("error while creating remote source: %s", err)
	}

	authClient := authtypes.NewQueryClient(source.GrpcConn)
	bankClient := banktypes.NewQueryClient(source.GrpcConn)
	distrClient := distrtypes.NewQueryClient(source.GrpcConn)
	govClient := govtypes.NewQueryClient(source.GrpcConn)
	mintClient := minttypes.NewQueryClient(source.GrpcConn)
	slashingClient := slashingtypes.NewQueryClient(source.GrpcConn)
	stakingClient := stakingtypes.NewQueryClient(source.GrpcConn)

	return &Sources{
		BankSource: remotebanksource.NewSource(source, bankClient),
	}, nil
}
