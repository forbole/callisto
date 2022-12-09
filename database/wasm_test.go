package database_test

import (
	"encoding/json"
	"time"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
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

	err := suite.database.SaveWasmParams(types.NewWasmParams(wasmParams, 10))
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

func (suite *DbTestSuite) TestSaveWasmCodes() error {
	suite.getAccount("cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs")
	address, _ := sdk.AccAddressFromBech32("cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs")

	code := wasmtypes.Code{
		CodeID:    1,
		CodeInfo:  wasmtypes.NewCodeInfo(nil, address, wasmtypes.AllowEverybody),
		CodeBytes: nil,
		Pinned:    true,
	}

	expected := types.NewWasmCode("", code.CodeBytes, &code.CodeInfo.InstantiateConfig, code.CodeID, 10)

	err := suite.database.SaveWasmCodes([]types.WasmCode{expected})
	suite.Require().NoError(err)

	var rows []dbtypes.WasmCodeRow
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM wasm_code`)
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1)

	suite.Require().Equal(expected.Sender, rows[0].Sender)
	suite.Require().Equal(expected.CodeID, rows[0].CodeID)
	suite.Require().Equal(expected.Height, rows[0].Height)

	var storedAccessConfig *wasmtypes.AccessConfig
	err = json.Unmarshal([]byte(rows[0].InstantiatePermission), &storedAccessConfig)
	suite.Require().NoError(err)
	suite.Require().Equal(expected.InstantiatePermission, storedAccessConfig)

	return nil
}

func (suite *DbTestSuite) TestSaveWasmContracts() error {
	suite.getAccount("cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs")
	address, _ := sdk.AccAddressFromBech32("cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs")

	contract := wasmtypes.Contract{
		ContractAddress: "cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs",
		ContractInfo: wasmtypes.NewContractInfo(
			1,
			address,
			address,
			"label",
			&wasmtypes.AbsoluteTxPosition{
				BlockHeight: 10,
				TxIndex:     1,
			},
		),
		ContractState: []wasmtypes.Model{},
	}

	instantiatedAt := time.Now()
	expected := types.NewWasmContract(
		contract.ContractInfo.Creator,
		contract.ContractInfo.Admin,
		contract.ContractInfo.CodeID,
		contract.ContractInfo.Label,
		[]byte(`{"key":"val"}`),
		sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(10))),
		contract.ContractAddress,
		"",
		instantiatedAt,
		contract.ContractInfo.Creator,
		"",
		contract.ContractState,
		10,
	)

	err := suite.database.SaveWasmContracts([]types.WasmContract{expected})
	suite.Require().NoError(err)

	return nil
}
