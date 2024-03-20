package gov

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	eventsutil "github.com/forbole/bdjuno/v4/utils/events"
)

// ProposalIDFromEvent returns the proposal id from the given events
func ProposalIDFromEvents(events sdk.StringEvents) (uint64, error) {
	for _, event := range events {
		attribute, ok := eventsutil.FindAttributeByKey(event, govtypes.AttributeKeyProposalID)
		if ok {
			return strconv.ParseUint(attribute.Value, 10, 64)
		}
	}

	return 0, fmt.Errorf("no proposal id found")
}

// VoteOptionFromEvents returns the vote option from the given events
func VoteOptionFromEvents(events sdk.StringEvents) (govtypesv1.VoteOption, error) {
	for _, event := range events {
		attribute, ok := eventsutil.FindAttributeByKey(event, govtypes.AttributeKeyOption)
		if ok {
			return parseVoteOption(attribute.Value)
		}
	}

	return 0, fmt.Errorf("no vote option found")
}

// parseVoteOption returns the vote option from the given string
// option value in string could be 2 cases, for example:
// 1. "{\"option\":1,\"weight\":\"1.000000000000000000\"}"
// 2. "option:VOTE_OPTION_NO weight:\"1.000000000000000000\""
func parseVoteOption(optionValue string) (govtypesv1.VoteOption, error) {
	// try parse option value as json
	type voteOptionJSON struct {
		Option govtypesv1.VoteOption `json:"option"`
	}
	var voteOptionParsedJSON voteOptionJSON
	err := json.Unmarshal([]byte(optionValue), &voteOptionParsedJSON)
	if err == nil {
		return voteOptionParsedJSON.Option, nil
	}

	// try parse option value as string
	// option:VOTE_OPTION_NO weight:"1.000000000000000000"
	voteOptionParsed := strings.Split(optionValue, " ")
	voteOption, err := govtypesv1.VoteOptionFromString(strings.ReplaceAll(voteOptionParsed[0], "option:", ""))
	if err != nil {
		return 0, fmt.Errorf("failed to parse vote option %s: %s", optionValue, err)
	}

	return voteOption, nil
}
