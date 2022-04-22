package database_test

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	profilestypes "github.com/desmos-labs/desmos/v2/x/profiles/types"

	dbtypes "github.com/forbole/bdjuno/v3/database/types"
	"github.com/forbole/bdjuno/v3/types"
)

func (suite *DbTestSuite) TestBigDipperDb_SaveProfilesParams() {
	profilesParams := profilestypes.Params{
		Nickname: profilestypes.NewNicknameParams(sdk.NewInt(1), sdk.NewInt(100)),
		DTag:     profilestypes.NewDTagParams("abc", sdk.NewInt(1), sdk.NewInt(100)),
		Bio:      profilestypes.NewBioParams(sdk.NewInt(100)),
		Oracle:   profilestypes.NewOracleParams(1, 1, 1, 1, 1, sdk.NewCoin("band", sdk.NewInt(1))),
	}
	err := suite.database.SaveProfilesParams(types.NewProfilesParams(profilesParams, 10))
	suite.Require().NoError(err)

	var rows []dbtypes.ProfilesParamsRow
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM profiles_params`)
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1)

	var stored profilestypes.Params
	err = json.Unmarshal([]byte(rows[0].Params), &stored)
	suite.Require().NoError(err)
	suite.Require().Equal(profilesParams, stored)
}
