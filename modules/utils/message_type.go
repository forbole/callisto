package utils

import (
	"fmt"
	"strings"
)

func GetModuleNameFromTypeURL(input string) string {
	moduleName := strings.Split(input, ".")
	if len(moduleName) > 1 {
		switch {
		case strings.Contains(moduleName[0], "cosmos"):
			return moduleName[1] // e.g. "cosmos.bank.v1beta1.MsgSend" => "bank"
		case strings.Contains(moduleName[0], "ibc"):
			return fmt.Sprintf("%s %s %s", moduleName[0], moduleName[1], moduleName[2]) // e.g. "ibc.core.client.v1.MsgUpdateClient" => "ibc core client"
		default:
			return fmt.Sprintf("%s %s", moduleName[0], moduleName[1]) // e.g. "cosmwasm.wasm.v1.MsgExecuteContract" => "cosmwasm wasm"
		}
	}

	return ""
}

func GetMsgFromTypeURL(input string) string {
	messageName := strings.Split(input, ".")
	if len(messageName) > 1 {
		return messageName[len(messageName)-1] // e.g. "cosmos.bank.v1beta1.MsgSend" => "MsgSend"
	}
	return ""
}
