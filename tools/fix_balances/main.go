/*
 * Copyright (c) 2023. Business Process Technologies. All rights reserved.
 */

package main

import (
	"context"

	"git.ooo.ua/vipcoin/lib/filter"

	"fix_balances/repository"
)

// main - при парсинге нашего блокчейна bdjuno, возникли расхождения в балансах наших кошельков.
// Эти расхождения возникают из-за многопоточного парсинга, что влечет нарушение консистентности данных.
// Этот скрипт необходим для выравнивания наших балансов. Мы получаем кошельки из нашего блокчейна вместе с их балансами
// и обновляем балансы кошельков в базе данных bdjuno.
func main() {
	log, db, chainClient, err := initDependencies()
	if err != nil {
		log.Fatalf("Failed to initialize dependencies: %v", err)
	}

	repo := repository.NewRepository(db, chainClient)

	junoWallets, err := repo.GetWallets(context.Background(), filter.NewFilter())
	if err != nil {
		log.Fatalf("Failed to get wallets: %v", err)
	}

	chainWallets, err := repo.GetAllWalletFromChain(context.Background())
	if err != nil {
		log.Fatalf("Failed to get wallets: %v", err)
	}

	for i, _ := range junoWallets {
		chainWallet, found := chainWallets[junoWallets[i].Address]
		if !found {
			continue
		}
		junoWallets[i].Balance.Balance = chainWallet.Balance
	}

	if err = repo.UpdateWallets(context.Background(), junoWallets...); err != nil {
		log.Fatalf("Failed to update users: %v", err)
	}
}
