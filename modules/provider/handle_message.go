package provider

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	junotypes "github.com/forbole/juno/v3/types"
	providertypes "github.com/ovrclk/akash/x/provider/types/v1beta2"
)

// HandleMsg implements modules.MessageModule
func (m *Module) HandleMsg(_ int, msg sdk.Msg, tx *junotypes.Tx) error {
	switch cosmosMsg := msg.(type) {
	case *providertypes.MsgCreateProvider:
		return m.db.SaveProviders([]providertypes.Provider{
			{
				Owner:      cosmosMsg.Owner,
				HostURI:    cosmosMsg.HostURI,
				Attributes: cosmosMsg.Attributes,
				Info:       cosmosMsg.Info,
			},
		}, tx.Height)
	case *providertypes.MsgDeleteProvider:
		return m.db.DeleteProvider(cosmosMsg.Owner)
	case *providertypes.MsgUpdateProvider:
		return m.db.SaveProviders([]providertypes.Provider{
			{
				Owner:      cosmosMsg.Owner,
				HostURI:    cosmosMsg.HostURI,
				Attributes: cosmosMsg.Attributes,
				Info:       cosmosMsg.Info,
			},
		}, tx.Height)
	}

	return nil
}