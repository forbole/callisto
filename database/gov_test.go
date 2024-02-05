package database_test

import (
	"fmt"
	"time"

	gov "github.com/cosmos/cosmos-sdk/x/gov/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"

	"github.com/forbole/callisto/v4/testutils"
	"github.com/forbole/callisto/v4/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"

	dbtypes "github.com/forbole/callisto/v4/database/types"
)

func (suite *DbTestSuite) TestBigDipperDb_SaveGovParams() {
	params := govtypesv1.Params{
		MinDeposit:                 []sdk.Coin{sdk.NewCoin("uatom", sdk.NewInt(1000))},
		MaxDepositPeriod:           testutils.NewDurationPointer(time.Duration(int64(300000000000))),
		VotingPeriod:               testutils.NewDurationPointer(time.Duration(int64(300000))),
		Quorum:                     "0.5",
		Threshold:                  "0.3",
		VetoThreshold:              "0.15",
		MinInitialDepositRatio:     "0",
		BurnVoteQuorum:             false,
		BurnProposalDepositPrevote: false,
		BurnVoteVeto:               false,
	}

	original := types.NewGovParams(&params, 10)

	err := suite.database.SaveGovParams(original)
	suite.Require().NoError(err)

	stored, err := suite.database.GetGovParams()
	suite.Require().NoError(err)
	suite.Require().Equal(original, stored)

	// ----------------------------------------------------------------------------------------------------------------
	// Try updating with a lower height
	params.BurnVoteQuorum = false
	updated := types.NewGovParams(&params, 9)

	err = suite.database.SaveGovParams(updated)
	suite.Require().NoError(err)

	stored, err = suite.database.GetGovParams()
	suite.Require().NoError(err)
	suite.Require().Equal(original, stored)

	// ----------------------------------------------------------------------------------------------------------------
	// Try updating with the same height
	params.BurnProposalDepositPrevote = true
	updated = types.NewGovParams(&params, 10)

	err = suite.database.SaveGovParams(updated)
	suite.Require().NoError(err)

	stored, err = suite.database.GetGovParams()
	suite.Require().NoError(err)
	suite.Require().Equal(updated, stored)

	// ----------------------------------------------------------------------------------------------------------------
	// Try updating with a higher height
	params.BurnVoteVeto = true
	updated = types.NewGovParams(&params, 11)

	err = suite.database.SaveGovParams(updated)
	suite.Require().NoError(err)

	stored, err = suite.database.GetGovParams()
	suite.Require().NoError(err)
	suite.Require().Equal(updated, stored)
}

// -------------------------------------------------------------------------------------------------------------------

func (suite *DbTestSuite) getProposalRow(id int) types.Proposal {
	proposer := suite.getAccount("cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs")

	msgAny, err := codectypes.NewAnyWithValue(&govtypesv1.MsgUpdateParams{
		Authority: authtypes.NewModuleAddress(gov.ModuleName).String(),
		Params: govtypesv1.Params{
			MinDeposit:                 []sdk.Coin{sdk.NewCoin("uatom", sdk.NewInt(1000))},
			MaxDepositPeriod:           testutils.NewDurationPointer(time.Duration(int64(300000000000))),
			VotingPeriod:               testutils.NewDurationPointer(time.Duration(int64(300000))),
			Quorum:                     "0.5",
			Threshold:                  "0.3",
			VetoThreshold:              "0.15",
			MinInitialDepositRatio:     "0",
			BurnVoteQuorum:             false,
			BurnProposalDepositPrevote: false,
			BurnVoteVeto:               false,
		},
	})
	suite.Require().NoError(err)

	proposal := types.NewProposal(
		uint64(id),
		fmt.Sprintf("Proposal %d", id),
		fmt.Sprintf("Description of proposal %d", id),
		fmt.Sprintf("Metadata of proposal %d", id),
		[]*codectypes.Any{msgAny},
		govtypesv1.StatusVotingPeriod.String(),
		time.Date(2020, 1, 1, 00, 00, 00, 000, time.UTC),
		time.Date(2020, 1, 1, 01, 00, 00, 000, time.UTC),
		testutils.NewTimePointer(time.Date(2020, 1, 1, 02, 00, 00, 000, time.UTC)),
		testutils.NewTimePointer(time.Date(2020, 1, 1, 03, 00, 00, 000, time.UTC)),
		proposer.String(),
	)

	err = suite.database.SaveProposals([]types.Proposal{proposal})
	suite.Require().NoError(err)

	return proposal
}

func (suite *DbTestSuite) TestBigDipperDb_SaveProposals() {
	proposer1 := suite.getAccount("cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs")
	msgAny, err := codectypes.NewAnyWithValue(&govtypesv1.MsgUpdateParams{
		Authority: authtypes.NewModuleAddress(gov.ModuleName).String(),
		Params: govtypesv1.Params{
			MinDeposit:                 []sdk.Coin{sdk.NewCoin("uatom", sdk.NewInt(1000))},
			MaxDepositPeriod:           testutils.NewDurationPointer(time.Duration(int64(300000000000))),
			VotingPeriod:               testutils.NewDurationPointer(time.Duration(int64(300000))),
			Quorum:                     "0.5",
			Threshold:                  "0.3",
			VetoThreshold:              "0.15",
			MinInitialDepositRatio:     "0",
			BurnVoteQuorum:             false,
			BurnProposalDepositPrevote: false,
			BurnVoteVeto:               false,
		},
	})
	suite.Require().NoError(err)

	proposer2 := suite.getAccount("cosmos184ma3twcfjqef6k95ne8w2hk80x2kah7vcwy4a")

	input := []types.Proposal{
		types.NewProposal(
			1,
			"Proposal Title 1",
			"Proposal Description 1",
			"Proposal Metadata 1",
			[]*codectypes.Any{msgAny},
			govtypesv1.StatusDepositPeriod.String(),
			time.Date(2020, 1, 1, 00, 00, 00, 000, time.UTC),
			time.Date(2020, 1, 1, 01, 00, 00, 000, time.UTC),
			testutils.NewTimePointer(time.Date(2020, 1, 1, 02, 00, 00, 000, time.UTC)),
			testutils.NewTimePointer(time.Date(2020, 1, 1, 03, 00, 00, 000, time.UTC)),
			proposer1.String(),
		),
		types.NewProposal(
			2,
			"Proposal Title 2",
			"Proposal Description 2",
			"Proposal Metadata 2",
			nil,
			govtypesv1.StatusPassed.String(),
			time.Date(2020, 1, 2, 00, 00, 00, 000, time.UTC),
			time.Date(2020, 1, 2, 01, 00, 00, 000, time.UTC),
			testutils.NewTimePointer(time.Date(2020, 1, 2, 02, 00, 00, 000, time.UTC)),
			testutils.NewTimePointer(time.Date(2020, 1, 2, 03, 00, 00, 000, time.UTC)),
			proposer2.String(),
		),
	}

	err = suite.database.SaveProposals(input)
	suite.Require().NoError(err)

	var proposalRow []dbtypes.ProposalRow
	err = suite.database.Sqlx.Select(&proposalRow, `SELECT * FROM proposal ORDER BY id`)
	suite.Require().NoError(err)

	expected := []dbtypes.ProposalRow{
		dbtypes.NewProposalRow(
			1,
			"Proposal Title 1",
			"Proposal Description 1",
			"Proposal Metadata 1",
			"[{\"@type\": \"/cosmos.gov.v1.MsgUpdateParams\", \"params\": {\"quorum\": \"0.5\", \"threshold\": \"0.3\", \"min_deposit\": [{\"denom\": \"uatom\", \"amount\": \"1000\"}], \"voting_period\": \"0.000300s\", \"burn_vote_veto\": false, \"veto_threshold\": \"0.15\", \"burn_vote_quorum\": false, \"max_deposit_period\": \"300s\", \"min_initial_deposit_ratio\": \"0\", \"burn_proposal_deposit_prevote\": false}, \"authority\": \"cosmos10d07y265gmmuvt4z0w9aw880jnsr700j6zn9kn\"}]",
			time.Date(2020, 1, 1, 00, 00, 00, 000, time.UTC),
			time.Date(2020, 1, 1, 01, 00, 00, 000, time.UTC),
			testutils.NewTimePointer(time.Date(2020, 1, 1, 02, 00, 00, 000, time.UTC)),
			testutils.NewTimePointer(time.Date(2020, 1, 1, 03, 00, 00, 000, time.UTC)),
			proposer1.String(),
			govtypesv1.StatusDepositPeriod.String(),
		),
		dbtypes.NewProposalRow(
			2,
			"Proposal Title 2",
			"Proposal Description 2",
			"Proposal Metadata 2",
			"[]",
			time.Date(2020, 1, 2, 00, 00, 00, 000, time.UTC),
			time.Date(2020, 1, 2, 01, 00, 00, 000, time.UTC),
			testutils.NewTimePointer(time.Date(2020, 1, 2, 02, 00, 00, 000, time.UTC)),
			testutils.NewTimePointer(time.Date(2020, 1, 2, 03, 00, 00, 000, time.UTC)),
			proposer2.String(),
			govtypesv1.StatusPassed.String(),
		),
	}
	for i, expect := range expected {
		suite.Require().True(proposalRow[i].Equals(expect))
	}
}

func (suite *DbTestSuite) TestBigDipperDb_GetProposal() {
	proposer := suite.getAccount("cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs")
	msgAny, err := codectypes.NewAnyWithValue(&govtypesv1.MsgUpdateParams{
		Authority: authtypes.NewModuleAddress(gov.ModuleName).String(),
		Params: govtypesv1.Params{
			MinDeposit:                 []sdk.Coin{sdk.NewCoin("uatom", sdk.NewInt(1000))},
			MaxDepositPeriod:           testutils.NewDurationPointer(time.Duration(int64(300000000000))),
			VotingPeriod:               testutils.NewDurationPointer(time.Duration(int64(300000))),
			Quorum:                     "0.5",
			Threshold:                  "0.3",
			VetoThreshold:              "0.15",
			MinInitialDepositRatio:     "0",
			BurnVoteQuorum:             false,
			BurnProposalDepositPrevote: false,
			BurnVoteVeto:               false,
		},
	})
	suite.Require().NoError(err)
	proposal := types.NewProposal(
		1,
		"Proposal Title 1",
		"Proposal Description 1",
		"Proposal Metadata 1",
		[]*codectypes.Any{msgAny},
		govtypesv1.StatusDepositPeriod.String(),
		time.Date(2020, 1, 1, 00, 00, 00, 000, time.UTC),
		time.Date(2020, 1, 1, 01, 00, 00, 000, time.UTC),
		testutils.NewTimePointer(time.Date(2020, 1, 1, 02, 00, 00, 000, time.UTC)),
		testutils.NewTimePointer(time.Date(2020, 1, 1, 03, 00, 00, 000, time.UTC)),
		proposer.String(),
	)
	input := []types.Proposal{proposal}

	err = suite.database.SaveProposals(input)
	suite.Require().NoError(err)

	var rows []dbtypes.ProposalRow
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM proposal`)
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1)
}

func (suite *DbTestSuite) TestBigDipperDb_GetOpenProposalsIds() {
	proposer1 := suite.getAccount("cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs")
	proposer2 := suite.getAccount("cosmos184ma3twcfjqef6k95ne8w2hk80x2kah7vcwy4a")

	invalidProposal := types.NewProposal(
		6,
		"Proposal Title 6",
		"Proposal Description 6",
		"Proposal Metadata 6",
		nil,
		types.ProposalStatusInvalid,
		time.Date(2020, 1, 2, 00, 00, 00, 000, time.UTC),
		time.Date(2020, 1, 2, 01, 00, 00, 000, time.UTC),
		testutils.NewTimePointer(time.Date(2020, 1, 2, 02, 00, 00, 000, time.UTC)),
		testutils.NewTimePointer(time.Date(2020, 1, 2, 03, 00, 00, 000, time.UTC)),
		proposer2.String(),
	)

	input := []types.Proposal{
		types.NewProposal(
			1,
			"Proposal Title 2",
			"Proposal Description 2",
			"Proposal Metadata 2",
			nil,
			govtypesv1.StatusVotingPeriod.String(),
			time.Date(2020, 1, 1, 00, 00, 00, 000, time.UTC),
			time.Date(2020, 1, 1, 01, 00, 00, 000, time.UTC),
			testutils.NewTimePointer(time.Date(2020, 1, 1, 02, 00, 00, 000, time.UTC)),
			testutils.NewTimePointer(time.Date(2020, 1, 1, 03, 00, 00, 000, time.UTC)),
			proposer1.String(),
		),
		types.NewProposal(
			2,
			"Proposal Title 2",
			"Proposal Description 2",
			"Proposal Metadata 2",
			nil,
			govtypesv1.StatusDepositPeriod.String(),
			time.Date(2020, 1, 1, 00, 00, 00, 000, time.UTC),
			time.Date(2020, 1, 1, 01, 00, 00, 000, time.UTC),
			testutils.NewTimePointer(time.Date(2020, 1, 1, 02, 00, 00, 000, time.UTC)),
			testutils.NewTimePointer(time.Date(2020, 1, 1, 03, 00, 00, 000, time.UTC)),
			proposer1.String(),
		),
		types.NewProposal(
			3,
			"Proposal Title 3",
			"Proposal Description 3",
			"Proposal Metadata 3",
			nil,
			govtypesv1.StatusPassed.String(),
			time.Date(2020, 1, 2, 00, 00, 00, 000, time.UTC),
			time.Date(2020, 1, 2, 01, 00, 00, 000, time.UTC),
			testutils.NewTimePointer(time.Date(2020, 1, 2, 02, 00, 00, 000, time.UTC)),
			testutils.NewTimePointer(time.Date(2020, 1, 2, 03, 00, 00, 000, time.UTC)),
			proposer2.String(),
		),
		types.NewProposal(
			5,
			"Proposal Title 5",
			"Proposal Description 5",
			"Proposal Metadata 5",
			nil,
			govtypesv1.StatusRejected.String(),
			time.Date(2020, 1, 2, 00, 00, 00, 000, time.UTC),
			time.Date(2020, 1, 2, 01, 00, 00, 000, time.UTC),
			testutils.NewTimePointer(time.Date(2020, 1, 2, 02, 00, 00, 000, time.UTC)),
			testutils.NewTimePointer(time.Date(2020, 1, 2, 03, 00, 00, 000, time.UTC)),
			proposer2.String(),
		),
		invalidProposal,
	}

	err := suite.database.SaveProposals(input)
	suite.Require().NoError(err)

	timeBeforeDepositEnd := invalidProposal.DepositEndTime.Add(-1 * time.Hour)
	ids, err := suite.database.GetOpenProposalsIds(timeBeforeDepositEnd)
	suite.Require().NoError(err)
	suite.Require().Equal([]uint64{1, 2, 6}, ids)
}

func (suite *DbTestSuite) TestBigDipperDb_UpdateProposal() {
	proposal := suite.getProposalRow(1)
	proposer, err := sdk.AccAddressFromBech32(proposal.Proposer)
	suite.Require().NoError(err)

	timestamp1 := testutils.NewTimePointer(time.Date(2020, 1, 1, 00, 00, 00, 000, time.UTC))
	timestamp2 := testutils.NewTimePointer(time.Date(2020, 1, 1, 01, 00, 00, 000, time.UTC))

	update := types.NewProposalUpdate(
		proposal.ID,
		govtypesv1.StatusPassed.String(),
		timestamp1,
		timestamp2,
	)

	err = suite.database.UpdateProposal(update)
	suite.Require().NoError(err)

	expected := dbtypes.NewProposalRow(
		proposal.ID,
		"Proposal 1",
		"Description of proposal 1",
		"Metadata of proposal 1",
		"[{\"@type\": \"/cosmos.gov.v1.MsgUpdateParams\", \"params\": {\"quorum\": \"0.5\", \"threshold\": \"0.3\", \"min_deposit\": [{\"denom\": \"uatom\", \"amount\": \"1000\"}], \"voting_period\": \"0.000300s\", \"burn_vote_veto\": false, \"veto_threshold\": \"0.15\", \"burn_vote_quorum\": false, \"max_deposit_period\": \"300s\", \"min_initial_deposit_ratio\": \"0\", \"burn_proposal_deposit_prevote\": false}, \"authority\": \"cosmos10d07y265gmmuvt4z0w9aw880jnsr700j6zn9kn\"}]",
		proposal.SubmitTime,
		proposal.DepositEndTime,
		timestamp1,
		timestamp2,
		proposer.String(),
		govtypesv1.StatusPassed.String(),
	)

	var stored dbtypes.ProposalRow
	err = suite.database.SQL.Get(&stored, `SELECT * FROM proposal LIMIT 1`)
	suite.Require().NoError(err)
	suite.Require().True(expected.Equals(stored))
}

// -------------------------------------------------------------------------------------------------------------------

func (suite *DbTestSuite) TestBigDipperDb_SaveDeposits() {
	_ = suite.getBlock(9)
	_ = suite.getBlock(10)
	_ = suite.getBlock(11)

	proposal := suite.getProposalRow(1)

	depositor := suite.getAccount("cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs")
	amount := sdk.NewCoins(sdk.NewCoin("desmos", sdk.NewInt(10000)))
	txHash := "D40FE0C386FA85677FFB9B3C4CECD54CF2CD7ABECE4EF15FAEF328FCCBF4C3A8"

	depositor2 := suite.getAccount("cosmos184ma3twcfjqef6k95ne8w2hk80x2kah7vcwy4a")
	amount2 := sdk.NewCoins(sdk.NewCoin("desmos", sdk.NewInt(30000)))
	txHash2 := "40A9812A137256E88593E19428E006C01D87DB35F60F8D14739B4A46AC3C67A5"

	depositor3 := suite.getAccount("cosmos1gyds87lg3m52hex9yqta2mtwzw89pfukx3jl7g")
	amount3 := sdk.NewCoins(sdk.NewCoin("desmos", sdk.NewInt(50000)))
	txHash3 := "086CFE10741EF3800DB7F72B1666DE298DD40913BBB84C5530C87AF5EDE8027A"

	timestamp1 := time.Date(2020, 1, 1, 15, 00, 00, 000, time.UTC)
	timestamp2 := time.Date(2020, 1, 1, 16, 00, 00, 000, time.UTC)
	timestamp3 := time.Date(2020, 1, 1, 17, 00, 00, 000, time.UTC)

	deposit := []types.Deposit{
		types.NewDeposit(proposal.ID, depositor.String(), amount, timestamp1, txHash, 10),
		types.NewDeposit(proposal.ID, depositor2.String(), amount2, timestamp2, txHash2, 10),
		types.NewDeposit(proposal.ID, depositor3.String(), amount3, timestamp3, txHash3, 10),
	}

	err := suite.database.SaveDeposits(deposit)
	suite.Require().NoError(err)

	expected := []dbtypes.DepositRow{
		dbtypes.NewDepositRow(1, depositor.String(), dbtypes.NewDbCoins(amount), timestamp1, txHash, 10),
		dbtypes.NewDepositRow(1, depositor2.String(), dbtypes.NewDbCoins(amount2), timestamp2, txHash2, 10),
		dbtypes.NewDepositRow(1, depositor3.String(), dbtypes.NewDbCoins(amount3), timestamp3, txHash3, 10),
	}
	var result []dbtypes.DepositRow
	err = suite.database.Sqlx.Select(&result, `SELECT * FROM proposal_deposit`)
	suite.Require().NoError(err)
	for i, r := range result {
		suite.Require().True(expected[i].Equals(r))
	}

	// ----------------------------------------------------------------------------------------------------------------
	// Update values

	deposit = []types.Deposit{
		types.NewDeposit(proposal.ID, depositor.String(), amount, timestamp1, "8E6EA32C656A6EED84132425533E897D458F1B877080DF842B068C4AS92WP01A", 9),
		types.NewDeposit(proposal.ID, depositor2.String(), amount2, timestamp2, txHash2, 11),
		types.NewDeposit(proposal.ID, depositor3.String(), sdk.NewCoins(sdk.NewCoin("desmos", sdk.NewInt(30))), timestamp3, txHash3, 11),
	}

	err = suite.database.SaveDeposits(deposit)
	suite.Require().NoError(err)

	expected = []dbtypes.DepositRow{
		dbtypes.NewDepositRow(1, depositor.String(), dbtypes.NewDbCoins(amount), timestamp1, txHash, 10),
		dbtypes.NewDepositRow(1, depositor.String(), dbtypes.NewDbCoins(amount), timestamp1, "8E6EA32C656A6EED84132425533E897D458F1B877080DF842B068C4AS92WP01A", 9),
		dbtypes.NewDepositRow(1, depositor2.String(), dbtypes.NewDbCoins(amount2), timestamp2, txHash2, 11),
		dbtypes.NewDepositRow(1, depositor3.String(), dbtypes.NewDbCoins(sdk.NewCoins(sdk.NewCoin("desmos", sdk.NewInt(30)))), timestamp3, txHash3, 11),
	}

	result = []dbtypes.DepositRow{}
	err = suite.database.SQL.Select(&result, `SELECT * FROM proposal_deposit`)
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

	timestamp := time.Date(2020, 1, 1, 15, 00, 00, 000, time.UTC)

	vote := types.NewVote(1, voter.String(), govtypesv1.OptionYes, "0.5", timestamp, 1)
	err := suite.database.SaveVote(vote)
	suite.Require().NoError(err)

	vote2 := types.NewVote(1, voter.String(), govtypesv1.OptionNo, "0.5", timestamp, 1)
	err = suite.database.SaveVote(vote2)
	suite.Require().NoError(err)

	expected := []dbtypes.VoteRow{
		dbtypes.NewVoteRow(int64(proposal.ID), voter.String(), govtypesv1.OptionYes.String(), "0.5", timestamp, 1),
		dbtypes.NewVoteRow(int64(proposal.ID), voter.String(), govtypesv1.OptionNo.String(), "0.5", timestamp, 1),
	}

	var result []dbtypes.VoteRow
	err = suite.database.SQL.Select(&result, `SELECT * FROM proposal_vote`)
	suite.Require().NoError(err)
	suite.Require().Len(result, 2)
	for i, r := range result {
		suite.Require().True(expected[i].Equals(r))
	}

	// Update with lower height should not change option
	vote = types.NewVote(1, voter.String(), govtypesv1.OptionYes, "0.7", timestamp, 0)
	err = suite.database.SaveVote(vote)
	suite.Require().NoError(err)

	vote2 = types.NewVote(1, voter.String(), govtypesv1.OptionNo, "0.3", timestamp, 0)
	err = suite.database.SaveVote(vote2)
	suite.Require().NoError(err)

	result = []dbtypes.VoteRow{}
	err = suite.database.SQL.Select(&result, `SELECT * FROM proposal_vote`)
	suite.Require().NoError(err)
	suite.Require().Len(result, 2)
	for i, r := range result {
		suite.Require().True(expected[i].Equals(r))
	}

	// Update with same height should change option
	vote = types.NewVote(1, voter.String(), govtypesv1.OptionYes, "0.6", timestamp, 1)
	err = suite.database.SaveVote(vote)
	suite.Require().NoError(err)

	vote2 = types.NewVote(1, voter.String(), govtypesv1.OptionNo, "0.4", timestamp, 1)
	err = suite.database.SaveVote(vote2)
	suite.Require().NoError(err)

	expected = []dbtypes.VoteRow{
		dbtypes.NewVoteRow(int64(proposal.ID), voter.String(), govtypesv1.OptionYes.String(), "0.6", timestamp, 1),
		dbtypes.NewVoteRow(int64(proposal.ID), voter.String(), govtypesv1.OptionNo.String(), "0.4", timestamp, 1),
	}

	result = []dbtypes.VoteRow{}
	err = suite.database.Sqlx.Select(&result, `SELECT * FROM proposal_vote`)
	suite.Require().NoError(err)
	suite.Require().Len(result, 2)
	for i, r := range result {
		suite.Require().True(expected[i].Equals(r))
	}

	// Update with higher height should change option
	vote = types.NewVote(1, voter.String(), govtypesv1.OptionYes, "0.6", timestamp, 3)
	err = suite.database.SaveVote(vote)
	suite.Require().NoError(err)

	vote2 = types.NewVote(1, voter.String(), govtypesv1.OptionNo, "0.4", timestamp, 3)
	err = suite.database.SaveVote(vote2)
	suite.Require().NoError(err)

	expected = []dbtypes.VoteRow{
		dbtypes.NewVoteRow(int64(proposal.ID), voter.String(), govtypesv1.OptionYes.String(), "0.6", timestamp, 3),
		dbtypes.NewVoteRow(int64(proposal.ID), voter.String(), govtypesv1.OptionNo.String(), "0.4", timestamp, 3),
	}
	result = []dbtypes.VoteRow{}
	err = suite.database.Sqlx.Select(&result, `SELECT * FROM proposal_vote`)
	suite.Require().NoError(err)
	suite.Require().Len(result, 2)
	for i, r := range result {
		suite.Require().True(expected[i].Equals(r))
	}
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

	snapshot := types.NewProposalStakingPoolSnapshot(1, types.NewPoolSnapshot(
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

	err = suite.database.SaveProposalStakingPoolSnapshot(types.NewProposalStakingPoolSnapshot(1, types.NewPoolSnapshot(
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

	err = suite.database.SaveProposalStakingPoolSnapshot(types.NewProposalStakingPoolSnapshot(1, types.NewPoolSnapshot(
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

	err = suite.database.SaveProposalStakingPoolSnapshot(types.NewProposalStakingPoolSnapshot(1, types.NewPoolSnapshot(
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
			stakingtypes.Bonded,
			false,
			10,
		),
		types.NewProposalValidatorStatusSnapshot(
			1,
			validator2.GetConsAddr(),
			100,
			stakingtypes.Unbonding,
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
			stakingtypes.Bonded,
			true,
			9,
		),
		types.NewProposalValidatorStatusSnapshot(
			1,
			validator2.GetConsAddr(),
			700,
			stakingtypes.Unbonding,
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
			stakingtypes.Bonded,
			true,
			10,
		),
		types.NewProposalValidatorStatusSnapshot(
			1,
			validator2.GetConsAddr(),
			700,
			stakingtypes.Unbonding,
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
			stakingtypes.Unspecified,
			false,
			11,
		),
		types.NewProposalValidatorStatusSnapshot(
			1,
			validator2.GetConsAddr(),
			700000,
			stakingtypes.Unbonded,
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

func (suite *DbTestSuite) TestBigDipperDb_SaveSoftwareUpgradePlan() {
	_ = suite.getProposalRow(1)

	// ----------------------------------------------------------------------------------------------------------------
	// Save software upgrade plan at height 10 with upgrade height at 100
	var plan = upgradetypes.Plan{
		Name:   "name",
		Height: 100,
		Info:   "info",
	}

	err := suite.database.SaveSoftwareUpgradePlan(1, plan, 10)
	suite.Require().NoError(err)

	var rows []dbtypes.SoftwareUpgradePlanRow
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM software_upgrade_plan`)
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1)
	suite.Require().Equal(rows, []dbtypes.SoftwareUpgradePlanRow{
		dbtypes.NewSoftwareUpgradePlanRow(1, plan.Name, plan.Height, plan.Info, 10),
	})

	// ----------------------------------------------------------------------------------------------------------------
	// Update software upgrade plan with lower height
	planEdit1 := upgradetypes.Plan{
		Name:   "name_edit_1",
		Height: 101,
		Info:   "info_edit_1",
	}

	err = suite.database.SaveSoftwareUpgradePlan(1, planEdit1, 9)
	suite.Require().NoError(err)

	rows = []dbtypes.SoftwareUpgradePlanRow{}
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM software_upgrade_plan`)
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1)
	suite.Require().Equal(rows, []dbtypes.SoftwareUpgradePlanRow{
		dbtypes.NewSoftwareUpgradePlanRow(1, plan.Name, plan.Height, plan.Info, 10),
	})

	// ----------------------------------------------------------------------------------------------------------------
	// Update software upgrade plan with same height
	planEdit2 := upgradetypes.Plan{
		Name:   "name_edit_2",
		Height: 102,
		Info:   "info_edit_2",
	}

	err = suite.database.SaveSoftwareUpgradePlan(1, planEdit2, 10)
	suite.Require().NoError(err)

	rows = []dbtypes.SoftwareUpgradePlanRow{}
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM software_upgrade_plan`)
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1)
	suite.Require().Equal(rows, []dbtypes.SoftwareUpgradePlanRow{
		dbtypes.NewSoftwareUpgradePlanRow(1, planEdit2.Name, planEdit2.Height, planEdit2.Info, 10),
	})

	// ----------------------------------------------------------------------------------------------------------------
	// Update software upgrade plan with higher height
	planEdit3 := upgradetypes.Plan{
		Name:   "name_edit_3",
		Height: 103,
		Info:   "info_edit_3",
	}

	err = suite.database.SaveSoftwareUpgradePlan(1, planEdit3, 11)
	suite.Require().NoError(err)

	rows = []dbtypes.SoftwareUpgradePlanRow{}
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM software_upgrade_plan`)
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1)
	suite.Require().Equal(rows, []dbtypes.SoftwareUpgradePlanRow{
		dbtypes.NewSoftwareUpgradePlanRow(1, planEdit3.Name, planEdit3.Height, planEdit3.Info, 11),
	})
}

func (suite *DbTestSuite) TestBigDipperDb_DeleteSoftwareUpgradePlan() {
	_ = suite.getProposalRow(1)

	// Save software upgrade plan at height 10 with upgrade height at 100
	var plan = upgradetypes.Plan{
		Name:   "name",
		Height: 100,
		Info:   "info",
	}

	err := suite.database.SaveSoftwareUpgradePlan(1, plan, 10)
	suite.Require().NoError(err)

	// Delete software upgrade plan
	err = suite.database.DeleteSoftwareUpgradePlan(1)
	suite.Require().NoError(err)

	var rows []dbtypes.SoftwareUpgradePlanRow
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM software_upgrade_plan`)
	suite.Require().NoError(err)
	suite.Require().Len(rows, 0)

}

func (suite *DbTestSuite) TestBigDipperDb_CheckSoftwareUpgradePlan() {
	_ = suite.getProposalRow(1)

	// Save software upgrade plan at height 10 with upgrade height at 100
	var plan = upgradetypes.Plan{
		Name: "name",
		// the Height here is the upgrade height
		Height: 100,
		Info:   "info",
	}

	err := suite.database.SaveSoftwareUpgradePlan(1, plan, 10)
	suite.Require().NoError(err)

	// Check software upgrade plan at existing height
	exist, err := suite.database.CheckSoftwareUpgradePlan(100)
	suite.Require().NoError(err)
	suite.Require().Equal(true, exist)

	// Check software upgrade plan at non-existing height
	exist, err = suite.database.CheckSoftwareUpgradePlan(11)
	suite.Require().NoError(err)
	suite.Require().Equal(false, exist)
}
