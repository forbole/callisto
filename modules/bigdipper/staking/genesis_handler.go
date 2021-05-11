package staking

import (
	"encoding/json"
	"fmt"

	"github.com/forbole/bdjuno/modules/common/staking"

	tmtypes "github.com/tendermint/tendermint/types"

	bigdipperdb "github.com/forbole/bdjuno/database/bigdipper"
	bstakingcommon "github.com/forbole/bdjuno/modules/bigdipper/staking/common"
	bstakingtypes "github.com/forbole/bdjuno/modules/bigdipper/staking/types"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/rs/zerolog/log"
)

func HandleGenesis(
	doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage, cdc codec.Marshaler, db *bigdipperdb.Db,
) error {
	log.Debug().Str("module", "staking").Msg("parsing genesis")

	err := staking.HandleGenesis(doc, appState, cdc, db)
	if err != nil {
		return err
	}

	err = parseStakingState(doc, appState, cdc, db)
	if err != nil {
		return err
	}

	return nil
}

// parseStakingState parses the staking genesis state and stores the data properly
func parseStakingState(
	doc *tmtypes.GenesisDoc, appState map[string]json.RawMessage, cdc codec.Marshaler, db *bigdipperdb.Db,
) error {
	// Read the genesis state
	var genState stakingtypes.GenesisState
	err := cdc.UnmarshalJSON(appState[stakingtypes.ModuleName], &genState)
	if err != nil {
		return err
	}

	// Save the description
	err = saveValidatorDescription(doc, genState.Validators, db)
	if err != nil {
		return fmt.Errorf("error while storing staking genesis validator descriptions: %s", err)
	}

	err = saveValidatorsCommissions(doc.InitialHeight, genState.Validators, db)
	if err != nil {
		return fmt.Errorf("error while storing staking genesis validators commissions: %s", err)
	}

	return nil
}

// saveValidatorsCommissions save the initial commission for each validator
func saveValidatorsCommissions(height int64, validators stakingtypes.Validators, db *bigdipperdb.Db) error {
	for _, account := range validators {
		err := db.SaveValidatorCommission(bstakingtypes.NewValidatorCommission(
			account.OperatorAddress,
			&account.Commission.Rate,
			&account.MinSelfDelegation,
			height,
		))
		if err != nil {
			return err
		}
	}

	return nil
}

// saveValidatorDescription saves the description for the given validators
func saveValidatorDescription(doc *tmtypes.GenesisDoc, validators stakingtypes.Validators, db *bigdipperdb.Db) error {
	for _, account := range validators {
		description, err := bstakingcommon.GetValidatorDescription(
			account.OperatorAddress,
			account.Description,
			doc.InitialHeight,
		)
		if err != nil {
			return err
		}

		err = db.SaveValidatorDescription(description)
		if err != nil {
			return err
		}
	}

	return nil
}
