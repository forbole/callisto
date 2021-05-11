package consensus

import (
	"fmt"
	juno "github.com/desmos-labs/juno/types"

	bigdipperdb "github.com/forbole/bdjuno/database/bigdipper"
	"github.com/forbole/bdjuno/modules/bigdipper/consensus/types"

	"github.com/rs/zerolog/log"

	"github.com/desmos-labs/juno/client"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

// ListenOperation allows to start listening to new consensus events properly
func ListenOperation(cp *client.Proxy, db *bigdipperdb.Db) {
	events := []string{
		tmtypes.EventNewRound,
		tmtypes.EventNewRoundStep,
		tmtypes.EventCompleteProposal,
		tmtypes.EventVote,
		tmtypes.EventPolka,
		tmtypes.EventValidBlock,
	}

	// This channel will be used to gather all the events
	var eventChan = make(chan tmctypes.ResultEvent, 10)

	for _, event := range events {
		go subscribeConsensusEvent(event, cp, eventChan)
	}

	for event := range eventChan {
		handleEvent(event, db)
	}
}

// subscribeConsensusEvent allows to subscribe to the consensus event having the given name,
// and returns a read-only channel emitting all the events
func subscribeConsensusEvent(event string, cp *client.Proxy, eventChan chan<- tmctypes.ResultEvent) {
	query := fmt.Sprintf("tm.event = '%s'", event)

	subscriber := fmt.Sprintf("%s-event-%s", juno.Cfg.GetRPCConfig().ClientName, event)
	eventCh, cancel, err := cp.SubscribeEvents(subscriber, query)
	if err != nil {
		log.Error().Str("module", "consensus").Err(err).Msg("error while subscribing to event")
		return
	}
	defer cancel()

	for event := range eventCh {
		eventChan <- event
	}
}

// handleEvent handles the given event storing its data inside the database properly
func handleEvent(event tmctypes.ResultEvent, db *bigdipperdb.Db) {
	consEvent := mapEvent(event)
	if consEvent == nil {
		return
	}

	// Save the event
	log.Debug().Str("module", "consensus").
		Int64("height", consEvent.Height).
		Int32("round", consEvent.Round).
		Str("step", consEvent.Step).
		Msg("saving consensus event")

	err := db.SaveConsensus(consEvent)
	if err != nil {
		log.Error().Str("module", "consensus").Err(err).Msg("error while saving consensus event")
	}
}

// mapEvent converts the given ResultEvent to a ConsensusEvent instance
func mapEvent(event tmctypes.ResultEvent) *types.ConsensusEvent {
	switch data := event.Data.(type) {
	case tmtypes.EventDataNewRound:
		return types.NewConsensusEvent(data.Height, data.Round, data.Step)

	case tmtypes.EventDataRoundState:
		return types.NewConsensusEvent(data.Height, data.Round, data.Step)

	case tmtypes.EventDataCompleteProposal:
		return types.NewConsensusEvent(data.Height, data.Round, data.Step)

	case tmtypes.EventDataVote:
		return types.NewConsensusEvent(data.Vote.Height, data.Vote.Round, tmtypes.EventVote)

	default:
		return nil
	}
}
