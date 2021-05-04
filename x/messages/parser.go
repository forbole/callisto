package messages

import (
	junomessages "github.com/desmos-labs/juno/modules/messages"
)

// AddressesParser represents a MessageAddressesParser able to parse Cosmos and custom chain messages
var AddressesParser = junomessages.JoinMessageParsers(
	junomessages.CosmosMessageAddressesParser,
	desmosMessageAddressesParser,
)
