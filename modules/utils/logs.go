package utils

import (
	"fmt"
	"strconv"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func GetValueFromLogs(index uint32, logs sdk.ABCIMessageLogs, eventType, attributeKey string) string {
	for _, log := range logs {
		if log.MsgIndex != index {
			continue
		}

		for _, event := range log.Events {
			if event.Type != eventType {
				continue
			}

			for _, attr := range event.Attributes {
				if attr.Key == attributeKey {
					return strings.ReplaceAll(attr.Value, "\"", "")
				}
			}
		}
	}

	return ""
}

func GetUint64FromLogs(index int, logs sdk.ABCIMessageLogs, txHash, eventType, attributeKey string) (uint64, error) {
	valueStr := GetValueFromLogs(uint32(index), logs, eventType, attributeKey)
	if valueStr == "" {
		return 0, fmt.Errorf("attribute %s for event %s not found in tx %s", attributeKey, eventType, txHash)
	}

	value, err := strconv.ParseUint(valueStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse %s from tx %s to uint64", valueStr, txHash)
	}

	return value, nil
}
