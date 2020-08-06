package types

import "time"

// ProposalRow represents a single row inside the proposal table
type ProposalRow struct {
	Title           string    `db:"title"`
	Description     string    `db:"description"`
	ProposalRoute   string    `db:"proposal_route"`
	ProposalType    string    `db:"proposal_type"`
	ProposalID      uint64    `db:"proposal_id"`
	SubmitTime      time.Time `db:"submit_time"`
	DepositEndTime  time.Time `db:"deposit_end_time"`
	VotingStartTime time.Time `db:"voting_start_time"`
	VotingEndTime   time.Time `db:"voting_end_time" `
	Proposer        string    `db:"proposer"`
	Status          string    `db:"status"`
}

// NewCommunityPoolRow allows to easily create a new NewCommunityPoolRow
func NewProposalRow(title string,
	description string,
	proposalRoute string,
	proposalType string,
	proposalID uint64,
	submitTime time.Time,
	depositEndTime time.Time,
	votingStartTime time.Time,
	votingEndTime time.Time,
	proposer string,
	status string) ProposalRow {
	return ProposalRow{
		Title:           title,
		Description:     description,
		ProposalRoute:   proposalRoute,
		ProposalType:    proposalType,
		ProposalID:      proposalID,
		SubmitTime:      submitTime,
		DepositEndTime:  depositEndTime,
		VotingStartTime: votingStartTime,
		VotingEndTime:   votingEndTime,
		Proposer:        proposer,
		Status:          status,
	}
}

// Equals return true if two ProposalRow are the same
func (w ProposalRow) Equals(v ProposalRow) bool {
	return (w.Title == v.Title &&
		w.Description == v.Description &&
		w.ProposalRoute == v.ProposalRoute &&
		w.ProposalType == v.ProposalType &&
		w.ProposalID == v.ProposalID &&
		w.SubmitTime.Equal(v.SubmitTime) &&
		w.DepositEndTime.Equal(v.DepositEndTime) &&
		w.VotingStartTime.Equal(v.VotingStartTime) &&
		w.VotingEndTime.Equal(v.VotingEndTime) &&
		w.Proposer == v.Proposer &&
		w.Status == v.Status)
}

type TallyResultRow struct {
	ProposalID int64     `db:"proposal_id"`
	Yes        int64     `db:"yes"`
	Abstain    int64     `db:"abstain"`
	No         int64     `db:"no"`
	NoWithVeto int64     `db:"no_with_vet"`
	Height     int64     `db:"height"`
	Timestamp  time.Time `db:"timestamp"`
}

func NewTallyResultRow(
	proposalID int64,
	yes int64,
	abstain int64,
	no int64,
	noWithVeto int64,
	height int64,
	timestamp time.Time,
) TallyResultRow {
	return TallyResultRow{
		ProposalID: proposalID,
		Yes:        yes,
		Abstain:    abstain,
		No:         no,
		NoWithVeto: noWithVeto,
		Height:     height,
		Timestamp:  timestamp,
	}
}

func (w TallyResultRow) Equals(v TallyResultRow) bool {
	return w.ProposalID == v.ProposalID &&
		w.Yes == v.Yes &&
		w.Abstain == v.Abstain &&
		w.No == v.No &&
		w.NoWithVeto == v.NoWithVeto &&
		w.Height == v.Height &&
		w.Timestamp.Equal(v.Timestamp)
}

type VoteRow struct {
	ProposalId int64     `db:"proposal_id"`
	Voter      string    `db:"voter"`
	Option     string    `db:"option"`
	Height     int64     `db:"height"`
	Timestamp  time.Time `db:"timestamp"`
}

func NewVoteRow(
	proposalId int64,
	voter string,
	option string,
	height int64,
	timestamp time.Time,
) VoteRow {
	return VoteRow{
		ProposalId: proposalId,
		Voter:      voter,
		Option:     option,
		Height:     height,
		Timestamp:  timestamp,
	}
}

func (w VoteRow) Equals(v VoteRow) bool {
	return w.ProposalId == v.ProposalId &&
		w.Voter == v.Voter &&
		w.Option == v.Option &&
		w.Height == v.Height &&
		w.Timestamp.Equal(v.Timestamp)
}
