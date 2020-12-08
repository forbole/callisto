package consensus

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/forbole/bdjuno/database"

	"github.com/desmos-labs/juno/client"
	constypes "github.com/forbole/bdjuno/x/consensus/types"
	"github.com/rs/zerolog/log"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"
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

	var channels []<-chan tmctypes.ResultEvent
	for index, event := range events {
		channels = append(channels, subscribeConsensusEvent(index, event, cp))
	}

	merged := merge(channels...)
	go func() {
		for event := range merged {

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
func subscribeConsensusEvent(index int, event string, cp *client.Proxy) <-chan tmctypes.ResultEvent {
	query := fmt.Sprintf("tm.event = '%s'", event)
	subscriber := fmt.Sprintf("juno-client-consensus-%d", index)

	eventCh, cancel, err := cp.SubscribeEvents(subscriber, query)
	defer cancel()

	if err != nil {
		log.Fatal().Err(err).Send()
	}

	return eventCh
}

// merge takes a list of read-only channels and merges them into a single read-only channel
func merge(cs ...<-chan tmctypes.ResultEvent) <-chan tmctypes.ResultEvent {
	out := make(chan tmctypes.ResultEvent)

	var wg sync.WaitGroup
	wg.Add(len(cs))

	for _, c := range cs {
		go func(c <-chan tmctypes.ResultEvent) {
			for v := range c {
				out <- v
			}
			wg.Done()
		}(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
