package provider

import (
	providertypes "github.com/akash-network/akash-api/go/node/provider/v1beta3"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/forbole/bdjuno/v4/types"
	junotypes "github.com/forbole/juno/v4/types"
)

// HandleMsg implements modules.MessageModule
func (m *Module) HandleMsg(_ int, msg sdk.Msg, tx *junotypes.Tx) error {
	switch cosmosMsg := msg.(type) {
	case *providertypes.MsgCreateProvider:
		return m.db.SaveProviders([]*types.Provider{
			{
				OwnerAddress: cosmosMsg.Owner,
				HostURI:      cosmosMsg.HostURI,
				Attributes:   cosmosMsg.Attributes,
				Info:         cosmosMsg.Info,
			},
		}, tx.Height)
	case *providertypes.MsgDeleteProvider:
		return m.db.DeleteProvider(cosmosMsg.Owner)
	case *providertypes.MsgUpdateProvider:
		return m.db.SaveProviders([]*types.Provider{
			{
				OwnerAddress: cosmosMsg.Owner,
				HostURI:      cosmosMsg.HostURI,
				Attributes:   cosmosMsg.Attributes,
				Info:         cosmosMsg.Info,
			},
		}, tx.Height)
	}

	return nil
}
