package types

import (
	"strconv"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

const (
	ProposalStatusInvalid = "PROPOSAL_STATUS_INVALID"
)

type DepositParams struct {
	MinDeposit       sdk.Coins `json:"min_deposit,omitempty" yaml:"min_deposit"`
	MaxDepositPeriod string    `json:"max_deposit_period,omitempty" yaml:"max_deposit_period"`
}

func NewdepositParam(d govtypes.DepositParams) DepositParams {
	return DepositParams{
		MinDeposit:       d.MinDeposit,
		MaxDepositPeriod: strconv.FormatInt(d.MaxDepositPeriod.Nanoseconds(), 10) + "n",
	}
}

type VotingParams struct {
	VotingPeriod string `json:"voting_period,omitempty" yaml:"voting_period"`
}

func NewvotingParams(v govtypes.VotingParams) VotingParams {
	return VotingParams{
		VotingPeriod: strconv.FormatInt(v.VotingPeriod.Nanoseconds(), 10) + "n",
	}
}

type GovParams struct {
	DepositParams DepositParams        `json:"deposit_params" yaml:"deposit_params"`
	VotingParams  VotingParams         `json:"voting_params" yaml:"voting_params"`
	TallyParams   govtypes.TallyParams `json:"tally_params" yaml:"tally_params"`
	Height        int64                `json:"height" ymal:"height"`
}

func NewGovParams(v VotingParams, d DepositParams, t govtypes.TallyParams, height int64) *GovParams {
	return &GovParams{
		DepositParams: d,
		VotingParams:  v,
		TallyParams:   t,
		Height:        height,
	}
}

func (*VotingParams) ProtoMessage()  {}
func (*DepositParams) ProtoMessage() {}
func (m *VotingParams) Reset()       { *m = VotingParams{} }
func (m *DepositParams) Reset()      { *m = DepositParams{} }
func (m VotingParams) String() string {
	return m.VotingPeriod
}
func (m DepositParams) String() string {
	return m.MaxDepositPeriod
}

// --------------------------------------------------------------------------------------------------------------------

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

// -------------------------------------------------------------------------------------------------------------------

// ProposalStakingPoolSnapshot contains the data about a single staking pool snapshot to be associated with a proposal
type ProposalStakingPoolSnapshot struct {
	ProposalID uint64
	Pool       *Pool
}

// NewProposalStakingPoolSnapshot returns a new ProposalStakingPoolSnapshot instance
func NewProposalStakingPoolSnapshot(proposalID uint64, pool *Pool) ProposalStakingPoolSnapshot {
	return ProposalStakingPoolSnapshot{
		ProposalID: proposalID,
		Pool:       pool,
	}
}

// -------------------------------------------------------------------------------------------------------------------

// ProposalValidatorStatusSnapshot represents a single snapshot of the status of a validator associated
// with a single proposal
type ProposalValidatorStatusSnapshot struct {
	ProposalID           uint64
	ValidatorConsAddress string
	ValidatorVotingPower int64
	ValidatorStatus      int
	ValidatorJailed      bool
	Height               int64
}

// NewProposalValidatorStatusSnapshot returns a new ProposalValidatorStatusSnapshot instance
func NewProposalValidatorStatusSnapshot(
	proposalID uint64,
	validatorConsAddr string,
	validatorVotingPower int64,
	validatorStatus int,
	validatorJailed bool,
	height int64,
) ProposalValidatorStatusSnapshot {
	return ProposalValidatorStatusSnapshot{
		ProposalID:           proposalID,
		ValidatorStatus:      validatorStatus,
		ValidatorConsAddress: validatorConsAddr,
		ValidatorVotingPower: validatorVotingPower,
		ValidatorJailed:      validatorJailed,
		Height:               height,
	}
}
