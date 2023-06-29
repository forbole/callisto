package remote

import (
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	"github.com/forbole/juno/v5/node/remote"

	govsource "github.com/forbole/bdjuno/v4/modules/gov/source"
)

var (
	_ govsource.Source = &Source{}
)

// Source implements govsource.Source using a remote node
type Source struct {
	*remote.Source
	queryClient govtypesv1.QueryClient
}

// NewSource returns a new Source implementation
func NewSource(source *remote.Source, queryClient govtypesv1.QueryClient) *Source {
	return &Source{
		Source:      source,
		queryClient: queryClient,
	}
}

// Proposal implements govsource.Source
func (s Source) Proposal(height int64, id uint64) (*govtypesv1.Proposal, error) {
	res, err := s.queryClient.Proposal(
		remote.GetHeightRequestContext(s.Ctx, height),
		&govtypesv1.QueryProposalRequest{ProposalId: id},
	)
	if err != nil {
		return nil, err
	}

	return res.Proposal, err
}

// ProposalDeposit implements govsource.Source
func (s Source) ProposalDeposit(height int64, id uint64, depositor string) (*govtypesv1.Deposit, error) {
	res, err := s.queryClient.Deposit(
		remote.GetHeightRequestContext(s.Ctx, height),
		&govtypesv1.QueryDepositRequest{ProposalId: id, Depositor: depositor},
	)
	if err != nil {
		return nil, err
	}

	return res.Deposit, nil
}

// TallyResult implements govsource.Source
func (s Source) TallyResult(height int64, proposalID uint64) (*govtypesv1.TallyResult, error) {
	res, err := s.queryClient.TallyResult(
		remote.GetHeightRequestContext(s.Ctx, height),
		&govtypesv1.QueryTallyResultRequest{ProposalId: proposalID},
	)
	if err != nil {
		return nil, err
	}

	return res.Tally, nil
}

// Params implements govsource.Source
func (s Source) Params(height int64) (*govtypesv1.Params, error) {
	res, err := s.queryClient.Params(
		remote.GetHeightRequestContext(s.Ctx, height),
		&govtypesv1.QueryParamsRequest{ParamsType: govtypesv1.ParamDeposit},
	)
	if err != nil {
		return nil, err
	}

	return res.Params, nil
}
