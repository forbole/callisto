package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

const (
	ProposalStatusInvalid = "PROPOSAL_STATUS_INVALID"
)

// Proposal represents a single governance proposal
type Proposal struct {
	ProposalRoute   string
	ProposalType    string
	ProposalID      uint64
	Content         govtypes.Content
	Status          string
	SubmitTime      time.Time
	DepositEndTime  time.Time
	VotingStartTime time.Time
	VotingEndTime   time.Time
	Proposer        string
}

// NewProposal return a new Proposal instance
func NewProposal(
	proposalID uint64,
	proposalRoute string,
	proposalType string,
	content govtypes.Content,
	status string,
	submitTime time.Time,
	depositEndTime time.Time,
	votingStartTime time.Time,
	votingEndTime time.Time,
	proposer string,
) Proposal {
	return Proposal{
		Content:         content,
		ProposalRoute:   proposalRoute,
		ProposalType:    proposalType,
		ProposalID:      proposalID,
		Status:          status,
		SubmitTime:      submitTime,
		DepositEndTime:  depositEndTime,
		VotingStartTime: votingStartTime,
		VotingEndTime:   votingEndTime,
		Proposer:        proposer,
	}
}

// Equal tells whether p and other contain the same data
func (p Proposal) Equal(other Proposal) bool {
	return p.ProposalRoute == other.ProposalRoute &&
		p.ProposalType == other.ProposalType &&
		p.ProposalID == other.ProposalID &&
		p.Content.String() == other.Content.String() &&
		p.Status == other.Status &&
		p.SubmitTime.Equal(other.SubmitTime) &&
		p.DepositEndTime.Equal(other.DepositEndTime) &&
		p.VotingStartTime.Equal(other.VotingStartTime) &&
		p.VotingEndTime.Equal(other.VotingEndTime) &&
		p.Proposer == other.Proposer
}

// ProposalUpdate contains the data that should be used when updating a governance proposal
type ProposalUpdate struct {
	ProposalID      uint64
	Status          string
	VotingStartTime time.Time
	VotingEndTime   time.Time
}

// NewProposalUpdate allows to build a new ProposalUpdate instance
func NewProposalUpdate(
	proposalID uint64, status string, votingStartTime, votingEndTime time.Time,
) ProposalUpdate {
	return ProposalUpdate{
		ProposalID:      proposalID,
		Status:          status,
		VotingStartTime: votingStartTime,
		VotingEndTime:   votingEndTime,
	}
}

// -------------------------------------------------------------------------------------------------------------------

// Deposit contains the data of a single deposit made towards a proposal
type Deposit struct {
	ProposalID uint64
	Depositor  string
	Amount     sdk.Coins
	Height     int64
}

// NewDeposit return a new Deposit instance
func NewDeposit(
	proposalID uint64,
	depositor string,
	amount sdk.Coins,
	height int64,
) Deposit {
	return Deposit{
		ProposalID: proposalID,
		Depositor:  depositor,
		Amount:     amount,
		Height:     height,
	}
}

// -------------------------------------------------------------------------------------------------------------------

// Vote contains the data of a single proposal vote
type Vote struct {
	ProposalID uint64
	Voter      string
	Option     govtypes.VoteOption
	Height     int64
}

// NewVote return a new Vote instance
func NewVote(
	proposalID uint64,
	voter string,
	option govtypes.VoteOption,
	height int64,
) Vote {
	return Vote{
		ProposalID: proposalID,
		Voter:      voter,
		Option:     option,
		Height:     height,
	}
}

// -------------------------------------------------------------------------------------------------------------------

// TallyResult contains the data about the final results of a proposal
type TallyResult struct {
	ProposalID uint64
	Yes        int64
	Abstain    int64
	No         int64
	NoWithVeto int64
	Height     int64
}

// NewTallyResult return a new TallyResult instance
func NewTallyResult(
	proposalID uint64,
	yes int64,
	abstain int64,
	no int64,
	noWithVeto int64,
	height int64,
) TallyResult {
	return TallyResult{
		ProposalID: proposalID,
		Yes:        yes,
		Abstain:    abstain,
		No:         no,
		NoWithVeto: noWithVeto,
		Height:     height,
	}
}
