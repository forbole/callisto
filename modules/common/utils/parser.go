package utils

import (
	junomessages "github.com/desmos-labs/juno/modules/messages"
)

// AddressesParser represents a MessageAddressesParser able to parse Cosmos and custom chain utils
var AddressesParser = junomessages.JoinMessageParsers(
	junomessages.CosmosMessageAddressesParser,
	desmosMessageAddressesParser,
)
