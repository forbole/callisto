package margin

import (
	"fmt"

	margintypes "github.com/Sifchain/sifnode/x/margin/types"
	juno "github.com/forbole/juno/v3/types"
	abci "github.com/tendermint/tendermint/abci/types"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
)

// HandleBlock implements BlockModule
func (m *Module) HandleBlock(
	block *tmctypes.ResultBlock, results *tmctypes.ResultBlockResults, tx []*juno.Tx, _ *tmctypes.ResultValidators,
) error {
	// Get x/margin module events
	err := m.getMarginEvents(block.Block.Height, results.EndBlockEvents, results.BeginBlockEvents, tx)
	if err != nil {
		return fmt.Errorf("error while getting x/margin events: %s", err)
	}

	return nil
}

// getMarginEvents reads the events from txs and stores its values inside the database
func (m *Module) getMarginEvents(height int64, events []abci.Event, eventsBegin []abci.Event, tx []*juno.Tx) error {
	var msgs []*juno.Message
	var involvedAccounts []string

	for _, e := range tx {
		for _, ev := range e.Events {
			address, _ := juno.FindAttributeByKey(ev, "address")
			involvedAccounts = append(involvedAccounts, address.String())
			switch ev.Type {
			case margintypes.EventForceClose:
				msgs = append(msgs, juno.NewMessage(e.TxHash, 0, ev.Type, ev.String(), involvedAccounts, height))
			case margintypes.EventOpen:
				msgs = append(msgs, juno.NewMessage(e.TxHash, 0, ev.Type, ev.String(), involvedAccounts, height))
			case margintypes.EventInterestRateComputation:
				msgs = append(msgs, juno.NewMessage(e.TxHash, 0, ev.Type, ev.String(), involvedAccounts, height))
			case margintypes.EventRepayFund:
				msgs = append(msgs, juno.NewMessage(e.TxHash, 0, ev.Type, ev.String(), involvedAccounts, height))
			case margintypes.EventBelowRemovalThreshold:
				msgs = append(msgs, juno.NewMessage(e.TxHash, 0, ev.Type, ev.String(), involvedAccounts, height))
			case margintypes.EventAboveRemovalThreshold:
				msgs = append(msgs, juno.NewMessage(e.TxHash, 0, ev.Type, ev.String(), involvedAccounts, height))
			case margintypes.EventIncrementalPayFund:
				msgs = append(msgs, juno.NewMessage(e.TxHash, 0, ev.Type, ev.String(), involvedAccounts, height))
			case margintypes.EventMarginUpdateParams:
				msgs = append(msgs, juno.NewMessage(e.TxHash, 0, ev.Type, ev.String(), involvedAccounts, height))
			}
		}

	}

	for _, i := range msgs {
		err := m.db.SaveMessage(i)
		if err != nil {
			fmt.Errorf("error while saving messages %s", err)
		}
	}

	return nil
}
