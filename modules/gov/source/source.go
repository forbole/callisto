package source

import govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

type Source interface {
	Proposal(height int64, id uint64) (govtypes.Proposal, error)
	ProposalDeposit(height int64, id uint64, depositor string) (govtypes.Deposit, error)
	TallyResult(height int64, proposalID uint64) (govtypes.TallyResult, error)
	DepositParams(height int64) (govtypes.DepositParams, error)
	VotingParams(height int64) (govtypes.VotingParams, error)
	TallyParams(height int64) (govtypes.TallyParams, error)
}
