package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	actionstypes "github.com/forbole/bdjuno/v2/cmd/actions/types"
)

func Redelegation(w http.ResponseWriter, r *http.Request) {

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

	result, err := getRedelegation(actionPayload.Input.Address)
	if err != nil {
		errorHandler(w, err)
		return
	}

	data, _ := json.Marshal(result)
	w.Write(data)
}

func getRedelegation(address string) ([]actionstypes.Redelegation, error) {
	parseCtx, sources, err := getCtxAndSources()
	if err != nil {
		return nil, err
	}

	// Get latest node height
	height, err := parseCtx.Node.LatestHeight()
	if err != nil {
		return nil, fmt.Errorf("error while getting chain latest block height: %s", err)
	}

	// Get delegator's redelegations
	redelegations, err := sources.StakingSource.GetDelegatorRedelegations(height, address)
	if err != nil {
		return nil, fmt.Errorf("error while getting delegator redelegations: %s", err)
	}

	response := make([]actionstypes.Redelegation, len(redelegations))
	for index, redel := range redelegations {
		response[index] = actionstypes.Redelegation{
			DelegatorAddress:    redel.Redelegation.DelegatorAddress,
			ValidatorSrcAddress: redel.Redelegation.ValidatorSrcAddress,
			ValidatorDstAddress: redel.Redelegation.ValidatorSrcAddress,
			Entries:             redel.Redelegation.Entries,
		}
	}

	return response, nil
}
