package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	sdk "github.com/cosmos/cosmos-sdk/types"
	actionstypes "github.com/forbole/bdjuno/v2/cmd/actions/types"
)

func AccountBalance(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	var actionPayload actionstypes.AccountBalancesPayload
	err = json.Unmarshal(reqBody, &actionPayload)
	if err != nil {
		http.Error(w, "invalid payload: failed to unmarshal json", http.StatusInternalServerError)
		return
	}

	result, err := getAccountBalance(actionPayload.Input)
	if err != nil {
		graphQLError(w, err)
		return
	}

	data, _ := json.Marshal(result)
	w.Write(data)
}

func getAccountBalance(input actionstypes.AccountBalancesArgs) (actionstypes.Coins, error) {
	parseCtx, sources, err := getCtxAndSources()
	if err != nil {
		return actionstypes.Coins{}, err
	}

	height := input.Height
	if height == 0 {
		// Get latest height if height input is empty
		height, err = parseCtx.Node.LatestHeight()
		if err != nil {
			return actionstypes.Coins{}, fmt.Errorf("error while getting chain latest block height: %s", err)
		}
	}

	balances, err := sources.BankSource.GetBalances([]string{input.Address}, height)

	if err != nil {
		return actionstypes.Coins{}, err
	}

	var coins []sdk.Coin
	for _, bal := range balances {
		for _, coin := range bal.Balance {
			coins = append(coins, coin)
		}
	}

	return actionstypes.Coins{
		Coins: coins,
	}, nil
}
