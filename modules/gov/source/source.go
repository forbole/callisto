package source

import (
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govtypesv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

type Source interface {
	Proposal(height int64, id uint64) (govtypesv1beta1.Proposal, error)
	ProposalDeposit(height int64, id uint64, depositor string) (*govtypes.Deposit, error)
	TallyResult(height int64, proposalID uint64) (*govtypes.TallyResult, error)
	DepositParams(height int64) (*govtypes.DepositParams, error)
	VotingParams(height int64) (*govtypes.VotingParams, error)
	TallyParams(height int64) (*govtypes.TallyParams, error)
}
