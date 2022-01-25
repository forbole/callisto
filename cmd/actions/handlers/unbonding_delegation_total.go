package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	sdk "github.com/cosmos/cosmos-sdk/types"
	actionstypes "github.com/forbole/bdjuno/v2/cmd/actions/types"
	"github.com/forbole/bdjuno/v2/utils"
)

func UnbondingDelegationsTotal(w http.ResponseWriter, r *http.Request) {

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

	result, err := getUnbondingDelegationsTotalAmount(actionPayload.Input)
	if err != nil {
		errorHandler(w, err)
		return
	}

	data, _ := json.Marshal(result)
	w.Write(data)
}

func getUnbondingDelegationsTotalAmount(input actionstypes.PayloadArgs) (actionstypes.Balance, error) {
	parseCtx, sources, err := getCtxAndSources()
	if err != nil {
		return actionstypes.Balance{}, err
	}

	height, err := utils.GetHeight(parseCtx, input.Height)
	if err != nil {
		return actionstypes.Balance{}, fmt.Errorf("error while getting chain latest block height: %s", err)
	}

	// Get all unbonding delegations for given delegator address
	unbondingDelegations, err := sources.StakingSource.GetUnbondingDelegations(height, input.Address, nil)
	if err != nil {
		return actionstypes.Balance{}, fmt.Errorf("error while getting delegator delegations: %s", err)
	}

	var coins sdk.Coins
	var totalAmount int64

	// Get the bond denom type
	params, err := sources.StakingSource.GetParams(height)
	if err != nil {
		return actionstypes.Balance{}, fmt.Errorf("error while getting bond denom type: %s", err)
	}

	// Add up total value of unbonding delegations
	for _, eachUnbondingDelegation := range unbondingDelegations.UnbondingResponses {
		for _, entry := range eachUnbondingDelegation.Entries {
			totalAmount += entry.Balance.Int64()
		}
	}

	coins = append(coins, sdk.NewCoin(params.BondDenom, sdk.NewInt(totalAmount)))

	return actionstypes.Balance{
		Coins: coins,
	}, nil
}
