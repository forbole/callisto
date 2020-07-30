package operations

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/desmos-labs/juno/parse/client"
	"github.com/forbole/bdjuno/database"
	"github.com/forbole/bdjuno/x/staking/types"

	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/rs/zerolog/log"
	tmctypes "github.com/tendermint/tendermint/rpc/core/types"
)

// UpdateValidatorVotingPower fetches and stores into the database all the current validators' voting powers
func UpdateValidatorVotingPower(cp client.ClientProxy, db database.BigDipperDb) error {
	log.Debug().
		Str("module", "staking").
		Str("operation", " voting percentage").
		Msg("getting validators  voting percentage")

	// First, get the latest block height
	var block tmctypes.ResultBlock
	if err := cp.QueryLCD("/blocks/latest", &block); err != nil {
		return err
	}
	// Second, get the validators
	var validators rpc.ResultValidatorsOutput
	endpoint := fmt.Sprintf("/validatorsets/%d", block.Block.Height)
	_, err := cp.QueryLCDWithHeight(endpoint, &validators)
	if err != nil {
		return err
	}
	// Store the signing infos into the database
	log.Debug().
		Str("module", "staking").
		Str("operation", "uptime").
		Msg("saving  voting percentage")
	var votings []types.ValidatorVotingPower
	for _, validator := range validators.Validators {
		if found, _ := db.HasValidator(validator.Address.String()); !found {
			continue
		}
		consAddress, err := sdk.ConsAddressFromBech32(validator.Address.String())
		if err != nil {
			return err
		}
		votings = append(votings, types.NewValidatorVotingPower(
			consAddress,
			validator.VotingPower,
			block.Block.Height,
		))
	}

	if err := db.SaveVotingPowers(votings); err != nil {
		return err
	}
	return nil
}
