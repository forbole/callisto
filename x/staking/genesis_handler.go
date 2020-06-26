package staking


import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/desmos-labs/juno/parse/worker"
	"github.com/forbole/bdjuno/database"
	"github.com/rs/zerolog/log"
	"github.com/tendermint/tendermint/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/staking/exported"


)

func GenesisHandler(codec *codec.Codec, genDoc *types.GenesisDoc, appState map[string]json.RawMessage, w worker.Worker) error {
	log.Debug().Str("module", "auth").Msg("parsing genesis")

	db, ok := w.Db.(database.BigDipperDb)
	if !ok {
		return fmt.Errorf("given database instance is not a BigDipperDb")
	}

	var stakingGenesisState staking.GenesisState
	if err := codec.UnmarshalJSON(appState[staking.ModuleName], &stakingGenesisState); err != nil {
		return err
	}

	if err := codec.UnmarshalJSON(appState[staking.ModuleName], &stakingGenesisState); err != nil {
		return err
	}

	// Store the accounts
	accounts := make([]exported.DelegationI, len(stakingGenesisState.Validators))
	for index, account := range stakingGenesisState.Validators {
		accounts[index] = account.(exported.Account)
		selfAddress := sdk.AccAddress(account[index].Bytes())
	}




	return nil
}