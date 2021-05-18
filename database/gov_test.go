package database_test

import (
	"fmt"
	"time"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/gogo/protobuf/proto"

	"github.com/forbole/bdjuno/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	dbtypes "github.com/forbole/bdjuno/database/types"
)

func (suite *DbTestSuite) getProposalRow(id int) types.Proposal {
	proposer := suite.getAccount("cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs")

	title := fmt.Sprintf("title%d", id)
	description := fmt.Sprintf("description%d", id)
	proposalRoute := fmt.Sprintf("proposalRoute%d", id)
	proposalType := fmt.Sprintf("proposalType%d", id)

	proposal := types.NewProposal(
		uint64(id),
		proposalRoute,
		proposalType,
		govtypes.NewTextProposal(title, description),
		govtypes.StatusPassed,
		time.Date(2020, 1, 1, 00, 00, 00, 000, time.UTC),
		time.Date(2020, 1, 1, 01, 00, 00, 000, time.UTC),
		time.Date(2020, 1, 1, 02, 00, 00, 000, time.UTC),
		time.Date(2020, 1, 1, 03, 00, 00, 000, time.UTC),
		proposer.String(),
	)

	err := suite.database.SaveProposals([]types.Proposal{proposal})
	suite.Require().NoError(err)

	return proposal
}

func (suite *DbTestSuite) encodeProposalContent(content govtypes.Content) string {
	protoContent, ok := content.(proto.Message)
	suite.Require().True(ok)

	anyContent, err := codectypes.NewAnyWithValue(protoContent)
	suite.Require().NoError(err)

	contentBz, err := suite.database.EncodingConfig.Marshaler.MarshalJSON(anyContent)
	suite.Require().NoError(err)

	return string(contentBz)
}

func (suite *DbTestSuite) TestBigDipperDb_SaveProposals() {
	proposer1 := suite.getAccount("cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs")
	proposer2 := suite.getAccount("cosmos184ma3twcfjqef6k95ne8w2hk80x2kah7vcwy4a")

	content1 := govtypes.NewTextProposal("title", "description")
	content2 := govtypes.NewTextProposal("title1", "description1")
	input := []types.Proposal{
		types.NewProposal(
			1,
			"proposalRoute",
			"proposalType",
			content1,
			govtypes.StatusDepositPeriod,
			time.Date(2020, 1, 1, 00, 00, 00, 000, time.UTC),
			time.Date(2020, 1, 1, 01, 00, 00, 000, time.UTC),
			time.Date(2020, 1, 1, 02, 00, 00, 000, time.UTC),
			time.Date(2020, 1, 1, 03, 00, 00, 000, time.UTC),
			proposer1.String(),
		),
		types.NewProposal(
			2,
			"proposalRoute1",
			"proposalType1",
			content2,
			govtypes.StatusPassed,
			time.Date(2020, 1, 2, 00, 00, 00, 000, time.UTC),
			time.Date(2020, 1, 2, 01, 00, 00, 000, time.UTC),
			time.Date(2020, 1, 2, 02, 00, 00, 000, time.UTC),
			time.Date(2020, 1, 2, 03, 00, 00, 000, time.UTC),
			proposer2.String(),
		),
	}

	err := suite.database.SaveProposals(input)
	suite.Require().NoError(err)

	var proposalRow []dbtypes.ProposalRow
	err = suite.database.Sqlx.Select(&proposalRow, `SELECT * FROM proposal ORDER BY id`)
	suite.Require().NoError(err)

	expected := []dbtypes.ProposalRow{
		dbtypes.NewProposalRow(
			1,
			"proposalRoute",
			"proposalType",
			"title",
			"description",
			suite.encodeProposalContent(content1),
			time.Date(2020, 1, 1, 00, 00, 00, 000, time.UTC),
			time.Date(2020, 1, 1, 01, 00, 00, 000, time.UTC),
			time.Date(2020, 1, 1, 02, 00, 00, 000, time.UTC),
			time.Date(2020, 1, 1, 03, 00, 00, 000, time.UTC),
			proposer1.String(),
			govtypes.StatusDepositPeriod.String(),
		),
		dbtypes.NewProposalRow(
			2,
			"proposalRoute1",
			"proposalType1",
			"title1",
			"description1",
			suite.encodeProposalContent(content2),
			time.Date(2020, 1, 2, 00, 00, 00, 000, time.UTC),
			time.Date(2020, 1, 2, 01, 00, 00, 000, time.UTC),
			time.Date(2020, 1, 2, 02, 00, 00, 000, time.UTC),
			time.Date(2020, 1, 2, 03, 00, 00, 000, time.UTC),
			proposer2.String(),
			govtypes.StatusPassed.String(),
		),
	}
	for i, expect := range expected {
		suite.Require().True(proposalRow[i].Equals(expect))
	}
}

func (suite *DbTestSuite) TestBigDipperDb_GetOpenProposalsIds() {
	proposer1 := suite.getAccount("cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs")
	proposer2 := suite.getAccount("cosmos184ma3twcfjqef6k95ne8w2hk80x2kah7vcwy4a")

	content1 := govtypes.NewTextProposal("title", "description")
	content2 := govtypes.NewTextProposal("title1", "description1")
	input := []types.Proposal{
		types.NewProposal(
			1,
			"proposalRoute",
			"proposalType",
			content1,
			govtypes.StatusVotingPeriod,
			time.Date(2020, 1, 1, 00, 00, 00, 000, time.UTC),
			time.Date(2020, 1, 1, 01, 00, 00, 000, time.UTC),
			time.Date(2020, 1, 1, 02, 00, 00, 000, time.UTC),
			time.Date(2020, 1, 1, 03, 00, 00, 000, time.UTC),
			proposer1.String(),
		),
		types.NewProposal(
			2,
			"proposalRoute",
			"proposalType",
			content1,
			govtypes.StatusDepositPeriod,
			time.Date(2020, 1, 1, 00, 00, 00, 000, time.UTC),
			time.Date(2020, 1, 1, 01, 00, 00, 000, time.UTC),
			time.Date(2020, 1, 1, 02, 00, 00, 000, time.UTC),
			time.Date(2020, 1, 1, 03, 00, 00, 000, time.UTC),
			proposer1.String(),
		),
		types.NewProposal(
			3,
			"proposalRoute1",
			"proposalType1",
			content2,
			govtypes.StatusPassed,
			time.Date(2020, 1, 2, 00, 00, 00, 000, time.UTC),
			time.Date(2020, 1, 2, 01, 00, 00, 000, time.UTC),
			time.Date(2020, 1, 2, 02, 00, 00, 000, time.UTC),
			time.Date(2020, 1, 2, 03, 00, 00, 000, time.UTC),
			proposer2.String(),
		),
		types.NewProposal(
			5,
			"proposalRoute1",
			"proposalType1",
			content2,
			govtypes.StatusRejected,
			time.Date(2020, 1, 2, 00, 00, 00, 000, time.UTC),
			time.Date(2020, 1, 2, 01, 00, 00, 000, time.UTC),
			time.Date(2020, 1, 2, 02, 00, 00, 000, time.UTC),
			time.Date(2020, 1, 2, 03, 00, 00, 000, time.UTC),
			proposer2.String(),
		),
	}

	err := suite.database.SaveProposals(input)
	suite.Require().NoError(err)

	ids, err := suite.database.GetOpenProposalsIds()
	suite.Require().NoError(err)
	suite.Require().Equal([]uint64{1, 2}, ids)
}

func (suite *DbTestSuite) TestBigDipperDb_UpdateProposal() {
	proposal := suite.getProposalRow(1)
	proposer, err := sdk.AccAddressFromBech32(proposal.Proposer)
	suite.Require().NoError(err)

	update := types.NewProposalUpdate(
		proposal.ProposalID,
		govtypes.StatusPassed,
		time.Date(2020, 1, 1, 00, 00, 00, 000, time.UTC),
		time.Date(2020, 1, 1, 01, 00, 00, 000, time.UTC),
	)

	err = suite.database.UpdateProposal(update)
	suite.Require().NoError(err)

	expected := dbtypes.NewProposalRow(
		proposal.ProposalID,
		proposal.ProposalRoute,
		proposal.ProposalType,
		proposal.Content.GetTitle(),
		proposal.Content.GetDescription(),
		suite.encodeProposalContent(proposal.Content),
		proposal.SubmitTime,
		proposal.DepositEndTime,
		time.Date(2020, 1, 1, 00, 00, 00, 000, time.UTC),
		time.Date(2020, 1, 1, 01, 00, 00, 000, time.UTC),
		proposer.String(),
		govtypes.StatusPassed.String(),
	)

	var result []dbtypes.ProposalRow
	err = suite.database.Sqlx.Select(&result, `SELECT * FROM proposal`)
	suite.Require().NoError(err)
	suite.Require().Len(result, 1)
	for _, r := range result {
		suite.Require().True(expected.Equals(r))
	}
}

// -------------------------------------------------------------------------------------------------------------------

func (suite *DbTestSuite) TestBigDipperDb_SaveVote() {
	proposal := suite.getProposalRow(1)
	voter := suite.getAccount("cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs")

	vote := types.NewVote(1, voter.String(), govtypes.OptionYes, 1)
	err := suite.database.SaveVote(vote)
	suite.Require().NoError(err)

	expected := dbtypes.NewVoteRow(int64(proposal.ProposalID), voter.String(), govtypes.OptionYes.String(), 1)

	var result []dbtypes.VoteRow
	err = suite.database.Sqlx.Select(&result, `SELECT * FROM proposal_vote`)
	suite.Require().NoError(err)
	suite.Require().Len(result, 1)
	suite.Require().True(expected.Equals(result[0]))
}

func (suite *DbTestSuite) TestBigDipperDb_SaveTallyResults() {
	suite.getProposalRow(1)
	suite.getProposalRow(2)

	input := []types.TallyResult{
		types.NewTallyResult(1, 1, 1, 1, 1, 1),
		types.NewTallyResult(2, 2, 2, 2, 2, 2),
	}

	err := suite.database.SaveTallyResults(input)
	suite.Require().NoError(err)

	var result []dbtypes.TallyResultRow
	err = suite.database.Sqlx.Select(&result, `SELECT * FROM proposal_tally_result`)
	suite.Require().NoError(err)

	expected := []dbtypes.TallyResultRow{
		dbtypes.NewTallyResultRow(1, 1, 1, 1, 1, 1),
		dbtypes.NewTallyResultRow(2, 2, 2, 2, 2, 2),
	}

	for i, r := range result {
		suite.Require().True(r.Equals(expected[i]))
	}
}

// -------------------------------------------------------------------------------------------------------------------

func (suite *DbTestSuite) TestBigDipperDb_SaveDeposits() {
	proposal := suite.getProposalRow(1)

	depositor := suite.getAccount("cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs")
	amount := sdk.NewCoins(sdk.NewCoin("desmos", sdk.NewInt(10000)))

	depositor2 := suite.getAccount("cosmos184ma3twcfjqef6k95ne8w2hk80x2kah7vcwy4a")
	amount2 := sdk.NewCoins(sdk.NewCoin("desmos", sdk.NewInt(30000)))

	deposit := []types.Deposit{
		types.NewDeposit(proposal.ProposalID, depositor.String(), amount, 10),
		types.NewDeposit(proposal.ProposalID, depositor2.String(), amount2, 10),
	}

	err := suite.database.SaveDeposits(deposit)
	suite.Require().NoError(err)

	expected := []dbtypes.DepositRow{
		dbtypes.NewDepositRow(1, depositor.String(), dbtypes.NewDbCoins(amount), 10),
		dbtypes.NewDepositRow(1, depositor2.String(), dbtypes.NewDbCoins(amount2), 10),
	}
	var result []dbtypes.DepositRow
	err = suite.database.Sqlx.Select(&result, `SELECT * FROM proposal_deposit`)
	suite.Require().NoError(err)
	for i, r := range result {
		suite.Require().True(expected[i].Equals(r))
	}
}
