package database_test

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	gov "github.com/cosmos/cosmos-sdk/x/gov/types"
	dbtypes "github.com/forbole/bdjuno/database/types"
	"github.com/forbole/bdjuno/x/gov/types"
)

func (suite *DbTestSuite) TestBigDipperDb_SaveProposals() {
	proposer1 := suite.getDelegator("cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs")
	submitTime1, err := time.Parse(time.RFC3339, "2020-10-10T15:00:00Z")
	suite.Require().NoError(err)
	depositEndTime1, err := time.Parse(time.RFC3339, "2020-10-15T15:00:00Z")
	suite.Require().NoError(err)
	votingStartTime1, err := time.Parse(time.RFC3339, "2020-10-20T15:00:00Z")
	suite.Require().NoError(err)
	votingEndTime1, err := time.Parse(time.RFC3339, "2020-10-25T15:00:00Z")
	suite.Require().NoError(err)
	status1, err := gov.ProposalStatusFromString("DepositPeriod")
	suite.Require().NoError(err)

	proposer2 := suite.getDelegator("cosmos184ma3twcfjqef6k95ne8w2hk80x2kah7vcwy4a")
	submitTime2, err := time.Parse(time.RFC3339, "2020-12-10T15:00:00Z")
	suite.Require().NoError(err)
	depositEndTime2, err := time.Parse(time.RFC3339, "2020-12-15T15:00:00Z")
	suite.Require().NoError(err)
	votingStartTime2, err := time.Parse(time.RFC3339, "2020-12-20T15:00:00Z")
	suite.Require().NoError(err)
	votingEndTime2, err := time.Parse(time.RFC3339, "2020-12-25T15:00:00Z")
	suite.Require().NoError(err)
	status2, err := gov.ProposalStatusFromString("Passed")
	suite.Require().NoError(err)

	input := []types.Proposal{types.NewProposal("title",
		"description",
		"proposalRoute",
		"proposalType",
		1,
		status1,
		submitTime1,
		depositEndTime1,
		votingStartTime1,
		votingEndTime1,
		proposer1,
	), types.NewProposal("title1",
		"description1",
		"proposalRoute1",
		"proposalType1",
		2,
		status2,
		submitTime2,
		depositEndTime2,
		votingStartTime2,
		votingEndTime2,
		proposer2)}
	err = suite.database.SaveProposals(input)
	suite.Require().NoError(err)

	var proposalRow []dbtypes.ProposalRow
	err = suite.database.Sqlx.Select(&proposalRow, `SELECT * FROM proposal ORDER BY proposal_id ASC`)
	suite.Require().NoError(err)

	expected := []dbtypes.ProposalRow{dbtypes.NewProposalRow("title",
		"description",
		"proposalRoute",
		"proposalType",
		1,
		submitTime1,
		depositEndTime1,
		votingStartTime1,
		votingEndTime1,
		proposer1.String(),
		status1.String(),
	), dbtypes.NewProposalRow("title1",
		"description1",
		"proposalRoute1",
		"proposalType1",
		2,
		submitTime2,
		depositEndTime2,
		votingStartTime2,
		votingEndTime2,
		proposer2.String(),
		status2.String(),
	)}
	for i, expect := range expected {
		suite.Require().True(proposalRow[i].Equals(expect))
	}
}

func (suite *DbTestSuite) TestBigDipperDb_SaveProposal() {
	proposer1 := suite.getDelegator("cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs")
	submitTime1, err := time.Parse(time.RFC3339, "2020-10-10T15:00:00Z")
	suite.Require().NoError(err)
	depositEndTime1, err := time.Parse(time.RFC3339, "2020-10-15T15:00:00Z")
	suite.Require().NoError(err)
	votingStartTime1, err := time.Parse(time.RFC3339, "2020-10-20T15:00:00Z")
	suite.Require().NoError(err)
	votingEndTime1, err := time.Parse(time.RFC3339, "2020-10-25T15:00:00Z")
	suite.Require().NoError(err)
	status1, err := gov.ProposalStatusFromString("DepositPeriod")
	suite.Require().NoError(err)

	input := types.NewProposal("title",
		"description",
		"proposalRoute",
		"proposalType",
		1,
		status1,
		submitTime1,
		depositEndTime1,
		votingStartTime1,
		votingEndTime1,
		proposer1,
	)

	err = suite.database.SaveProposal(input)
	suite.Require().NoError(err)

	var proposalRow []dbtypes.ProposalRow
	err = suite.database.Sqlx.Select(&proposalRow, `SELECT * FROM proposal`)
	suite.Require().NoError(err)
	suite.Require().Len(proposalRow, 1)

	expect := dbtypes.NewProposalRow("title",
		"description",
		"proposalRoute",
		"proposalType",
		1,
		submitTime1,
		depositEndTime1,
		votingStartTime1,
		votingEndTime1,
		proposer1.String(),
		status1.String(),
	)

	suite.Require().True(expect.Equals(proposalRow[0]))
}

func (suite *DbTestSuite) TestBigDipperDb_SaveTallyResult() {
	suite.getProposalRow(1)

	timestamp, err := time.Parse(time.RFC3339, "2020-10-30T15:00:00Z")
	suite.Require().NoError(err)
	tally := types.NewTallyResult(1, 1, 1, 1, 1, 1, timestamp)
	err = suite.database.SaveTallyResult(tally)
	suite.Require().NoError(err)

	expected := dbtypes.NewTallyResultRow(1, 1, 1, 1, 1, 1, timestamp)

	var result []dbtypes.TallyResultRow
	err = suite.database.Sqlx.Select(&result, `SELECT * FROM tally_result`)
	suite.Require().NoError(err)
	suite.Require().Len(result, 1)
	suite.Require().True(expected.Equals(result[0]))
}

func (suite *DbTestSuite) TestBigDipperDb_SaveTallyResults() {
	suite.getProposalRow(1)
	suite.getProposalRow(2)
	timestamp1, err := time.Parse(time.RFC3339, "2020-10-30T15:00:00Z")
	suite.Require().NoError(err)

	timestamp2, err := time.Parse(time.RFC3339, "2020-10-31T15:00:00Z")
	suite.Require().NoError(err)

	input := []types.TallyResult{
		types.NewTallyResult(1, 1, 1, 1, 1, 1, timestamp1),
		types.NewTallyResult(2, 2, 2, 2, 2, 2, timestamp2)}

	err = suite.database.SaveTallyResults(input)
	suite.Require().NoError(err)

	var result []dbtypes.TallyResultRow
	err = suite.database.Sqlx.Select(&result, `SELECT * FROM tally_result`)
	suite.Require().NoError(err)

	expected := []dbtypes.TallyResultRow{
		dbtypes.NewTallyResultRow(1, 1, 1, 1, 1, 1, timestamp1),
		dbtypes.NewTallyResultRow(2, 2, 2, 2, 2, 2, timestamp2)}

	for i, r := range result {
		suite.Require().True(r.Equals(expected[i]))
	}
}

func (suite *DbTestSuite) TestBigDipperDb_SaveVote() {
	proposal := suite.getProposalRow(1)
	voter := suite.getDelegator("cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs")
	timestamp1, err := time.Parse(time.RFC3339, "2020-10-30T15:00:00Z")
	suite.Require().NoError(err)

	vote := types.NewVote(1, voter, gov.OptionYes, 1, timestamp1)
	err = suite.database.SaveVote(vote)
	suite.Require().NoError(err)

	expected := dbtypes.NewVoteRow(int64(proposal.ProposalID), voter.String(), gov.OptionYes.String(), 1, timestamp1)

	var result []dbtypes.VoteRow
	err = suite.database.Sqlx.Select(&result, `SELECT * FROM vote`)
	suite.Require().NoError(err)
	suite.Require().Len(result, 1)
	suite.Require().True(expected.Equals(result[0]))
}

func (suite *DbTestSuite) TestBigDipperDb_SaveDeposit() {
	proposal := suite.getProposalRow(1)
	depositor := suite.getDelegator("cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs")
	amount := sdk.NewCoins(
		sdk.NewCoin("desmos", sdk.NewInt(10000)),
	)
	total := sdk.NewCoins(
		sdk.NewCoin("desmos", sdk.NewInt(20000)),
	)
	timestamp, err := time.Parse(time.RFC3339, "2020-10-30T15:00:00Z")
	suite.Require().NoError(err)

	deposit := types.NewDeposit(proposal.ProposalID, depositor, amount, total, 10, timestamp)
	err = suite.database.SaveDeposit(deposit)
	suite.Require().NoError(err)

	expected := dbtypes.NewDepositRow(1, depositor.String(), dbtypes.NewDbCoins(amount), dbtypes.NewDbCoins(total), 10, timestamp)
	var result []dbtypes.DepositRow
	err = suite.database.Sqlx.Select(&result, `SELECT * FROM deposit`)
	suite.Require().NoError(err)
	suite.Require().Len(result, 1)
	suite.Require().True(expected.Equals(result[0]))
}

func (suite *DbTestSuite) TestBigDipperDb_SaveDeposits() {
	proposal := suite.getProposalRow(1)
	depositor := suite.getDelegator("cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs")
	amount := sdk.NewCoins(
		sdk.NewCoin("desmos", sdk.NewInt(10000)),
	)
	total := sdk.NewCoins(
		sdk.NewCoin("desmos", sdk.NewInt(20000)),
	)
	timestamp, err := time.Parse(time.RFC3339, "2020-10-30T15:00:00Z")
	suite.Require().NoError(err)

	depositor2 := suite.getDelegator("cosmos184ma3twcfjqef6k95ne8w2hk80x2kah7vcwy4a")
	amount2 := sdk.NewCoins(
		sdk.NewCoin("desmos", sdk.NewInt(30000)),
	)
	total2 := sdk.NewCoins(
		sdk.NewCoin("desmos", sdk.NewInt(60000)),
	)

	deposit := []types.Deposit{
		types.NewDeposit(proposal.ProposalID, depositor, amount, total, 10, timestamp),
		types.NewDeposit(proposal.ProposalID, depositor2, amount2, total2, 10, timestamp),
	}

	err = suite.database.SaveDeposits(deposit)
	suite.Require().NoError(err)

	expected := []dbtypes.DepositRow{
		dbtypes.NewDepositRow(1, depositor.String(), dbtypes.NewDbCoins(amount), dbtypes.NewDbCoins(total), 10, timestamp),
		dbtypes.NewDepositRow(1, depositor2.String(), dbtypes.NewDbCoins(amount2), dbtypes.NewDbCoins(total2), 10, timestamp),
	}
	var result []dbtypes.DepositRow
	err = suite.database.Sqlx.Select(&result, `SELECT * FROM deposit`)
	suite.Require().NoError(err)
	for i, r := range result {
		suite.Require().True(expected[i].Equals(r))
	}
}

func (suite *DbTestSuite) TestBigDipperDb_UpdateProposal() {
	proposal := suite.getProposalRow(1)
	proposer, err := sdk.AccAddressFromBech32(proposal.Proposer)
	suite.Require().NoError(err)
	votingStartTime2, err := time.Parse(time.RFC3339, "2020-12-20T15:00:00Z")
	suite.Require().NoError(err)
	votingEndTime2, err := time.Parse(time.RFC3339, "2020-12-25T15:00:00Z")
	suite.Require().NoError(err)
	status2, err := gov.ProposalStatusFromString("Passed")
	suite.Require().NoError(err)
	update := types.NewProposal(proposal.Title,
		proposal.Description,
		proposal.ProposalRoute,
		proposal.ProposalType,
		proposal.ProposalID,
		status2,
		proposal.SubmitTime,
		proposal.DepositEndTime,
		votingStartTime2,
		votingEndTime2,
		proposer)

	err = suite.database.UpdateProposal(update)
	suite.Require().NoError(err)
	expected := dbtypes.NewProposalRow(proposal.Title,
		proposal.Description,
		proposal.ProposalRoute,
		proposal.ProposalType,
		proposal.ProposalID,
		proposal.SubmitTime,
		proposal.DepositEndTime,
		votingStartTime2,
		votingEndTime2,
		proposer.String(),
		status2.String())
	var result []dbtypes.ProposalRow
	err = suite.database.Sqlx.Select(&result, `SELECT * FROM proposal`)
	suite.Require().NoError(err)
	suite.Require().Len(result, 1)
	for _, r := range result {
		suite.Require().True(expected.Equals(r))
	}
}

func (suite *DbTestSuite) getProposalRow(id int) dbtypes.ProposalRow {

	proposer1 := suite.getDelegator("cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs")
	submitTime1, err := time.Parse(time.RFC3339, "2020-10-10T15:00:00Z")
	suite.Require().NoError(err)
	depositEndTime1, err := time.Parse(time.RFC3339, "2020-10-15T15:00:00Z")
	suite.Require().NoError(err)
	votingStartTime1, err := time.Parse(time.RFC3339, "2020-10-20T15:00:00Z")
	suite.Require().NoError(err)
	votingEndTime1, err := time.Parse(time.RFC3339, "2020-10-25T15:00:00Z")
	suite.Require().NoError(err)
	status1, err := gov.ProposalStatusFromString("DepositPeriod")
	suite.Require().NoError(err)

	title := "title" + string(id)
	description := "description" + string(id)
	proposalRoute := "proposalRoute" + string(id)
	proposalType := "proposalType" + string(id)

	proposal := dbtypes.NewProposalRow(title,
		description,
		proposalRoute,
		proposalType,
		uint64(id),
		submitTime1,
		depositEndTime1,
		votingStartTime1,
		votingEndTime1,
		proposer1.String(),
		status1.String(),
	)
	_, err = suite.database.Sqlx.Exec(`INSERT INTO proposal 
	(title, description ,proposer,proposal_route ,proposal_type,proposal_id,
		status,submit_time ,deposit_end_time ,voting_start_time,voting_end_time) VALUES 
		($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`, title,
		description,
		proposer1.String(),
		proposalRoute,
		proposalType,
		uint64(id),
		status1.String(),
		submitTime1,
		depositEndTime1,
		votingStartTime1,
		votingEndTime1)
	suite.Require().NoError(err)

	return proposal
}
