package remote

import (
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/desmos-labs/juno/node/remote"
	govsource "github.com/forbole/bdjuno/modules/gov/source"
)

var (
	_ govsource.Source = &Source{}
)

// Source implements govsource.Source using a remote node
type Source struct {
	*remote.Source
	govClient govtypes.QueryClient
}

// NewSource returns a new Source implementation
func NewSource(source *remote.Source, govClient govtypes.QueryClient) *Source {
	return &Source{
		Source:    source,
		govClient: govClient,
	}
}

// Proposal implements govsource.Source
func (s Source) Proposal(height int64, id uint64) (govtypes.Proposal, error) {
	res, err := s.govClient.Proposal(
		s.Ctx,
		&govtypes.QueryProposalRequest{ProposalId: id},
		remote.GetHeightRequestHeader(height),
	)
	if err != nil {
		return govtypes.Proposal{}, err
	}

	return res.Proposal, err
}

// ProposalDeposit implements govsource.Source
func (s Source) ProposalDeposit(height int64, id uint64, depositor string) (govtypes.Deposit, error) {
	res, err := s.govClient.Deposit(
		s.Ctx,
		&govtypes.QueryDepositRequest{ProposalId: id, Depositor: depositor},
		remote.GetHeightRequestHeader(height),
	)
	if err != nil {
		return govtypes.Deposit{}, err
	}

	return res.Deposit, nil
}

// TallyResult implements govsource.Source
func (s Source) TallyResult(height int64, proposalID uint64) (govtypes.TallyResult, error) {
	res, err := s.govClient.TallyResult(
		s.Ctx,
		&govtypes.QueryTallyResultRequest{ProposalId: proposalID},
		remote.GetHeightRequestHeader(height),
	)
	if err != nil {
		return govtypes.TallyResult{}, err
	}

	return res.Tally, nil
}

// DepositParams implements govsource.Source
func (s Source) DepositParams(height int64) (govtypes.DepositParams, error) {
	res, err := s.govClient.Params(
		s.Ctx,
		&govtypes.QueryParamsRequest{ParamsType: govtypes.ParamDeposit},
		remote.GetHeightRequestHeader(height),
	)
	if err != nil {
		return govtypes.DepositParams{}, err
	}

	return res.DepositParams, nil
}

// VotingParams implements govsource.Source
func (s Source) VotingParams(height int64) (govtypes.VotingParams, error) {
	res, err := s.govClient.Params(
		s.Ctx,
		&govtypes.QueryParamsRequest{ParamsType: govtypes.ParamVoting},
		remote.GetHeightRequestHeader(height),
	)
	if err != nil {
		return govtypes.VotingParams{}, err
	}

	return res.VotingParams, nil
}

// TallyParams implements govsource.Source
func (s Source) TallyParams(height int64) (govtypes.TallyParams, error) {
	res, err := s.govClient.Params(
		s.Ctx,
		&govtypes.QueryParamsRequest{ParamsType: govtypes.ParamTallying},
		remote.GetHeightRequestHeader(height),
	)
	if err != nil {
		return govtypes.TallyParams{}, err
	}

	return res.TallyParams, nil
}
