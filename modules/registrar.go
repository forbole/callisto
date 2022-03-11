package modules

import (
	"fmt"
	"os"

	"github.com/cosmos/cosmos-sdk/simapp"
	emoneyapp "github.com/e-money/em-ledger"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/forbole/juno/v2/modules/pruning"
	"github.com/forbole/juno/v2/modules/telemetry"

	"github.com/cosmos/cosmos-sdk/simapp/params"
	"github.com/forbole/juno/v2/node/remote"

	"github.com/forbole/bdjuno/v2/modules/authority"
	"github.com/forbole/bdjuno/v2/modules/gov"
	"github.com/forbole/bdjuno/v2/modules/inflation"
	"github.com/forbole/bdjuno/v2/modules/mint"
	"github.com/forbole/bdjuno/v2/modules/slashing"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	authoritytypes "github.com/e-money/em-ledger/x/authority/types"
	inflationtypes "github.com/e-money/em-ledger/x/inflation/types"
	"github.com/forbole/juno/v2/node/local"

	jmodules "github.com/forbole/juno/v2/modules"
	"github.com/forbole/juno/v2/modules/messages"
	"github.com/forbole/juno/v2/modules/registrar"

	"github.com/forbole/bdjuno/v2/utils"

	nodeconfig "github.com/forbole/juno/v2/node/config"

	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/forbole/bdjuno/v2/database"
	"github.com/forbole/bdjuno/v2/modules/auth"
	authoritysource "github.com/forbole/bdjuno/v2/modules/authority/source"
	localauthoritysource "github.com/forbole/bdjuno/v2/modules/authority/source/local"
	remoteauthoritysource "github.com/forbole/bdjuno/v2/modules/authority/source/remote"
	"github.com/forbole/bdjuno/v2/modules/bank"
	banksource "github.com/forbole/bdjuno/v2/modules/bank/source"
	localbanksource "github.com/forbole/bdjuno/v2/modules/bank/source/local"
	remotebanksource "github.com/forbole/bdjuno/v2/modules/bank/source/remote"
	"github.com/forbole/bdjuno/v2/modules/consensus"
	"github.com/forbole/bdjuno/v2/modules/distribution"

	distrsource "github.com/forbole/bdjuno/v2/modules/distribution/source"
	localdistrsource "github.com/forbole/bdjuno/v2/modules/distribution/source/local"
	remotedistrsource "github.com/forbole/bdjuno/v2/modules/distribution/source/remote"

	govsource "github.com/forbole/bdjuno/v2/modules/gov/source"
	remotegovsource "github.com/forbole/bdjuno/v2/modules/gov/source/remote"
	inflationsource "github.com/forbole/bdjuno/v2/modules/inflation/source"
	localinflationsource "github.com/forbole/bdjuno/v2/modules/inflation/source/local"
	remoteinflationsource "github.com/forbole/bdjuno/v2/modules/inflation/source/remote"
	"github.com/forbole/bdjuno/v2/modules/modules"
	"github.com/forbole/bdjuno/v2/modules/pricefeed"
	slashingsource "github.com/forbole/bdjuno/v2/modules/slashing/source"
	localslashingsource "github.com/forbole/bdjuno/v2/modules/slashing/source/local"
	remoteslashingsource "github.com/forbole/bdjuno/v2/modules/slashing/source/remote"
	"github.com/forbole/bdjuno/v2/modules/staking"
	stakingsource "github.com/forbole/bdjuno/v2/modules/staking/source"

	mintsource "github.com/forbole/bdjuno/v2/modules/mint/source"
	remotemintsource "github.com/forbole/bdjuno/v2/modules/mint/source/remote"
	localstakingsource "github.com/forbole/bdjuno/v2/modules/staking/source/local"
	remotestakingsource "github.com/forbole/bdjuno/v2/modules/staking/source/remote"
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

	sources, err := BuildSources(ctx.JunoConfig.Node, ctx.EncodingConfig)
	if err != nil {
		panic(err)
	}

	authModule := auth.NewModule(r.parser, cdc, db)
	authorityModule := authority.NewModule(cdc, sources.AuthoritySource, db)
	bankModule := bank.NewModule(r.parser, sources.BankSource, cdc, db)
	consensusModule := consensus.NewModule(db)
	distrModule := distribution.NewModule(sources.DistrSource, cdc, db)
	inflationModule := inflation.NewModule(cdc, sources.InflationSource, db)
	mintModule := mint.NewModule(sources.MintSource, cdc, db)
	slashingModule := slashing.NewModule(sources.SlashingSource, cdc, db)
	stakingModule := staking.NewModule(sources.StakingSource, slashingModule, cdc, db)
	govModule := gov.NewModule(sources.GovSource, authModule, distrModule, mintModule, slashingModule, stakingModule, cdc, db)

	return []jmodules.Module{
		messages.NewModule(r.parser, cdc, ctx.Database),
		telemetry.NewModule(ctx.JunoConfig),
		pruning.NewModule(ctx.JunoConfig, db, ctx.Logger),

		authModule,
		authorityModule,
		bankModule,
		consensusModule,
		distrModule,
		govModule,
		inflationModule,
		mintModule,
		modules.NewModule(ctx.JunoConfig.Chain, db),
		pricefeed.NewModule(ctx.JunoConfig, cdc, db),
		slashingModule,
		stakingModule,
	}
}

type Sources struct {
	AuthoritySource authoritysource.Source
	BankSource      banksource.Source
	DistrSource     distrsource.Source

	GovSource       govsource.Source
	InflationSource inflationsource.Source
	MintSource      mintsource.Source

	SlashingSource slashingsource.Source
	StakingSource  stakingsource.Source
}

func BuildSources(nodeCfg nodeconfig.Config, encodingConfig *params.EncodingConfig) (*Sources, error) {
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
	source, err := local.NewSource(cfg.Home, encodingConfig)
	if err != nil {
		return nil, err
	}

	emoneyApp := emoneyapp.NewApp(
		log.NewTMLogger(log.NewSyncWriter(os.Stdout)), source.StoreDB, nil, true, map[int64]bool{},
		cfg.Home, 0, emoneyapp.MakeEncodingConfig(), simapp.EmptyAppOptions{},
	)

	sources := &Sources{
		AuthoritySource: localauthoritysource.NewSource(source, authoritytypes.QueryServer(emoneyApp.AuthorityKeeper)),
		BankSource:      localbanksource.NewSource(source, banktypes.QueryServer(emoneyApp.BankKeeper)),
		DistrSource:     localdistrsource.NewSource(source, distrtypes.QueryServer(emoneyApp.DistrKeeper)),
		InflationSource: localinflationsource.NewSource(source, inflationtypes.QueryServer(emoneyApp.InflationKeeper)),
		SlashingSource:  localslashingsource.NewSource(source, slashingtypes.QueryServer(emoneyApp.SlashingKeeper)),
		StakingSource:   localstakingsource.NewSource(source, stakingkeeper.Querier{Keeper: emoneyApp.StakingKeeper}),
	}

	// Mount and initialize the stores
	err = source.MountKVStores(emoneyApp, "keys")
	if err != nil {
		return nil, err
	}

	err = source.MountTransientStores(emoneyApp, "tkeys")
	if err != nil {
		return nil, err
	}

	err = source.MountMemoryStores(emoneyApp, "memKeys")
	if err != nil {
		return nil, err
	}

	err = source.InitStores()
	if err != nil {
		return nil, err
	}

	return sources, nil
}

func buildRemoteSources(cfg *remote.Details) (*Sources, error) {
	source, err := remote.NewSource(cfg.GRPC)
	if err != nil {
		return nil, fmt.Errorf("error while creating remote source: %s", err)
	}

	return &Sources{
		AuthoritySource: remoteauthoritysource.NewSource(source, authoritytypes.NewQueryClient(source.GrpcConn)),
		BankSource:      remotebanksource.NewSource(source, banktypes.NewQueryClient(source.GrpcConn)),
		DistrSource:     remotedistrsource.NewSource(source, distrtypes.NewQueryClient(source.GrpcConn)),
		GovSource:       remotegovsource.NewSource(source, govtypes.NewQueryClient(source.GrpcConn)),
		InflationSource: remoteinflationsource.NewSource(source, inflationtypes.NewQueryClient(source.GrpcConn)),
		MintSource:      remotemintsource.NewSource(source, minttypes.NewQueryClient(source.GrpcConn)),
		SlashingSource:  remoteslashingsource.NewSource(source, slashingtypes.NewQueryClient(source.GrpcConn)),
		StakingSource:   remotestakingsource.NewSource(source, stakingtypes.NewQueryClient(source.GrpcConn)),
	}, nil
}
