package gov

import (
	"fmt"
	"strconv"

	juno "github.com/forbole/juno/v5/types"

	tmctypes "github.com/cometbft/cometbft/rpc/core/types"

	abci "github.com/cometbft/cometbft/abci/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/rs/zerolog/log"
)

// HandleBlock implements modules.BlockModule
func (m *Module) HandleBlock(
	b *tmctypes.ResultBlock, blockResults *tmctypes.ResultBlockResults, txs []*juno.Tx, _ *tmctypes.ResultValidators,
) error {
	txEvents := collectTxEvents(txs)
	err := m.updateProposalsStatus(b.Block.Height, txEvents, blockResults.EndBlockEvents)
	if err != nil {
		log.Error().Str("module", "gov").Int64("height", b.Block.Height).
			Err(err).Msg("error while updating proposals")
	}

	return nil
}

// updateProposalsStatus updates the status of proposals if they have been included in the EndBlockEvents or status
// was changed from deposit to voting
func (m *Module) updateProposalsStatus(height int64, txEvents, endBlockEvents []abci.Event) error {
	var ids []uint64
	// check if EndBlockEvents contains active_proposal event
	endBlockIDs, err := findProposalIDsInEvents(endBlockEvents, govtypes.EventTypeActiveProposal, govtypes.AttributeKeyProposalID)
	if err != nil {
		return err
	}
	ids = append(ids, endBlockIDs...)

	// the proposal changes state from the submit to voting
	idsInSubmitTxs, err := findProposalIDsInEvents(txEvents, govtypes.EventTypeSubmitProposal, govtypes.AttributeKeyVotingPeriodStart)
	if err != nil {
		return err
	}
	ids = append(ids, idsInSubmitTxs...)

	// the proposal changes state from the deposit to voting
	idsInDepositTxs, err := findProposalIDsInEvents(txEvents, govtypes.EventTypeProposalDeposit, govtypes.AttributeKeyVotingPeriodStart)
	if err != nil {
		return err
	}
	ids = append(ids, idsInDepositTxs...)

	// update status for proposals IDs stored in ids array
	for _, id := range ids {
		err := m.UpdateProposalStatus(height, id)
		if err != nil {
			return fmt.Errorf("error while updating proposal %d status: %s", id, err)
		}
	}

	return nil
}

func findProposalIDsInEvents(events []abci.Event, eventType, attrKey string) ([]uint64, error) {
	ids := make([]uint64, 0)
	for _, event := range events {
		if event.Type != eventType {
			continue
		}
		for _, attr := range event.Attributes {
			if attr.Key != attrKey {
				continue
			}
			// parse proposal ID from []byte to unit64
			id, err := strconv.ParseUint(attr.Value, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("error while parsing proposal id: %s", err)
			}
			// add proposal ID to ids array
			ids = append(ids, id)
		}
	}

	return ids, nil
}

func collectTxEvents(txs []*juno.Tx) []abci.Event {
	events := make([]abci.Event, 0)
	for _, tx := range txs {
		for _, ev := range tx.Events {
			events = append(events, abci.Event{Type: ev.Type, Attributes: ev.Attributes})
		}
	}

	return events
}
