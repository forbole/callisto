package local

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	"github.com/forbole/juno/v5/node/local"

	govsource "github.com/forbole/callisto/v4/modules/gov/source"
)

var (
	_ govsource.Source = &Source{}
)

// Source implements govsource.Source by using a local node
type Source struct {
	*local.Source
	queryClient govtypesv1.QueryServer
}

// NewSource returns a new Source instance
func NewSource(source *local.Source, govKeeper govtypesv1.QueryServer) *Source {
	return &Source{
		Source:      source,
		queryClient: govKeeper,
	}
}

// Proposal implements govsource.Source
func (s Source) Proposal(height int64, id uint64) (*govtypesv1.Proposal, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return nil, fmt.Errorf("error while loading height: %s", err)
	}

	res, err := s.queryClient.Proposal(sdk.WrapSDKContext(ctx), &govtypesv1.QueryProposalRequest{ProposalId: id})
	if err != nil {
		return nil, err
	}

	return res.Proposal, nil
}

// ProposalDeposit implements govsource.Source
func (s Source) ProposalDeposit(height int64, id uint64, depositor string) (*govtypesv1.Deposit, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return nil, fmt.Errorf("error while loading height: %s", err)
	}

	res, err := s.queryClient.Deposit(sdk.WrapSDKContext(ctx), &govtypesv1.QueryDepositRequest{ProposalId: id, Depositor: depositor})
	if err != nil {
		return nil, err
	}

	return res.Deposit, nil
}

// TallyResult implements govsource.Source
func (s Source) TallyResult(height int64, proposalID uint64) (*govtypesv1.TallyResult, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return nil, fmt.Errorf("error while loading height: %s", err)
	}

	res, err := s.queryClient.TallyResult(sdk.WrapSDKContext(ctx), &govtypesv1.QueryTallyResultRequest{ProposalId: proposalID})
	if err != nil {
		return nil, err
	}

	return res.Tally, nil
}

// Params implements govsource.Source
func (s Source) Params(height int64) (*govtypesv1.Params, error) {
	ctx, err := s.LoadHeight(height)
	if err != nil {
		return nil, fmt.Errorf("error while loading height: %s", err)
	}

	res, err := s.queryClient.Params(sdk.WrapSDKContext(ctx), &govtypesv1.QueryParamsRequest{ParamsType: govtypesv1.ParamDeposit})
	if err != nil {
		return nil, err
	}

	return res.Params, nil
}
