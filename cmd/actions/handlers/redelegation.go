package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/cosmos/cosmos-sdk/types/query"
	actionstypes "github.com/forbole/bdjuno/v2/cmd/actions/types"
)

func Redelegation(w http.ResponseWriter, r *http.Request) {

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

	result, err := getRedelegation(actionPayload.Input)
	if err != nil {
		errorHandler(w, err)
		return
	}

	data, _ := json.Marshal(result)
	w.Write(data)
}

func getRedelegation(input actionstypes.StakingArgs) (actionstypes.RedelegationResponse, error) {
	parseCtx, sources, err := getCtxAndSources()
	if err != nil {
		return actionstypes.RedelegationResponse{}, err
	}

	// Get latest node height
	height, err := parseCtx.Node.LatestHeight()
	if err != nil {
		return actionstypes.RedelegationResponse{}, fmt.Errorf("error while getting chain latest block height: %s", err)
	}

	pagination := &query.PageRequest{
		Offset:     input.Offset,
		Limit:      input.Limit,
		CountTotal: input.CountTotal,
	}

	// Get delegator's redelegations
	redelegations, err := sources.StakingSource.GetDelegatorRedelegations(height, input.Address, pagination)
	if err != nil {
		return actionstypes.RedelegationResponse{}, fmt.Errorf("error while getting delegator redelegations: %s", err)
	}

	redelegationsList := make([]actionstypes.Redelegation, len(redelegations.RedelegationResponses))
	for index, del := range redelegations.RedelegationResponses {
		redelegationsList[index] = actionstypes.Redelegation{
			DelegatorAddress:    del.Redelegation.DelegatorAddress,
			ValidatorSrcAddress: del.Redelegation.ValidatorSrcAddress,
			ValidatorDstAddress: del.Redelegation.ValidatorDstAddress,
		}

		for indexEntry, entry := range del.Entries {
			redelegationsList[index].
				RedelegationEntries[indexEntry] = actionstypes.RedelegationEntry{
				CompletionTime: entry.RedelegationEntry.CompletionTime,
				Balance:        entry.Balance,
			}
		}
	}

	return actionstypes.RedelegationResponse{
		Redelegations: redelegationsList,
		Pagination:    redelegations.Pagination,
	}, nil
}
