package database_test

import (
	"encoding/json"
	"time"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	dbtypes "github.com/forbole/callisto/v4/database/types"
	"github.com/forbole/callisto/v4/types"
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

	codeType := types.NewWasmCode("", code.CodeBytes, &code.CodeInfo.InstantiateConfig, code.CodeID, 10)

	err := suite.database.SaveWasmCodes([]types.WasmCode{codeType})
	suite.Require().NoError(err)

	var rows []dbtypes.WasmCodeRow
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM wasm_code`)
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1)

	suite.Require().Equal(codeType.Sender, rows[0].Sender)
	suite.Require().Equal(codeType.CodeID, rows[0].CodeID)
	suite.Require().Equal(codeType.Height, rows[0].Height)

	var storedAccessConfig *wasmtypes.AccessConfig
	err = json.Unmarshal([]byte(rows[0].InstantiatePermission), &storedAccessConfig)
	suite.Require().NoError(err)
	suite.Require().Equal(codeType.InstantiatePermission, storedAccessConfig)

	return nil
}

func (suite *DbTestSuite) TestSaveWasmContracts() error {
	suite.getAccount("cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs")
	address, _ := sdk.AccAddressFromBech32("cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs")

	// Get code
	suite.getWasmCode(1, address)

	// Make a contract
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

	instantiatedAt := time.Date(2020, 10, 10, 15, 00, 00, 000, time.UTC)
	contractType := types.NewWasmContract(
		contract.ContractInfo.Creator,
		contract.ContractInfo.Admin,
		contract.ContractInfo.CodeID,
		contract.ContractInfo.Label,
		make(wasmtypes.RawContractMessage, 0),
		sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(10))),
		contract.ContractAddress,
		"",
		instantiatedAt,
		contract.ContractInfo.Creator,
		"",
		contract.ContractState,
		10,
	)

	err := suite.database.SaveWasmContracts([]types.WasmContract{contractType})
	suite.Require().NoError(err)

	// Verify the data
	dbCoins := dbtypes.NewDbCoins(sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(10))))
	expected := dbtypes.NewWasmContractRow(
		contract.ContractInfo.Creator,
		contract.ContractInfo.Admin,
		contract.ContractInfo.CodeID,
		contract.ContractInfo.Label,
		make(wasmtypes.RawContractMessage, 0),
		&dbCoins,
		contract.ContractAddress,
		"",
		instantiatedAt,
		contract.ContractInfo.Creator,
		"",
		10,
	)

	var rows []dbtypes.WasmContractRow
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM wasm_contract`)
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1)
	suite.Require().True(rows[0].Equals(expected))

	return nil
}

func (suite *DbTestSuite) TestSaveWasmExecuteContracts() error {
	suite.getAccount("cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs")
	address, _ := sdk.AccAddressFromBech32("cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs")

	// Get Code
	suite.getWasmCode(1, address)

	// Get Contract
	suite.getWasmContract(address)

	// Store contract execution
	msg := wasmtypes.MsgExecuteContract{
		Sender:   "cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs",
		Contract: "cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs",
		Msg:      make(wasmtypes.RawContractMessage, 0),
		Funds:    sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(10))),
	}
	executedAt := time.Date(2020, 10, 10, 15, 00, 00, 000, time.UTC)
	err := suite.database.SaveWasmExecuteContracts(
		[]types.WasmExecuteContract{
			types.NewWasmExecuteContract(
				msg.Sender, msg.Contract, msg.Msg, msg.Funds,
				"", executedAt, 10,
			),
		})
	suite.Require().NoError(err)

	// Verify the data
	dbCoins := dbtypes.NewDbCoins(sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(10))))
	expected := dbtypes.NewWasmExecuteContractRow(
		msg.Sender,
		msg.Contract,
		msg.Msg,
		&dbCoins,
		"",
		executedAt,
		10,
	)

	var rows []dbtypes.WasmExecuteContractRow
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM wasm_execute_contract`)
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1)
	suite.Require().True(rows[0].Equals(expected))

	return nil
}

func (suite *DbTestSuite) TestUpdateContractWithMsgMigrateContract() error {
	suite.getAccount("cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs")
	address, _ := sdk.AccAddressFromBech32("cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs")

	suite.getWasmCode(1, address)
	suite.getWasmCode(2, address)

	suite.getWasmContract(address)

	// Migrate contract to Code ID 2
	msg := &wasmtypes.MsgMigrateContract{
		Sender:   "cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs",
		Contract: "cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs",
		CodeID:   2,
		Msg:      make(wasmtypes.RawContractMessage, 0),
	}

	err := suite.database.UpdateContractWithMsgMigrateContract(msg.Sender, msg.Contract, msg.CodeID, msg.Msg, "")
	suite.Require().NoError(err)

	dbCoins := dbtypes.NewDbCoins(sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(10))))
	expected := dbtypes.NewWasmContractRow(
		"cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs",
		"cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs",
		msg.CodeID,
		"label",
		make(wasmtypes.RawContractMessage, 0),
		&dbCoins,
		"cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs",
		"",
		time.Date(2020, 10, 10, 15, 00, 00, 000, time.UTC),
		"cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs",
		"",
		10,
	)

	var rows []dbtypes.WasmContractRow
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM wasm_contract`)
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1)
	suite.Require().True(rows[0].Equals(expected))

	return nil
}

func (suite *DbTestSuite) TestUpdateContractAdmin() error {
	suite.getAccount("cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs")
	address1, _ := sdk.AccAddressFromBech32("cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs")

	suite.getAccount("cosmos1ltzt0z992ke6qgmtjxtygwzn36km4cy6cqdknt")

	suite.getWasmCode(1, address1)
	suite.getWasmContract(address1)

	msg := &wasmtypes.MsgUpdateAdmin{
		Sender:   "cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs",
		NewAdmin: "cosmos1ltzt0z992ke6qgmtjxtygwzn36km4cy6cqdknt",
		Contract: "cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs",
	}

	err := suite.database.UpdateContractAdmin(msg.Sender, msg.Contract, msg.NewAdmin)
	suite.Require().NoError(err)

	dbCoins := dbtypes.NewDbCoins(sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(10))))
	expected := dbtypes.NewWasmContractRow(
		"cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs",
		"cosmos1ltzt0z992ke6qgmtjxtygwzn36km4cy6cqdknt",
		1,
		"label",
		make(wasmtypes.RawContractMessage, 0),
		&dbCoins,
		"cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs",
		"",
		time.Date(2020, 10, 10, 15, 00, 00, 000, time.UTC),
		"cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs",
		"",
		10,
	)

	var rows []dbtypes.WasmContractRow
	err = suite.database.Sqlx.Select(&rows, `SELECT * FROM wasm_contract`)
	suite.Require().NoError(err)
	suite.Require().Len(rows, 1)
	suite.Require().True(rows[0].Equals(expected))

	return nil
}

// getWasmCode saves inside the database a x/wasm code
func (suite *DbTestSuite) getWasmCode(codeID uint64, addr sdk.AccAddress) {
	codeInfo := wasmtypes.NewCodeInfo(nil, addr, wasmtypes.AllowEverybody)
	err := suite.database.SaveWasmCodes(
		[]types.WasmCode{types.NewWasmCode("", nil, &codeInfo.InstantiateConfig, codeID, 10)},
	)
	suite.Require().NoError(err)
}

// getWasmCode saves inside the database a x/wasm contract
func (suite *DbTestSuite) getWasmContract(address sdk.AccAddress) {
	suite.getWasmCode(1, address)

	err := suite.database.SaveWasmContracts([]types.WasmContract{types.NewWasmContract(
		"cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs",
		"cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs",
		1,
		"label",
		make(wasmtypes.RawContractMessage, 0),
		sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(10))),
		"cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs",
		"",
		time.Date(2020, 10, 10, 15, 00, 00, 000, time.UTC),
		"cosmos1z4hfrxvlgl4s8u4n5ngjcw8kdqrcv43599amxs",
		"",
		[]wasmtypes.Model{},
		10,
	)})

	suite.Require().NoError(err)
}
