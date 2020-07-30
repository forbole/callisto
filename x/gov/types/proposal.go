package types

import (
	"math/big"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
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
	Yes        big.Int
	Abstain    big.Int
	No         big.Int
	NoWithVeto big.Int
	Height     int64
	Timestamp  time.Time
}

// NewTallyResult return a new TallyResult instance
func NewTallyResult(
	proposalID uint64,
	yes big.Int,
	abstain big.Int,
	no big.Int,
	noWithVeto big.Int,
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
