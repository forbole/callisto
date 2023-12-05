package types

import (
	"database/sql/driver"
	"fmt"
	"regexp"
	"strings"

	"git.ooo.ua/vipcoin/lib/errs"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/lib/pq"
)

var (
	coinsRegEx    = regexp.MustCompile(`\(([^,]+),(\d+)\)`)
	sendDataRegEx = regexp.MustCompile(`^\(([^,]+),\{(.+)\}\)$`)
)

type (
	// DbSendData represents the information stored inside the database about a MsgMultiSend or other types.
	DbSendData struct {
		Address string
		Coins   DbCoins
	}

	// DbSendDataList represents a list of DbSendData
	DbSendDataList []*DbSendData
)

// NewDbSendData builds a DbSendData.
func NewDbSendData(address string, coins DbCoins) *DbSendData {
	return &DbSendData{
		Address: address,
		Coins:   coins,
	}
}

// NewDbSendDataByInput builds a DbSendData by an SDK Input
func NewDbSendDataByInput(in types.Input) *DbSendData {
	return &DbSendData{
		Address: in.Address,
		Coins:   NewDbCoins(in.Coins),
	}
}

// NewDbSendDataListByInputs builds a DbSendData by a list of SDK Input
func NewDbSendDataListByInputs(list []types.Input) []*DbSendData {
	res := make(DbSendDataList, 0, len(list))

	for _, in := range list {
		res = append(res, NewDbSendDataByInput(in))
	}

	return res
}

// NewDbSendDataByOutput builds a DbSendData by a list of SDK Output
func NewDbSendDataByOutput(out types.Output) *DbSendData {
	return &DbSendData{
		Address: out.Address,
		Coins:   NewDbCoins(out.Coins),
	}
}

// NewDbSendDataListByOutputs builds a DbSendData by a list of SDK Output
func NewDbSendDataListByOutputs(list []types.Output) []*DbSendData {
	res := make(DbSendDataList, 0, len(list))

	for _, out := range list {
		res = append(res, NewDbSendDataByOutput(out))
	}

	return res
}

// ToInput - mapping func to a domain model.
func (d DbSendData) ToInput() types.Input {
	return types.Input{
		Address: d.Address,
		Coins:   d.Coins.ToCoins(),
	}
}

// ToInputList - mapping func to a domain model.
func (d DbSendDataList) ToInputList() []types.Input {
	res := make([]types.Input, 0, len(d))

	for _, in := range d {
		res = append(res, in.ToInput())
	}

	return res
}

// ToOutput - mapping func to a domain model.
func (d DbSendData) ToOutput() types.Output {
	return types.Output{
		Address: d.Address,
		Coins:   d.Coins.ToCoins(),
	}
}

// ToOutputList - mapping func to a domain model.
func (d DbSendDataList) ToOutputList() []types.Output {
	res := make([]types.Output, 0, len(d))

	for _, out := range d {
		res = append(res, out.ToOutput())
	}

	return res
}

// FromPqStringArrayToInputs - mapping func to a domain model.
// Example arg: (ovg18p9heyy3m4dsq7fj86p6v9yzzx8a64f86eec7u,{(ovg,5),(stovg,2)})
func FromPqStringArrayToInputs(arr pq.StringArray) ([]types.Input, error) {
	res := make([]types.Input, 0, len(arr))

	for _, v := range arr {
		match := sendDataRegEx.FindStringSubmatch(v) // e.g.: (ovg18p9hey...,{(ovg,5),(stovg,2)})
		if match != nil {
			address := match[1] // e.g.: ovg18p9hey...

			// parse coins
			coins := make(sdk.Coins, 0)
			for _, coinStr := range strings.Split(match[2], ",") { // e.g.: {(ovg,5),(stovg,2)})
				coin, err := parseCoin(coinStr) // e.g.: (ovg,5)
				if err != nil {
					return nil, err
				}

				coins = append(coins, coin)
			}

			res = append(res, types.Input{
				Address: address,
				Coins:   coins,
			})
		}
	}

	return res, nil
}

// FromPqStringArrayToOutputs - mapping func to a domain model.
// Example arg: (ovg18p9heyy3m4dsq7fj86p6v9yzzx8a64f86eec7u,{(ovg,5),(stovg,2)})
func FromPqStringArrayToOutputs(arr pq.StringArray) ([]types.Output, error) {
	res := make([]types.Output, 0, len(arr))

	for _, v := range arr {
		match := sendDataRegEx.FindStringSubmatch(v) // e.g.: (ovg18p9hey...,{(ovg,5),(stovg,2)})
		if match != nil {
			address := match[1] // e.g.: ovg18p9hey...

			// parse coins
			coins := make(sdk.Coins, 0)
			for _, coinStr := range strings.Split(match[2], ",") { // e.g.: {(ovg,5),(stovg,2)})
				coin, err := parseCoin(coinStr) // e.g.: (ovg,5)
				if err != nil {
					return nil, err
				}

				coins = append(coins, coin)
			}

			res = append(res, types.Output{
				Address: address,
				Coins:   coins,
			})
		}
	}

	return res, nil
}

// Value implements driver.Valuer
func (d *DbSendData) Value() (driver.Value, error) {
	coins, err := d.Coins.Value()
	if err != nil {
		return nil, err
	}

	return fmt.Sprintf("(%s,%s)", d.Address, coins), nil
}

// Scan implements sql.Scanner
func (coin *DbSendData) Scan(src interface{}) error {
	strValue := string(src.([]byte))
	strValue = strings.ReplaceAll(strValue, `"`, "")
	strValue = strings.ReplaceAll(strValue, "{", "")
	strValue = strings.ReplaceAll(strValue, "}", "")
	strValue = strings.ReplaceAll(strValue, "(", "")
	strValue = strings.ReplaceAll(strValue, ")", "")

	values := strings.Split(strValue, ",")

	address := values[0]
	coins, err := sdk.ParseCoinsNormalized(values[1])
	if err != nil {
		return err
	}

	*coin = DbSendData{Address: address, Coins: NewDbCoins(coins)}
	return nil
}

// Value implements driver.Valuer, e.g.: {(ovg,100000000),(stovg,500000000000)}
func (coins *DbCoins) Value() (driver.Value, error) {
	res := make(pq.StringArray, 0, len(coins.ToCoins()))
	for _, coin := range coins.ToCoins() {
		res = append(res, fmt.Sprintf("(%s,%s)", coin.Denom, coin.Amount.String()))
	}

	return fmt.Sprintf("{%s}", strings.Join(res, ",")), nil
}

// parseCoin - parses a single coin, e.g.: (ovg,100000000)
func parseCoin(coinStr string) (sdk.Coin, error) {
	match := coinsRegEx.FindStringSubmatch(coinStr)
	if match != nil {
		denom := match[1]
		amountStr := match[2]

		amount, ok := sdk.NewIntFromString(amountStr)
		if !ok {
			return sdk.Coin{}, errs.Internal{Cause: "invalid amount"}
		}

		return sdk.NewCoin(denom, amount), nil
	}

	return sdk.Coin{}, errs.Internal{Cause: "invalid coin"}
}

// FromPqStringArrayToCoins - converts pq.StringArray to sdk.Coins, e.g.: {(ovg,10),(stovg,5)}
func FromPqStringArrayToCoins(arr pq.StringArray) (sdk.Coins, error) {
	res := make(sdk.Coins, 0)
	for _, v := range arr {
		coin, err := parseCoin(v)
		if err != nil {
			return nil, err
		}

		res = append(res, coin)
	}

	return res, nil
}
