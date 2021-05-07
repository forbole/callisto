package bigdipper_test

import (
	"fmt"
	"time"

	"github.com/forbole/bdjuno/database/types"

	bgovtypes "github.com/forbole/bdjuno/modules/bigdipper/gov/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	dbtypes "github.com/forbole/bdjuno/database/bigdipper/types"
)

func (suite *DbTestSuite) TestBigDipperDb_SaveProposals() {
	proposer1 := suite.getAccount("cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs")
	proposer2 := suite.getAccount("cosmos184ma3twcfjqef6k95ne8w2hk80x2kah7vcwy4a")

	input := []bgovtypes.Proposal{
		bgovtypes.NewProposal("title",
			"description",
			"proposalRoute",
			"proposalType",
			1,
			govtypes.StatusDepositPeriod,
			time.Date(2020, 1, 1, 00, 00, 00, 000, time.UTC),
			time.Date(2020, 1, 1, 01, 00, 00, 000, time.UTC),
			time.Date(2020, 1, 1, 02, 00, 00, 000, time.UTC),
			time.Date(2020, 1, 1, 03, 00, 00, 000, time.UTC),
			proposer1.String(),
		),
		bgovtypes.NewProposal("title1",
			"description1",
			"proposalRoute1",
			"proposalType1",
			2,
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
	err = suite.database.Sqlx.Select(&proposalRow, `SELECT * FROM proposal ORDER BY proposal_id`)
	suite.Require().NoError(err)

	expected := []dbtypes.ProposalRow{dbtypes.NewProposalRow("title",
		"description",
		"proposalRoute",
		"proposalType",
		1,
		time.Date(2020, 1, 1, 00, 00, 00, 000, time.UTC),
		time.Date(2020, 1, 1, 01, 00, 00, 000, time.UTC),
		time.Date(2020, 1, 1, 02, 00, 00, 000, time.UTC),
		time.Date(2020, 1, 1, 03, 00, 00, 000, time.UTC),
		proposer1.String(),
		govtypes.StatusDepositPeriod.String(),
	), dbtypes.NewProposalRow("title1",
		"description1",
		"proposalRoute1",
		"proposalType1",
		2,
		time.Date(2020, 1, 2, 00, 00, 00, 000, time.UTC),
		time.Date(2020, 1, 2, 01, 00, 00, 000, time.UTC),
		time.Date(2020, 1, 2, 02, 00, 00, 000, time.UTC),
		time.Date(2020, 1, 2, 03, 00, 00, 000, time.UTC),
		proposer2.String(),
		govtypes.StatusPassed.String(),
	)}
	for i, expect := range expected {
		suite.Require().True(proposalRow[i].Equals(expect))
	}
}

func (suite *DbTestSuite) TestBigDipperDb_SaveProposal() {
	proposer1 := suite.getAccount("cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs")

	input := bgovtypes.NewProposal("title",
		"description",
		"proposalRoute",
		"proposalType",
		1,
		govtypes.StatusDepositPeriod,
		time.Date(2020, 1, 2, 00, 00, 00, 000, time.UTC),
		time.Date(2020, 1, 2, 01, 00, 00, 000, time.UTC),
		time.Date(2020, 1, 2, 02, 00, 00, 000, time.UTC),
		time.Date(2020, 1, 2, 03, 00, 00, 000, time.UTC),
		proposer1.String(),
	)

	err := suite.database.SaveProposal(input)
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
		time.Date(2020, 1, 2, 00, 00, 00, 000, time.UTC),
		time.Date(2020, 1, 2, 01, 00, 00, 000, time.UTC),
		time.Date(2020, 1, 2, 02, 00, 00, 000, time.UTC),
		time.Date(2020, 1, 2, 03, 00, 00, 000, time.UTC),
		proposer1.String(),
		govtypes.StatusDepositPeriod.String(),
	)

	suite.Require().True(expect.Equals(proposalRow[0]))
}

func (suite *DbTestSuite) TestBigDipperDb_SaveTallyResults() {
	suite.getProposalRow(1)
	suite.getProposalRow(2)

	input := []bgovtypes.TallyResult{
		bgovtypes.NewTallyResult(1, 1, 1, 1, 1, 1),
		bgovtypes.NewTallyResult(2, 2, 2, 2, 2, 2),
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

func (suite *DbTestSuite) TestBigDipperDb_SaveVote() {
	proposal := suite.getProposalRow(1)
	voter := suite.getAccount("cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs")

	vote := bgovtypes.NewVote(1, voter.String(), govtypes.OptionYes, 1)
	err := suite.database.SaveVote(vote)
	suite.Require().NoError(err)

	expected := dbtypes.NewVoteRow(int64(proposal.ProposalID), voter.String(), govtypes.OptionYes.String(), 1)

	var result []dbtypes.VoteRow
	err = suite.database.Sqlx.Select(&result, `SELECT * FROM proposal_vote`)
	suite.Require().NoError(err)
	suite.Require().Len(result, 1)
	suite.Require().True(expected.Equals(result[0]))
}

func (suite *DbTestSuite) TestBigDipperDb_SaveDeposit() {
	proposal := suite.getProposalRow(1)
	depositor := suite.getAccount("cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs")
	amount := sdk.NewCoins(sdk.NewCoin("desmos", sdk.NewInt(10000)))

	deposit := bgovtypes.NewDeposit(proposal.ProposalID, depositor.String(), amount, 10)
	err := suite.database.SaveDeposit(deposit)
	suite.Require().NoError(err)

	var result []dbtypes.DepositRow
	err = suite.database.Sqlx.Select(&result, `SELECT * FROM proposal_deposit`)
	suite.Require().NoError(err)
	suite.Require().Len(result, 1)
	suite.Require().True(result[0].Equals(dbtypes.NewDepositRow(
		1,
		depositor.String(),
		types.NewDbCoins(amount),
		10,
	)))
}

func (suite *DbTestSuite) TestBigDipperDb_SaveDeposits() {
	proposal := suite.getProposalRow(1)

	depositor := suite.getAccount("cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs")
	amount := sdk.NewCoins(sdk.NewCoin("desmos", sdk.NewInt(10000)))

	depositor2 := suite.getAccount("cosmos184ma3twcfjqef6k95ne8w2hk80x2kah7vcwy4a")
	amount2 := sdk.NewCoins(sdk.NewCoin("desmos", sdk.NewInt(30000)))

	deposit := []bgovtypes.Deposit{
		bgovtypes.NewDeposit(proposal.ProposalID, depositor.String(), amount, 10),
		bgovtypes.NewDeposit(proposal.ProposalID, depositor2.String(), amount2, 10),
	}

	err := suite.database.SaveDeposits(deposit)
	suite.Require().NoError(err)

	expected := []dbtypes.DepositRow{
		dbtypes.NewDepositRow(1, depositor.String(), types.NewDbCoins(amount), 10),
		dbtypes.NewDepositRow(1, depositor2.String(), types.NewDbCoins(amount2), 10),
	}
	var result []dbtypes.DepositRow
	err = suite.database.Sqlx.Select(&result, `SELECT * FROM proposal_deposit`)
	suite.Require().NoError(err)
	for i, r := range result {
		suite.Require().True(expected[i].Equals(r))
	}
}

func (suite *DbTestSuite) TestBigDipperDb_UpdateProposal() {
	proposal := suite.getProposalRow(1)
	proposer, err := sdk.AccAddressFromBech32(proposal.Proposer)
	suite.Require().NoError(err)

	update := bgovtypes.NewProposal(proposal.Title,
		proposal.Description,
		proposal.ProposalRoute,
		proposal.ProposalType,
		proposal.ProposalID,
		govtypes.StatusPassed,
		proposal.SubmitTime,
		proposal.DepositEndTime,
		time.Date(2020, 1, 1, 00, 00, 00, 000, time.UTC),
		time.Date(2020, 1, 1, 01, 00, 00, 000, time.UTC),
		proposer.String(),
	)

	err = suite.database.UpdateProposal(update)
	suite.Require().NoError(err)
	expected := dbtypes.NewProposalRow(proposal.Title,
		proposal.Description,
		proposal.ProposalRoute,
		proposal.ProposalType,
		proposal.ProposalID,
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

func (suite *DbTestSuite) getProposalRow(id int) bgovtypes.Proposal {
	proposer := suite.getAccount("cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs")

	title := fmt.Sprintf("title%d", id)
	description := fmt.Sprintf("description%d", id)
	proposalRoute := fmt.Sprintf("proposalRoute%d", id)
	proposalType := fmt.Sprintf("proposalType%d", id)

	proposal := bgovtypes.NewProposal(
		title,
		description,
		proposalRoute,
		proposalType,
		uint64(id),
		govtypes.StatusPassed,
		time.Date(2020, 1, 1, 00, 00, 00, 000, time.UTC),
		time.Date(2020, 1, 1, 01, 00, 00, 000, time.UTC),
		time.Date(2020, 1, 1, 02, 00, 00, 000, time.UTC),
		time.Date(2020, 1, 1, 03, 00, 00, 000, time.UTC),
		proposer.String(),
	)

	err := suite.database.SaveProposal(proposal)
	suite.Require().NoError(err)

	return proposal
}
