package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	actionstypes "github.com/forbole/bdjuno/v2/cmd/actions/types"
)

func AccountBalance(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	var actionPayload actionstypes.AccountBalancePayload
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

func getAccountBalance(input actionstypes.AccountBalancesArgs) (response []actionstypes.Coin, err error) {
	parseCtx, sources, err := getCtxAndSources()
	if err != nil {
		return response, err
	}

	height := input.Height
	if height == 0 {
		// Get latest height if height input is empty
		height, err = parseCtx.Node.LatestHeight()
		if err != nil {
			return response, fmt.Errorf("error while getting chain latest block height: %s", err)
		}
	}

	balances, err := sources.BankSource.GetBalances([]string{input.Address}, height)

	if err != nil {
		return response, err
	}

	for _, bal := range balances {
		for _, coin := range bal.Balance {
			response = append(response, actionstypes.Coin{
				Amount: coin.Amount.Int64(),
				Denom:  coin.Denom,
			})
		}
	}

	return response, nil
}
