package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/types/query"
	actionstypes "github.com/forbole/bdjuno/v2/cmd/actions/types"
	dbtypes "github.com/forbole/bdjuno/v2/database/types"

	"github.com/forbole/bdjuno/v2/utils"
)

func Delegation(w http.ResponseWriter, r *http.Request) {

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

	result, err := getDelegation(actionPayload.Input)
	if err != nil {
		errorHandler(w, err)
		return
	}

	data, _ := json.Marshal(result)
	w.Write(data)
}

func getDelegation(input actionstypes.PayloadArgs) (actionstypes.DelegationResponse, error) {
	parseCtx, sources, err := getCtxAndSources()
	if err != nil {
		return actionstypes.DelegationResponse{}, err
	}

	height, err := utils.GetHeight(parseCtx, input.Height)
	if err != nil {
		return actionstypes.DelegationResponse{}, fmt.Errorf("error while getting height: %s", err)
	}

	pagination := &query.PageRequest{
		Offset:     input.Offset,
		Limit:      input.Limit,
		CountTotal: input.CountTotal,
	}

	// Get delegator's total rewards
	res, err := sources.StakingSource.GetDelegationsWithPagination(height, input.Address, pagination)
	if err != nil {
		return actionstypes.DelegationResponse{}, fmt.Errorf("error while getting delegator delegations: %s", err)
	}

	delegations := make([]actionstypes.Delegation, len(res.DelegationResponses))
	for index, del := range res.DelegationResponses {
		delegations[index] = actionstypes.Delegation{
			DelegatorAddress: del.Delegation.DelegatorAddress,
			ValidatorAddress: del.Delegation.ValidatorAddress,
			Coins:            dbtypes.NewDbCoins([]sdk.Coin{del.Balance}),
		}
	}

	return actionstypes.DelegationResponse{
		Delegations: delegations,
		Pagination:  res.Pagination,
	}, nil
}
