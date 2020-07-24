package database_test

import (
	dbtypes "github.com/forbole/bdjuno/database/types"
	constypes "github.com/forbole/bdjuno/x/consensus/types"
)

func (suite *DbTestSuite) TestSaveConsensus() {
	events := []constypes.ConsensusEvent{
		constypes.NewConsensusEvent(1, 0, "First"),
		constypes.NewConsensusEvent(2, 0, "Second - Round 0 "),
		constypes.NewConsensusEvent(2, 1, "Second - Round 1"),
	}

	for _, event := range events {
		err := suite.database.SaveConsensus(event)
		suite.Require().NoError(err)

		var rows []dbtypes.ConsensusRow
		err = suite.database.Sqlx.Select(&rows, "SELECT * FROM consensus")
		suite.Require().NoError(err)

		// Make sure the consensus table only holds 1 value at a time with the correct data inside
		suite.Require().Len(rows, 1)
		suite.Require().True(rows[0].Equal(dbtypes.ConsensusRow{
			Height: event.Height,
			Round:  event.Round,
			Step:   event.Step,
		}))
	}
}
