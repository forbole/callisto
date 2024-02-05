package gov

import (
	"fmt"
	"strconv"
	"time"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	juno "github.com/forbole/juno/v5/types"
	"github.com/rs/zerolog/log"
	abci "github.com/tendermint/tendermint/abci/types"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
)

// HandleBlock implements modules.BlockModule
func (m *Module) HandleBlock(
	b *tmctypes.ResultBlock, blockResults *tmctypes.ResultBlockResults, txs []*juno.Tx, vals *tmctypes.ResultValidators,
) error {
	txEvents := collectTxEvents(txs)
	err := m.updateProposalsStatus(b.Block.Height, b.Block.Time, txEvents, blockResults.EndBlockEvents, vals)
	if err != nil {
		log.Error().Str("module", "gov").Int64("height", b.Block.Height).
			Err(err).Msg("error while updating proposals")
	}
	return nil
}

// updateProposalsStatus updates the status of proposals if they have been included in the EndBlockEvents or status
// was changed from deposit to voting
func (m *Module) updateProposalsStatus(height int64, blockTime time.Time, txEvents, endBlockEvents []abci.Event, blockVals *tmctypes.ResultValidators) error {
	var ids []uint64
	// check if EndBlockEvents contains active_proposal event
	endBlockIDs, err := findProposalIDsInEvents(endBlockEvents, govtypes.EventTypeActiveProposal, govtypes.AttributeKeyProposalID)
	if err != nil {
		return err
	}
	ids = append(ids, endBlockIDs...)

	// the proposal changes state from the deposit to voting
	txIDs, err := findProposalIDsInEvents(txEvents, govtypes.EventTypeProposalDeposit, govtypes.AttributeKeyVotingPeriodStart)
	if err != nil {
		return err
	}
	ids = append(ids, txIDs...)

	// update status for proposals IDs stored in ids array
	for _, id := range ids {
		err = m.UpdateProposal(height, blockTime, id)
		if err != nil {
			return fmt.Errorf("error while updating proposal: %s", err)
		}

		err = m.UpdateProposalValidatorStatusesSnapshot(height, blockVals, id)
		if err != nil {
			return fmt.Errorf("error while updating proposal validator statuses snapshots: %s", err)
		}

		err = m.UpdateProposalStakingPoolSnapshot(height, blockVals, id)
		if err != nil {
			return fmt.Errorf("error while updating proposal validator statuses snapshots: %s", err)
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
			if string(attr.Key) != attrKey {
				continue
			}
			// parse proposal ID from []byte to unit64
			id, err := strconv.ParseUint(string(attr.Value), 10, 64)
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
