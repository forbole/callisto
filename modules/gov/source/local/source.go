package local

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/desmos-labs/juno/node/local"
	govsource "github.com/forbole/bdjuno/modules/gov/source"
)

var (
	_ govsource.Source = &Source{}
)

// Source implements govsource.Source by using a local node
type Source struct {
	*local.Source
	k govkeeper.Keeper
}

// NewSource returns a new Source instance
func NewSource(source *local.Source, govKeeper govkeeper.Keeper) *Source {
	return &Source{
		Source: source,
		k:      govKeeper,
	}
}

// Proposal implements govsource.Source
func (s Source) Proposal(height int64, id uint64) (govtypes.Proposal, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return govtypes.Proposal{}, fmt.Errorf("error while loading height: %s", err)
	}

	res, err := s.k.Proposal(sdk.WrapSDKContext(ctx), &govtypes.QueryProposalRequest{ProposalId: id})
	if err != nil {
		return govtypes.Proposal{}, err
	}

	return res.Proposal, nil
}

// ProposalDeposit implements govsource.Source
func (s Source) ProposalDeposit(height int64, id uint64, depositor string) (govtypes.Deposit, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return govtypes.Deposit{}, fmt.Errorf("error while loading height: %s", err)
	}

	res, err := s.k.Deposit(sdk.WrapSDKContext(ctx), &govtypes.QueryDepositRequest{ProposalId: id, Depositor: depositor})
	if err != nil {
		return govtypes.Deposit{}, err
	}

	return res.Deposit, nil
}

// TallyResult implements govsource.Source
func (s Source) TallyResult(height int64, proposalID uint64) (govtypes.TallyResult, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return govtypes.TallyResult{}, fmt.Errorf("error while loading height: %s", err)
	}

	res, err := s.k.TallyResult(sdk.WrapSDKContext(ctx), &govtypes.QueryTallyResultRequest{ProposalId: proposalID})
	if err != nil {
		return govtypes.TallyResult{}, err
	}

	return res.Tally, nil
}

// DepositParams implements govsource.Source
func (s Source) DepositParams(height int64) (govtypes.DepositParams, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return govtypes.DepositParams{}, fmt.Errorf("error while loading height: %s", err)
	}

	res, err := s.k.Params(sdk.WrapSDKContext(ctx), &govtypes.QueryParamsRequest{ParamsType: govtypes.ParamDeposit})
	if err != nil {
		return govtypes.DepositParams{}, err
	}

	return res.DepositParams, nil
}

// VotingParams implements govsource.Source
func (s Source) VotingParams(height int64) (govtypes.VotingParams, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return govtypes.VotingParams{}, fmt.Errorf("error while loading height: %s", err)
	}

	res, err := s.k.Params(sdk.WrapSDKContext(ctx), &govtypes.QueryParamsRequest{ParamsType: govtypes.ParamVoting})
	if err != nil {
		return govtypes.VotingParams{}, err
	}

	return res.VotingParams, nil
}

// TallyParams implements govsource.Source
func (s Source) TallyParams(height int64) (govtypes.TallyParams, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return govtypes.TallyParams{}, fmt.Errorf("error while loading height: %s", err)
	}

	res, err := s.k.Params(sdk.WrapSDKContext(ctx), &govtypes.QueryParamsRequest{ParamsType: govtypes.ParamTallying})
	if err != nil {
		return govtypes.TallyParams{}, err
	}

	return res.TallyParams, nil
}
