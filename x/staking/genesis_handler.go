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
	dbtypes "github.com/forbole/bdjuno/database/types"

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
	err = InitialCommission(stakingGenesisState,db)
	if err!=nil{
		return err
	}
	return nil
}

func InitialCommission(stakingGenesisState staking.GenesisState,db database.BigDipperDb)error{
	// Store the accounts
	accounts := make([]dbtypes.ValidatorCommission, len(stakingGenesisState.Validators))
	for index, account := range stakingGenesisState.Validators {
		accounts[index] = dbtypes.ValidatorCommission{
			ValidatorAddress  : account.ConsAddress.String(),
	        Timestamp         : time.Time.now,
	        Commission        : account.Commission,
	        MinSelfDelegation : account.MinSelfDelegation,
	        Height            : 0,
		}
	}

	err := db.SaveValidatorCommissions(accounts)
	if err!=nil{
		return err
	}
	return nil
}

func InitialInformation(stakingGenesisState staking.GenesisState,db database.BigDipperDb)error{
	accounts := make([]dbtypes.ValidatorInfoRow, len(stakingGenesisState.Validators))
	for index, account := range stakingGenesisState.Validators {
		accounts[index] = dbtypes.ValidatorInfoRow{
			consensus_address    : account.ConsAddress.String(),
			operator_address     : account.OperatorAddress.String(),
			moniker             :  account.Description.Moniker,
			identity            : account.Description.Identity,
			website              :account.Description.Website,
			securityContact    :  account.Description.SecurityContact,
			details            :  account.Description.Details,
		}
	}

	err := db.SaveInitialValidatorInfo(accounts)
	if err!=nil{
		return err
	}
	return nil
}
}