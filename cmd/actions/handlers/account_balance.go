package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	actionstypes "github.com/forbole/bdjuno/v2/cmd/actions/types"
	"github.com/forbole/bdjuno/v2/utils"
)

func AccountBalance(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	var actionPayload actionstypes.Payload
	err = json.Unmarshal(reqBody, &actionPayload)
	if err != nil {
		http.Error(w, "invalid payload: failed to unmarshal json", http.StatusInternalServerError)
		return
	}

	result, err := getAccountBalance(actionPayload.Input)
	if err != nil {
		errorHandler(w, err)
		return
	}

	data, _ := json.Marshal(result)
	w.Write(data)
}

func getAccountBalance(input actionstypes.PayloadArgs) (response actionstypes.Balance, err error) {
	parseCtx, sources, err := getCtxAndSources()
	if err != nil {
		return response, err
	}

	height, err := utils.GetHeight(parseCtx, input.Height)
	if err != nil {
		return response, fmt.Errorf("error while getting height: %s", err)
	}

	balance, err := sources.BankSource.GetAccountBalance(input.Address, height)
	if err != nil {
		return response, err
	}

	return actionstypes.Balance{
		Coins: balance,
	}, nil
}
