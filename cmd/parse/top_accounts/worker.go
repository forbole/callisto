package top_accounts

import (
	topaccounts "github.com/forbole/bdjuno/v3/modules/top_accounts"
	"github.com/rs/zerolog/log"
)

type AddressQueue chan string

func NewQueue(size int) AddressQueue {
	return make(chan string, size)
}

type Worker struct {
	queue             AddressQueue
	topaccountsModule *topaccounts.Module
}

func NewWorker(queue AddressQueue, topaccountsModule *topaccounts.Module) Worker {
	return Worker{
		queue:             queue,
		topaccountsModule: topaccountsModule,
	}
}

func (w Worker) start() {
	for address := range w.queue {
		err := w.topaccountsModule.RefreshAll(address)
		if err != nil {
			log.Error().Str("account", address).Err(err).Msg("re-enqueueing failed address")

			go func(address string) {
				w.queue <- address
			}(address)
		}

	}
}
