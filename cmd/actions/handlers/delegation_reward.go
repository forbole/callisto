package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	actionstypes "github.com/forbole/bdjuno/v2/cmd/actions/types"
)

func DelegationReward(w http.ResponseWriter, r *http.Request) {

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

	result, err := getDelegationReward(actionPayload.Input.Address)
	if err != nil {
		graphQLError(w, err)
		return
	}

	data, _ := json.Marshal(result)
	w.Write(data)
}

func getDelegationReward(address string) (response []actionstypes.DelegatorReward, err error) {
	parseCtx, sources, err := getCtxAndSources()
	if err != nil {
		return response, err
	}

	// Get latest node height
	height, err := parseCtx.Node.LatestHeight()
	if err != nil {
		return response, fmt.Errorf("error while getting chain latest block height: %s", err)
	}

	// Get delegator's total rewards
	rewards, err := sources.DistrSource.DelegatorTotalRewards(address, height)
	if err != nil {
		return response, err
	}

	response = make([]actionstypes.DelegatorReward, len(rewards))
	for index, rew := range rewards {
		response[index] = actionstypes.DelegatorReward{
			DecCoins:   rew.Reward,
			ValAddress: rew.ValidatorAddress,
		}
	}

	return response, nil
}
