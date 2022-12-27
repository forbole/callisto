package top_accounts

import (
	"fmt"

	"github.com/forbole/bdjuno/v3/modules/bank"
	"github.com/forbole/bdjuno/v3/modules/distribution"
	"github.com/forbole/bdjuno/v3/modules/staking"
	topaccounts "github.com/forbole/bdjuno/v3/modules/top_accounts"
	"github.com/rs/zerolog/log"
)

type AddressQueue chan string

func NewQueue(size int) AddressQueue {
	return make(chan string, size)
}

type Worker struct {
	queue             AddressQueue
	bankModule        *bank.Module
	distriModule      *distribution.Module
	stakingModule     *staking.Module
	topaccountsModule *topaccounts.Module
}

func NewWorker(
	queue AddressQueue,
	bankModule *bank.Module, distriModule *distribution.Module,
	stakingModule *staking.Module, topaccountsModule *topaccounts.Module,
) Worker {
	return Worker{
		queue:             queue,
		bankModule:        bankModule,
		distriModule:      distriModule,
		stakingModule:     stakingModule,
		topaccountsModule: topaccountsModule,
	}
}

func (w Worker) start() {
	for address := range w.queue {
		err := w.refreshAll(address)
		if err != nil {
			log.Error().Str("account", address).Msg("re-enqueueing failed address")

			go func(address string) {
				w.queue <- address
			}(address)
		}
	}
}

func (w *Worker) refreshAll(address string) error {
	err := w.bankModule.UpdateBalances([]string{address}, 0)
	if err != nil {
		return fmt.Errorf("error while refreshing account balance of account %s", address)
	}

	err = w.stakingModule.RefreshDelegations(0, address)
	if err != nil {
		return fmt.Errorf("error while refreshing delegations of account %s", address)
	}

	err = w.stakingModule.RefreshRedelegations(0, address)
	if err != nil {
		return fmt.Errorf("error while refreshing redelegations of account %s", address)
	}

	err = w.stakingModule.RefreshUnbondings(0, address)
	if err != nil {
		return fmt.Errorf("error while refreshing unbonding delegations of account %s", address)
	}

	err = w.distriModule.RefreshDelegatorRewards(0, []string{address})
	if err != nil {
		return fmt.Errorf("error while refreshing rewards of account %s", address)
	}

	err = w.topaccountsModule.RefreshTopAccountsSum([]string{address})
	if err != nil {
		return fmt.Errorf("error while refreshing top account sum of account %s", address)
	}

	return nil
}
