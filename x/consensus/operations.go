package consensus

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/desmos-labs/juno/config"
	"github.com/desmos-labs/juno/db"
	"github.com/desmos-labs/juno/parse/client"
	"github.com/rs/zerolog/log"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

func ListeningOperation(_ config.Config, _ *codec.Codec, cp client.ClientProxy, _ db.Database) error {
	events := []string{
		//tmtypes.EventCompleteProposal,
		//tmtypes.EventLock,
		tmtypes.EventNewRound,
		tmtypes.EventNewRoundStep,
		//tmtypes.EventPolka,
		//tmtypes.EventRelock,
		//tmtypes.EventTimeoutPropose,
		//tmtypes.EventTimeoutWait,
		//tmtypes.EventUnlock,
		//tmtypes.EventValidBlock,
		//tmtypes.EventVote,
	}

	var channels []<-chan tmctypes.ResultEvent
	for index, event := range events {
		channels = append(channels, SubscribeConsensusEvent(index, event, cp))
	}

	merged := merge(channels...)
	go func() {
		for event := range merged {
			bz, err := json.Marshal(event)
			if err != nil {
				log.Fatal().Err(err).Send()
			}

			log.Debug().Msg(string(bz))
		}
	}()

	return nil
}

func SubscribeConsensusEvent(index int, event string, cp client.ClientProxy) <-chan tmctypes.ResultEvent {
	query := fmt.Sprintf("tm.event = '%s'", event)
	subscriber := fmt.Sprintf("juno-client-consensus-%d", index)

	eventCh, cancel, err := cp.SubscribeEvents(subscriber, query)
	defer cancel()

	if err != nil {
		log.Fatal().Err(err).Send()
	}

	return eventCh
}

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
