package margin

import (
	"fmt"

	margintypes "github.com/Sifchain/sifnode/x/margin/types"
	"github.com/forbole/bdjuno/v3/types"
	juno "github.com/forbole/juno/v3/types"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
)

// HandleBlock implements BlockModule
func (m *Module) HandleBlock(
	block *tmctypes.ResultBlock, results *tmctypes.ResultBlockResults, tx []*juno.Tx, _ *tmctypes.ResultValidators,
) error {
	// Get x/margin module events
	err := m.getMarginEvents(block.Block.Height, tx)
	if err != nil {
		return fmt.Errorf("error while getting x/margin events: %s", err)
	}

	return nil
}

// getMarginEvents reads the events from txs and stores its values inside the database
func (m *Module) getMarginEvents(height int64, tx []*juno.Tx) error {
	var events []types.MarginEvent
	var involvedAccounts []string

	for _, e := range tx {
		for _, event := range e.Events {
			address, _ := juno.FindAttributeByKey(event, "address")
			if len(address.String()) > 0 {
				involvedAccounts = append(involvedAccounts, address.String())
			}
			switch event.Type {
			case margintypes.EventForceClose:
				events = append(events, *types.NewMarginEvent(e.TxHash, 0, event.Type, event, involvedAccounts, height))
			case margintypes.EventOpen:
				events = append(events, *types.NewMarginEvent(e.TxHash, 0, event.Type, event, involvedAccounts, height))
			case margintypes.EventInterestRateComputation:
				events = append(events, *types.NewMarginEvent(e.TxHash, 0, event.Type, event, involvedAccounts, height))
			case margintypes.EventRepayFund:
				events = append(events, *types.NewMarginEvent(e.TxHash, 0, event.Type, event, involvedAccounts, height))
			case margintypes.EventBelowRemovalThreshold:
				events = append(events, *types.NewMarginEvent(e.TxHash, 0, event.Type, event, involvedAccounts, height))
			case margintypes.EventAboveRemovalThreshold:
				events = append(events, *types.NewMarginEvent(e.TxHash, 0, event.Type, event, involvedAccounts, height))
			case margintypes.EventIncrementalPayFund:
				events = append(events, *types.NewMarginEvent(e.TxHash, 0, event.Type, event, involvedAccounts, height))
			case margintypes.EventMarginUpdateParams:
				events = append(events, *types.NewMarginEvent(e.TxHash, 0, event.Type, event, involvedAccounts, height))
			}
		}

	}

	return m.db.SaveMarginEvent(events)
}
