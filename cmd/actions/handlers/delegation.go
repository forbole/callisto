package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	actionstypes "github.com/forbole/bdjuno/v2/cmd/actions/types"
)

func Delegation(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	var actionPayload actionstypes.AddressPayload
	err = json.Unmarshal(reqBody, &actionPayload)
	if err != nil {
		http.Error(w, "invalid payload: failed to unmarshal json", http.StatusInternalServerError)
		return
	}

	result, err := getDelegation(actionPayload.Input.Address)
	if err != nil {
		errorHandler(w, err)
		return
	}

	data, _ := json.Marshal(result)
	w.Write(data)
}

func getDelegation(address string) ([]actionstypes.Delegation, error) {
	parseCtx, sources, err := getCtxAndSources()
	if err != nil {
		return nil, err
	}

	// Get latest node height
	height, err := parseCtx.Node.LatestHeight()
	if err != nil {
		return nil, fmt.Errorf("error while getting chain latest block height: %s", err)
	}

	// Get delegator's total rewards
	delegations, err := sources.StakingSource.GetDelegatorDelegations(height, address)
	if err != nil {
		return nil, fmt.Errorf("error while getting delegator delegations: %s", err)
	}

	response := make([]actionstypes.Delegation, len(delegations))
	for index, del := range delegations {
		response[index] = actionstypes.Delegation{
			DelAddress: del.Delegation.DelegatorAddress,
			ValAddress: del.Delegation.ValidatorAddress,
			Coin:       del.Balance,
		}
	}

	return response, nil
}
