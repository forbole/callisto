package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	sdk "github.com/cosmos/cosmos-sdk/types"
	actionstypes "github.com/forbole/bdjuno/v2/cmd/actions/types"
)

func TotalDelegationAmount(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	var actionPayload actionstypes.StakingPayload
	err = json.Unmarshal(reqBody, &actionPayload)
	if err != nil {
		http.Error(w, "invalid payload: failed to unmarshal json", http.StatusInternalServerError)
		return
	}

	result, err := getTotalDelegationAmount(actionPayload.Input)
	if err != nil {
		errorHandler(w, err)
		return
	}

	data, _ := json.Marshal(result)
	w.Write(data)
}

func getTotalDelegationAmount(input actionstypes.StakingArgs) (actionstypes.Balance, error) {
	parseCtx, sources, err := getCtxAndSources()
	if err != nil {
		return actionstypes.Balance{}, err
	}

	// Get latest node height
	height, err := parseCtx.Node.LatestHeight()
	if err != nil {
		return actionstypes.Balance{}, fmt.Errorf("error while getting chain latest block height: %s", err)
	}

	// Get all  delegations for given delegator address
	delegationList, err := sources.StakingSource.GetDelegationsWithPagination(height, input.Address, nil)
	if err != nil {
		return actionstypes.Balance{}, fmt.Errorf("error while getting delegator delegations: %s", err)
	}

	var coinObject sdk.Coins

	// Add up total value of delegations
	for _, eachDelegation := range delegationList.DelegationResponses {
		for index, eachCoin := range coinObject {
			if eachCoin.Denom == eachDelegation.Balance.Denom {
				coinObject[index].Amount = coinObject[index].Amount.Add(eachDelegation.Balance.Amount)
			}
			if eachCoin.Denom != eachDelegation.Balance.Denom {
				coinObject = append(coinObject, sdk.NewCoin(eachDelegation.Balance.Denom, eachDelegation.Balance.Amount))
			}
		}
		if coinObject == nil {
			coinObject = append(coinObject, sdk.NewCoin(eachDelegation.Balance.Denom, eachDelegation.Balance.Amount))
		}
	}

	return actionstypes.Balance{
		Coins: coinObject,
	}, nil
}
