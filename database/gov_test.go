package database_test

import (
	"fmt"
	"time"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/gogo/protobuf/proto"

	"github.com/forbole/bdjuno/v3/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	dbtypes "github.com/forbole/bdjuno/v3/database/types"
)

func (suite *DbTestSuite) TestBigDipperDb_SaveGovParams() {
	votingParams := govtypes.NewVotingParams(time.Second * 10)
	tallyParams := govtypes.NewTallyParams(sdk.NewDec(10), sdk.NewDec(10), sdk.NewDec(10))
	depositParams := govtypes.NewDepositParams(sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(10))), time.Minute*5)
	original := types.NewGovParams(types.NewVotingParams(votingParams), types.NewDepositParam(depositParams), types.NewTallyParams(tallyParams), 10)

	err := suite.database.SaveGovParams(original)
	suite.Require().NoError(err)

	stored, err := suite.database.GetGovParams()
	suite.Require().NoError(err)
	suite.Require().Equal(original, stored)

	// ----------------------------------------------------------------------------------------------------------------
	// Try updating with a lower height
	depositParams = govtypes.NewDepositParams(sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(1000))), time.Minute*5)
	updated := types.NewGovParams(types.NewVotingParams(votingParams), types.NewDepositParam(depositParams), types.NewTallyParams(tallyParams), 9)

	err = suite.database.SaveGovParams(updated)
	suite.Require().NoError(err)

	stored, err = suite.database.GetGovParams()
	suite.Require().NoError(err)
	suite.Require().Equal(original, stored)

	// ----------------------------------------------------------------------------------------------------------------
	// Try updating with the same height	depositParams = govtypes.NewDepositParams(sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(1000))), time.Minute*5)
	updated = types.NewGovParams(types.NewVotingParams(votingParams), types.NewDepositParam(depositParams), types.NewTallyParams(tallyParams), 10)

	err = suite.database.SaveGovParams(updated)
	suite.Require().NoError(err)

	stored, err = suite.database.GetGovParams()
	suite.Require().NoError(err)
	suite.Require().Equal(updated, stored)

	// ----------------------------------------------------------------------------------------------------------------
	// Try updating with a higher height
	tallyParams = govtypes.NewTallyParams(sdk.NewDec(100), sdk.NewDec(100), sdk.NewDec(100))
	depositParams = govtypes.NewDepositParams(sdk.NewCoins(sdk.NewCoin("udesmos", sdk.NewInt(10000))), time.Minute*5)
	updated = types.NewGovParams(types.NewVotingParams(votingParams), types.NewDepositParam(depositParams), types.NewTallyParams(tallyParams), 11)

	err = suite.database.SaveGovParams(updated)
	suite.Require().NoError(err)

	stored, err = suite.database.GetGovParams()
	suite.Require().NoError(err)
	suite.Require().Equal(updated, stored)
}

// -------------------------------------------------------------------------------------------------------------------

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
		govtypes.StatusPassed.String(),
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
			govtypes.StatusDepositPeriod.String(),
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
			govtypes.StatusPassed.String(),
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

func (suite *DbTestSuite) TestBigDipperDb_GetProposal() {
	content := govtypes.NewTextProposal("title", "description")
	proposer := suite.getAccount("cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs")
	proposal := types.NewProposal(
		1,
		"proposalRoute",
		"proposalType",
		content,
		govtypes.StatusDepositPeriod.String(),
		time.Date(2020, 1, 1, 00, 00, 00, 000, time.UTC),
		time.Date(2020, 1, 1, 01, 00, 00, 000, time.UTC),
		time.Date(2020, 1, 1, 02, 00, 00, 000, time.UTC),
		time.Date(2020, 1, 1, 03, 00, 00, 000, time.UTC),
		proposer.String(),
	)
	input := []types.Proposal{proposal}

	err := suite.database.SaveProposals(input)
	suite.Require().NoError(err)

	stored, err := suite.database.GetProposal(1)
	suite.Require().NoError(err)
	suite.Require().True(proposal.Equal(*stored))
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
			govtypes.StatusVotingPeriod.String(),
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
			govtypes.StatusDepositPeriod.String(),
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
			govtypes.StatusPassed.String(),
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
			govtypes.StatusRejected.String(),
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
		govtypes.StatusPassed.String(),
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

func (suite *DbTestSuite) TestBigDipperDb_SaveDeposits() {
	_ = suite.getBlock(9)
	_ = suite.getBlock(10)
	_ = suite.getBlock(11)

	proposal := suite.getProposalRow(1)

	depositor := suite.getAccount("cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs")
	amount := sdk.NewCoins(sdk.NewCoin("desmos", sdk.NewInt(10000)))

	depositor2 := suite.getAccount("cosmos184ma3twcfjqef6k95ne8w2hk80x2kah7vcwy4a")
	amount2 := sdk.NewCoins(sdk.NewCoin("desmos", sdk.NewInt(30000)))

	depositor3 := suite.getAccount("cosmos1gyds87lg3m52hex9yqta2mtwzw89pfukx3jl7g")
	amount3 := sdk.NewCoins(sdk.NewCoin("desmos", sdk.NewInt(50000)))

	deposit := []types.Deposit{
		types.NewDeposit(proposal.ProposalID, depositor.String(), amount, 10),
		types.NewDeposit(proposal.ProposalID, depositor2.String(), amount2, 10),
		types.NewDeposit(proposal.ProposalID, depositor3.String(), amount3, 10),
	}

	err := suite.database.SaveDeposits(deposit)
	suite.Require().NoError(err)

	expected := []dbtypes.DepositRow{
		dbtypes.NewDepositRow(1, depositor.String(), dbtypes.NewDbCoins(amount), 10),
		dbtypes.NewDepositRow(1, depositor2.String(), dbtypes.NewDbCoins(amount2), 10),
		dbtypes.NewDepositRow(1, depositor3.String(), dbtypes.NewDbCoins(amount3), 10),
	}
	var result []dbtypes.DepositRow
	err = suite.database.Sqlx.Select(&result, `SELECT * FROM proposal_deposit`)
	suite.Require().NoError(err)
	for i, r := range result {
		suite.Require().True(expected[i].Equals(r))
	}

	// ----------------------------------------------------------------------------------------------------------------
	// Update values

	amount = sdk.NewCoins(sdk.NewCoin("desmos", sdk.NewInt(10)))
	amount2 = sdk.NewCoins(sdk.NewCoin("desmos", sdk.NewInt(20)))
	amount3 = sdk.NewCoins(sdk.NewCoin("desmos", sdk.NewInt(30)))

	deposit = []types.Deposit{
		types.NewDeposit(proposal.ProposalID, depositor.String(), amount, 9),
		types.NewDeposit(proposal.ProposalID, depositor2.String(), amount2, 10),
		types.NewDeposit(proposal.ProposalID, depositor3.String(), amount3, 11),
	}

	err = suite.database.SaveDeposits(deposit)
	suite.Require().NoError(err)

	expected = []dbtypes.DepositRow{
		dbtypes.NewDepositRow(1, depositor.String(), dbtypes.NewDbCoins(
			sdk.NewCoins(sdk.NewCoin("desmos", sdk.NewInt(10000)))), 10),
		dbtypes.NewDepositRow(1, depositor2.String(), dbtypes.NewDbCoins(amount2), 10),
		dbtypes.NewDepositRow(1, depositor3.String(), dbtypes.NewDbCoins(amount3), 11),
	}

	result = []dbtypes.DepositRow{}
	err = suite.database.Sqlx.Select(&result, `SELECT * FROM proposal_deposit`)
	suite.Require().NoError(err)
	for i, r := range result {
		suite.Require().True(expected[i].Equals(r))
	}
}

// -------------------------------------------------------------------------------------------------------------------

func (suite *DbTestSuite) TestBigDipperDb_SaveVote() {
	_ = suite.getBlock(0)
	_ = suite.getBlock(1)
	_ = suite.getBlock(2)

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

	// Update with lower height should not change option
	vote = types.NewVote(1, voter.String(), govtypes.OptionNo, 0)
	err = suite.database.SaveVote(vote)
	suite.Require().NoError(err)

	result = []dbtypes.VoteRow{}
	err = suite.database.Sqlx.Select(&result, `SELECT * FROM proposal_vote`)
	suite.Require().NoError(err)
	suite.Require().Len(result, 1)
	suite.Require().True(expected.Equals(result[0]))

	// Update with same height should change option
	vote = types.NewVote(1, voter.String(), govtypes.OptionAbstain, 1)
	err = suite.database.SaveVote(vote)
	suite.Require().NoError(err)

	expected = dbtypes.NewVoteRow(int64(proposal.ProposalID), voter.String(), govtypes.OptionAbstain.String(), 1)

	result = []dbtypes.VoteRow{}
	err = suite.database.Sqlx.Select(&result, `SELECT * FROM proposal_vote`)
	suite.Require().NoError(err)
	suite.Require().Len(result, 1)
	suite.Require().True(expected.Equals(result[0]))

	// Update with higher height should change option
	vote = types.NewVote(1, voter.String(), govtypes.OptionNoWithVeto, 2)
	err = suite.database.SaveVote(vote)
	suite.Require().NoError(err)

	expected = dbtypes.NewVoteRow(int64(proposal.ProposalID), voter.String(), govtypes.OptionNoWithVeto.String(), 2)

	result = []dbtypes.VoteRow{}
	err = suite.database.Sqlx.Select(&result, `SELECT * FROM proposal_vote`)
	suite.Require().NoError(err)
	suite.Require().Len(result, 1)
	suite.Require().True(expected.Equals(result[0]))
}

func (suite *DbTestSuite) TestBigDipperDb_SaveTallyResults() {
	suite.getProposalRow(1)
	suite.getProposalRow(2)
	suite.getProposalRow(3)

	// Store the data
	err := suite.database.SaveTallyResults([]types.TallyResult{
		types.NewTallyResult(1, "1", "1", "1", "1", 2),
		types.NewTallyResult(2, "2", "2", "2", "2", 2),
		types.NewTallyResult(3, "3", "3", "3", "3", 2),
	})
	suite.Require().NoError(err)

	// Verify the data
	var result []dbtypes.TallyResultRow
	err = suite.database.Sqlx.Select(&result, `SELECT * FROM proposal_tally_result`)
	suite.Require().NoError(err)

	expected := []dbtypes.TallyResultRow{
		dbtypes.NewTallyResultRow(1, "1", "1", "1", "1", 2),
		dbtypes.NewTallyResultRow(2, "2", "2", "2", "2", 2),
		dbtypes.NewTallyResultRow(3, "3", "3", "3", "3", 2),
	}
	for i, r := range result {
		suite.Require().True(r.Equals(expected[i]))
	}

	// ----------------------------------------------------------------------------------------------------------------
	// Update the data
	err = suite.database.SaveTallyResults([]types.TallyResult{
		types.NewTallyResult(1, "10", "10", "10", "10", 1),
		types.NewTallyResult(2, "20", "20", "20", "20", 2),
		types.NewTallyResult(3, "30", "30", "30", "30", 3),
	})
	suite.Require().NoError(err)

	// Verify the data
	result = []dbtypes.TallyResultRow{}
	err = suite.database.Sqlx.Select(&result, `SELECT * FROM proposal_tally_result`)
	suite.Require().NoError(err)

	expected = []dbtypes.TallyResultRow{
		dbtypes.NewTallyResultRow(1, "1", "1", "1", "1", 2),
		dbtypes.NewTallyResultRow(2, "20", "20", "20", "20", 2),
		dbtypes.NewTallyResultRow(3, "30", "30", "30", "30", 3),
	}
	for i, r := range result {
		suite.Require().True(r.Equals(expected[i]))
	}
}

// -------------------------------------------------------------------------------------------------------------------

func (suite *DbTestSuite) TestBigDipperDb_SaveProposalStakingPoolSnapshot() {
	_ = suite.getBlock(9)
	_ = suite.getBlock(10)
	_ = suite.getBlock(11)
	_ = suite.getProposalRow(1)

	// ----------------------------------------------------------------------------------------------------------------
	// Save snapshot

	snapshot := types.NewProposalStakingPoolSnapshot(1, types.NewPool(
		sdk.NewInt(100),
		sdk.NewInt(200),
		10,
	))
	err := suite.database.SaveProposalStakingPoolSnapshot(snapshot)
	suite.Require().NoError(err)

	var rows []dbtypes.ProposalStakingPoolSnapshotRow
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM proposal_staking_pool_snapshot`)
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1)
	suite.Require().Equal(rows[0], dbtypes.NewProposalStakingPoolSnapshotRow(
		1,
		100,
		200,
		10,
	))

	// ----------------------------------------------------------------------------------------------------------------
	// Update with lower height

	err = suite.database.SaveProposalStakingPoolSnapshot(types.NewProposalStakingPoolSnapshot(1, types.NewPool(
		sdk.NewInt(200),
		sdk.NewInt(500),
		9,
	)))
	suite.Require().NoError(err)

	rows = []dbtypes.ProposalStakingPoolSnapshotRow{}
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM proposal_staking_pool_snapshot`)
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1)
	suite.Require().Equal(rows[0], dbtypes.NewProposalStakingPoolSnapshotRow(
		1,
		100,
		200,
		10,
	))

	// ----------------------------------------------------------------------------------------------------------------
	// Update with same height

	err = suite.database.SaveProposalStakingPoolSnapshot(types.NewProposalStakingPoolSnapshot(1, types.NewPool(
		sdk.NewInt(500),
		sdk.NewInt(1000),
		10,
	)))
	suite.Require().NoError(err)

	rows = []dbtypes.ProposalStakingPoolSnapshotRow{}
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM proposal_staking_pool_snapshot`)
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1)
	suite.Require().Equal(rows[0], dbtypes.NewProposalStakingPoolSnapshotRow(
		1,
		500,
		1000,
		10,
	))

	// ----------------------------------------------------------------------------------------------------------------
	// Update with higher height

	err = suite.database.SaveProposalStakingPoolSnapshot(types.NewProposalStakingPoolSnapshot(1, types.NewPool(
		sdk.NewInt(1000),
		sdk.NewInt(2000),
		11,
	)))
	suite.Require().NoError(err)

	rows = []dbtypes.ProposalStakingPoolSnapshotRow{}
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM proposal_staking_pool_snapshot`)
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1)
	suite.Require().Equal(rows[0], dbtypes.NewProposalStakingPoolSnapshotRow(
		1,
		1000,
		2000,
		11,
	))
}

func (suite *DbTestSuite) TestBigDipperDb_SaveProposalValidatorsStatusesSnapshots() {
	_ = suite.getBlock(9)
	_ = suite.getBlock(10)
	_ = suite.getBlock(11)
	_ = suite.getProposalRow(1)

	validator1 := suite.getValidator(
		"cosmosvalcons1qqqqrezrl53hujmpdch6d805ac75n220ku09rl",
		"cosmosvaloper1rcp29q3hpd246n6qak7jluqep4v006cdsc2kkl",
		"cosmosvalconspub1zcjduepq7mft6gfls57a0a42d7uhx656cckhfvtrlmw744jv4q0mvlv0dypskehfk8",
	)
	validator2 := suite.getValidator(
		"cosmosvalcons1rtst6se0nfgjy362v33jt5d05crgdyhfvvvvay",
		"cosmosvaloper1jlr62guqwrwkdt4m3y00zh2rrsamhjf9num5xr",
		"cosmosvalconspub1zcjduepq5e8w7t7k9pwfewgrwy8vn6cghk0x49chx64vt0054yl4wwsmjgrqfackxm",
	)

	// ----------------------------------------------------------------------------------------------------------------
	// Save snapshots

	var snapshots = []types.ProposalValidatorStatusSnapshot{
		types.NewProposalValidatorStatusSnapshot(
			1,
			validator1.GetConsAddr(),
			100,
			int(stakingtypes.Bonded),
			false,
			10,
		),
		types.NewProposalValidatorStatusSnapshot(
			1,
			validator2.GetConsAddr(),
			100,
			int(stakingtypes.Unbonding),
			true,
			10,
		),
	}
	err := suite.database.SaveProposalValidatorsStatusesSnapshots(snapshots)
	suite.Require().NoError(err)

	var rows []dbtypes.ProposalValidatorVotingPowerSnapshotRow
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM proposal_validator_status_snapshot`)
	suite.Require().NoError(err)
	suite.Require().Len(rows, 2)
	suite.Require().Equal(rows, []dbtypes.ProposalValidatorVotingPowerSnapshotRow{
		dbtypes.NewProposalValidatorVotingPowerSnapshotRow(
			1,
			1,
			validator1.GetConsAddr(),
			100,
			3,
			false,
			10,
		),
		dbtypes.NewProposalValidatorVotingPowerSnapshotRow(
			2,
			1,
			validator2.GetConsAddr(),
			100,
			2,
			true,
			10,
		),
	})

	// ----------------------------------------------------------------------------------------------------------------
	// Update snapshots with lower height

	snapshots = []types.ProposalValidatorStatusSnapshot{
		types.NewProposalValidatorStatusSnapshot(
			1,
			validator1.GetConsAddr(),
			10,
			int(stakingtypes.Bonded),
			true,
			9,
		),
		types.NewProposalValidatorStatusSnapshot(
			1,
			validator2.GetConsAddr(),
			700,
			int(stakingtypes.Unbonding),
			true,
			9,
		),
	}
	err = suite.database.SaveProposalValidatorsStatusesSnapshots(snapshots)
	suite.Require().NoError(err)

	rows = []dbtypes.ProposalValidatorVotingPowerSnapshotRow{}
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM proposal_validator_status_snapshot`)
	suite.Require().NoError(err)
	suite.Require().Len(rows, 2)
	suite.Require().Equal(rows, []dbtypes.ProposalValidatorVotingPowerSnapshotRow{
		dbtypes.NewProposalValidatorVotingPowerSnapshotRow(
			1,
			1,
			validator1.GetConsAddr(),
			100,
			3,
			false,
			10,
		),
		dbtypes.NewProposalValidatorVotingPowerSnapshotRow(
			2,
			1,
			validator2.GetConsAddr(),
			100,
			2,
			true,
			10,
		),
	})

	// ----------------------------------------------------------------------------------------------------------------
	// Update snapshots with same height

	snapshots = []types.ProposalValidatorStatusSnapshot{
		types.NewProposalValidatorStatusSnapshot(
			1,
			validator1.GetConsAddr(),
			10,
			int(stakingtypes.Bonded),
			true,
			10,
		),
		types.NewProposalValidatorStatusSnapshot(
			1,
			validator2.GetConsAddr(),
			700,
			int(stakingtypes.Unbonding),
			true,
			10,
		),
	}
	err = suite.database.SaveProposalValidatorsStatusesSnapshots(snapshots)
	suite.Require().NoError(err)

	rows = []dbtypes.ProposalValidatorVotingPowerSnapshotRow{}
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM proposal_validator_status_snapshot`)
	suite.Require().NoError(err)
	suite.Require().Len(rows, 2)
	suite.Require().Equal(rows, []dbtypes.ProposalValidatorVotingPowerSnapshotRow{
		dbtypes.NewProposalValidatorVotingPowerSnapshotRow(
			1,
			1,
			validator1.GetConsAddr(),
			10,
			3,
			true,
			10,
		),
		dbtypes.NewProposalValidatorVotingPowerSnapshotRow(
			2,
			1,
			validator2.GetConsAddr(),
			700,
			2,
			true,
			10,
		),
	})

	// ----------------------------------------------------------------------------------------------------------------
	// Update snapshots with higher height

	snapshots = []types.ProposalValidatorStatusSnapshot{
		types.NewProposalValidatorStatusSnapshot(
			1,
			validator1.GetConsAddr(),
			100000,
			int(stakingtypes.Unspecified),
			false,
			11,
		),
		types.NewProposalValidatorStatusSnapshot(
			1,
			validator2.GetConsAddr(),
			700000,
			int(stakingtypes.Unbonded),
			false,
			11,
		),
	}
	err = suite.database.SaveProposalValidatorsStatusesSnapshots(snapshots)
	suite.Require().NoError(err)

	rows = []dbtypes.ProposalValidatorVotingPowerSnapshotRow{}
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM proposal_validator_status_snapshot`)
	suite.Require().NoError(err)
	suite.Require().Len(rows, 2)
	suite.Require().Equal(rows, []dbtypes.ProposalValidatorVotingPowerSnapshotRow{
		dbtypes.NewProposalValidatorVotingPowerSnapshotRow(
			1,
			1,
			validator1.GetConsAddr(),
			100000,
			0,
			false,
			11,
		),
		dbtypes.NewProposalValidatorVotingPowerSnapshotRow(
			2,
			1,
			validator2.GetConsAddr(),
			700000,
			1,
			false,
			11,
		),
	})
}
