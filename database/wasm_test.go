package database_test

import (
	"encoding/json"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	dbtypes "github.com/forbole/bdjuno/v3/database/types"
	"github.com/forbole/bdjuno/v3/types"
)

func (suite *DbTestSuite) TestSaveWasmParams() error {
	suite.getAccount("cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs")

	wasmParams := wasmtypes.Params{
		CodeUploadAccess: wasmtypes.AccessConfig{
			Permission: wasmtypes.AccessTypeEverybody,
			Address:    "cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs",
		},
		InstantiateDefaultPermission: wasmtypes.AccessTypeEverybody,
	}

	err := suite.database.SaveWasmParams(
		types.NewWasmParams(wasmParams, 10),
	)
	suite.Require().NoError(err)

	var rows []dbtypes.WasmParamsRow
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM wasm_params`)
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1)

	var stored wasmtypes.Params
	err = json.Unmarshal([]byte(rows[0].Params), &stored)
	suite.Require().NoError(err)
	suite.Require().Equal(wasmParams, stored)

	return nil
}
