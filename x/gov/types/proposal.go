package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	gov "github.com/cosmos/cosmos-sdk/x/gov/types"
)

// Proposal represent storing a gov.proposal
// For final tolly result, it stored in tally result as they share same proposal ID and VotingEndTime
type Proposal struct {
	Title           string
	Description     string
	ProposalRoute   string
	ProposalType    string
	ProposalID      uint64
	Status          string
	SubmitTime      time.Time
	DepositEndTime  time.Time
	TotalDeposit    sdk.Coins
	VotingStartTime time.Time
	VotingEndTime   time.Time
}

// NewProposal return a new Proposal instance
func NewProposal(
	title string,
	description string,
	proposalRoute string,
	proposalType string,
	proposalID uint64,
	status string,
	submitTime time.Time,
	depositEndTime time.Time,
	totalDeposit sdk.Coins,
	votingStartTime time.Time,
	votingEndTime time.Time,
) Proposal {
	return Proposal{
		Title:           title,
		Description:     description,
		ProposalRoute:   proposalRoute,
		ProposalType:    proposalType,
		ProposalID:      proposalID,
		Status:          status, //ProposalStatusFromString()
		SubmitTime:      submitTime,
		DepositEndTime:  depositEndTime,
		TotalDeposit:    totalDeposit,
		VotingStartTime: votingStartTime,
		VotingEndTime:   votingEndTime,
	}
}

//MsgVote
type TallyResult struct {
	ProposalID uint64
	Yes        int64
	Abstain    int64
	No         int64
	NoWithVeto int64
	Height     int64
	Timestamp  time.Time
}

// NewTallyResult return a new TallyResult instance
func NewTallyResult(
	proposalID uint64,
	yes int64,
	abstain int64,
	no int64,
	noWithVeto int64,
	height int64,
	timestamp time.Time,
) TallyResult {
	return TallyResult{
		ProposalID: proposalID,
		Yes:        yes,
		Abstain:    abstain,
		No:         no,
		NoWithVeto: noWithVeto,
		Height:     height,
		Timestamp:  timestamp,
	}
}

// Vote discribe a msgVote
type Vote struct {
	ProposalID uint64
	Voter      sdk.AccAddress
	Option     gov.VoteOption
	Height     int64
	Timestamp  time.Time
}

// NewVote return a new Vote instance
func NewVote(
	proposalID uint64,
	voter sdk.AccAddress,
	option gov.VoteOption,
	height int64,
	timestamp time.Time,
) Vote {
	return Vote{
		ProposalID: proposalID,
		Voter:      voter,
		Option:     option,
		Height:     height,
		Timestamp:  timestamp,
	}
}

//Deposit represent a message that a user do deposit action
type Deposit struct {
	ProposalID uint64
	Depositor  sdk.AccAddress
	Amount     sdk.Coins
	Height     int64
	Timestamp  time.Time
}

//NewDeposit return a new Deposit instance
func NewDeposit(
	proposalID uint64,
	depositor sdk.AccAddress,
	amount sdk.Coins,
	height int64,
	timestamp time.Time,
) Deposit {
	return Deposit{
		ProposalID: proposalID,
		Depositor:  depositor,
		Amount:     amount,
		Height:     height,
		Timestamp:  timestamp,
	}
}
