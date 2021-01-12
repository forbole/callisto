package consensus

import (
	"encoding/json"
	"fmt"

	"github.com/forbole/bdjuno/database"

	"github.com/desmos-labs/juno/client"
	"github.com/rs/zerolog/log"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"

	constypes "github.com/forbole/bdjuno/x/consensus/types"
)

// ListenOperation allows to start listening to new consensus events properly
func ListenOperation(cp *client.Proxy, db *database.BigDipperDb) error {
	events := []string{
		tmtypes.EventCompleteProposal,
		tmtypes.EventNewRound,
		tmtypes.EventNewRoundStep,
		tmtypes.EventPolka,
		tmtypes.EventValidBlock,
		tmtypes.EventVote,
	}

	go func() {
		var out = make(chan tmctypes.ResultEvent)
		for _, event := range events {
			err := subscribeConsensusEvent(event, cp, out)
			if err != nil {
				log.Fatal().Err(err)
			}
		}

		for event := range out {

			// Serialize the event data as JSON to later de-serialize it into our custom object
			bz, err := json.Marshal(event.Data)
			if err != nil {
				log.Fatal().Err(err).Send()
			}

			// De-serialize the data into our custom object
			var consEvent constypes.ConsensusEvent
			err = json.Unmarshal(bz, &consEvent)
			if err != nil {
				log.Fatal().Err(err).Send()
			}

			log.Debug().
				Str("module", "consensus").
				Int64("height", consEvent.Height).
				Int("round", consEvent.Round).
				Str("step", consEvent.Step).
				Msg("saving consensus")

			// Save the event
			err = db.SaveConsensus(consEvent)
			if err != nil {
				log.Fatal().Err(err).Send()
			}
		}
	}()

	return nil
}

// subscribeConsensusEvent allows to subscribe to the consensus event having the given name,
// and returns a read-only channel emitting all the events
func subscribeConsensusEvent(event string, cp *client.Proxy, out chan<- tmctypes.ResultEvent) error {
	query := fmt.Sprintf("tm.event = '%s'", event)
	subscriber := fmt.Sprintf("juno-client-consensus-%s", event)

	eventCh, cancel, err := cp.SubscribeEvents(subscriber, query)
	if err != nil {
		return err
	}
	defer cancel()

	go func() {
		for event := range eventCh {
			out <- event
		}
	}()

	return nil
}
