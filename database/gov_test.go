package database_test

import (
	"time"
	"github.com/forbole/bdjuno/x/gov/types"
	gov "github.com/cosmos/cosmos-sdk/x/gov/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	dbtypes "github.com/forbole/bdjuno/database/types"
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
	status1,err:=gov.ProposalStatusFromString("DepositPeriod")
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
	status2,err:=gov.ProposalStatusFromString("Passed")
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
	),types.NewProposal("title1",
		"description1",
		"proposalRoute1",
		"proposalType1",
		1,
		status2,
		submitTime2, 
		depositEndTime2,
		votingStartTime2,
		votingEndTime2, 
		proposer2)}
	err = suite.database.SaveProposals(input)
	suite.Require().NoError(err)

	var delegationRows []dbtypes.ValidatorDelegationRow
	err = suite.database.Sqlx.Select(&delegationRows, `SELECT * FROM validator_delegation`)
}
