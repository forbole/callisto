package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	actionstypes "github.com/forbole/bdjuno/v2/cmd/actions/types"
	"github.com/forbole/bdjuno/v2/utils"
)

func DelegationReward(w http.ResponseWriter, r *http.Request) {

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

	result, err := getDelegationReward(actionPayload.Input)
	if err != nil {
		errorHandler(w, err)
		return
	}

	data, _ := json.Marshal(result)
	w.Write(data)
}

func getDelegationReward(input actionstypes.PayloadArgs) (response []actionstypes.DelegationReward, err error) {
	parseCtx, sources, err := getCtxAndSources()
	if err != nil {
		return response, err
	}

	height, err := utils.GetHeight(parseCtx, input.Height)
	if err != nil {
		return response, fmt.Errorf("error while getting chain latest block height: %s", err)
	}

	// Get delegator's total rewards
	rewards, err := sources.DistrSource.DelegatorTotalRewards(input.Address, height)
	if err != nil {
		return response, err
	}

	delegationRewards := make([]actionstypes.DelegationReward, len(rewards))
	for index, rew := range rewards {
		delegationRewards[index] = actionstypes.DelegationReward{
			Coins:            rew.Reward,
			ValidatorAddress: rew.ValidatorAddress,
		}
	}

	return delegationRewards, nil
}
