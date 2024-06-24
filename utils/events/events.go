package events

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// FindEventByType returns the event with the given type
func FindEventByType(events sdk.StringEvents, eventType string) (sdk.StringEvent, bool) {
	for _, event := range events {
		if event.Type == eventType {
			return event, true
		}
	}
	return sdk.StringEvent{}, false
}

// FindAttributeByKey returns the attribute with the given key
func FindAttributeByKey(event sdk.StringEvent, key string) (sdk.Attribute, bool) {
	for _, attribute := range event.Attributes {
		if attribute.Key == key {
			return attribute, true
		}
	}
	return sdk.Attribute{}, false
}

// FindEventsByMsgIndex returns all events with the given msg index
func FindEventsByMsgIndex(events sdk.StringEvents, msgIndex int) sdk.StringEvents {
	var res sdk.StringEvents
	for _, event := range events {
		attribute, exist := FindAttributeByKey(event, "msg_index")
		if !exist {
			continue
		}

		if strconv.Itoa(msgIndex) == attribute.Value {
			res = append(res, event)
		}
	}
	return res
}
