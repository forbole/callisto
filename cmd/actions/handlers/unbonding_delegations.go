package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	actionstypes "github.com/forbole/bdjuno/v2/cmd/actions/types"
)

func UnbondingDelegations(w http.ResponseWriter, r *http.Request) {

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

	result, err := getUnbondingDelegations(actionPayload.Input.Address)
	if err != nil {
		errorHandler(w, err)
		return
	}

	data, _ := json.Marshal(result)
	w.Write(data)
}

func getUnbondingDelegations(address string) ([]actionstypes.UnbondingDelegation, error) {
	parseCtx, sources, err := getCtxAndSources()
	if err != nil {
		return nil, err
	}

	// Get latest node height
	height, err := parseCtx.Node.LatestHeight()
	if err != nil {
		return nil, fmt.Errorf("error while getting chain latest block height: %s", err)
	}

	// Get all unbonding delegations for given delegator address
	delegations, err := sources.StakingSource.GetUnbondingDelegations(height, address)
	if err != nil {
		return nil, fmt.Errorf("error while getting unbonding delegations for address %s: %s", address, err)
	}

	response := make([]actionstypes.UnbondingDelegation, len(delegations))
	for index, delegation := range delegations {
		response[index] = actionstypes.UnbondingDelegation{
			DelegatorAddress: delegation.DelegatorAddress,
			ValidatorAddress: delegation.ValidatorAddress,
			Entries: delegation.Entries,
		}
	}

	return response, nil
}
