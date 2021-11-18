package database_test

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	profilestypes "github.com/desmos-labs/desmos/x/profiles/types"
	dbtypes "github.com/forbole/bdjuno/v2/database/types"
	"github.com/forbole/bdjuno/v2/types"
)

func (suite *DbTestSuite) TestBigDipperDb_SaveProfilesParams() {
	nickNameParams := profilestypes.NewNicknameParams(sdk.NewInt(1), sdk.NewInt(100))
	dTagParams := profilestypes.NewDTagParams("abc", sdk.NewInt(1), sdk.NewInt(100))
	bioParams := profilestypes.NewBioParams(sdk.NewInt(100))
	oracleParams := profilestypes.NewOracleParams(1, 1, 1, 1, 1, "payer", sdk.NewCoin("band", sdk.NewInt(1)))
	profilesParams := profilestypes.Params{
		Nickname: nickNameParams,
		DTag:     dTagParams,
		Bio:      bioParams,
		Oracle:   oracleParams,
	}
	err := suite.database.SaveProfilesParams(types.NewProfilesParams(profilesParams, 10))
	suite.Require().NoError(err)

	var rows []dbtypes.ProfilesParamsRow
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM profiles_params`)
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1, "profiles_params table should contain only one row")

	var storedNicknameParams profilestypes.NicknameParams
	err = json.Unmarshal([]byte(rows[0].NickNameParams), &storedNicknameParams)
	suite.Require().NoError(err)
	suite.Require().Equal(profilesParams.Nickname, storedNicknameParams)

	var storedDTagParams profilestypes.DTagParams
	err = json.Unmarshal([]byte(rows[0].DTagParams), &storedDTagParams)
	suite.Require().NoError(err)
	suite.Require().Equal(profilesParams.DTag, storedDTagParams)

	var storedBioParams profilestypes.BioParams
	err = json.Unmarshal([]byte(rows[0].BioParams), &storedBioParams)
	suite.Require().NoError(err)
	suite.Require().Equal(profilesParams.Bio, storedBioParams)

	var storedOracleParams profilestypes.OracleParams
	err = json.Unmarshal([]byte(rows[0].OracleParams), &storedOracleParams)
	suite.Require().NoError(err)
	suite.Require().Equal(profilesParams.Oracle, storedOracleParams)

	suite.Require().Equal(int64(10), rows[0].Height)
}
