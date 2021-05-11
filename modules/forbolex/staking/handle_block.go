package staking

import (
	"github.com/cosmos/cosmos-sdk/codec"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	tmtypes "github.com/tendermint/tendermint/types"

	forbolexdb "github.com/forbole/bdjuno/database/forbolex"
	"github.com/forbole/bdjuno/modules/common/staking"
)

// HandleBlock handles the given block by storing specific information inside the database
func HandleBlock(block *tmtypes.Block, client stakingtypes.QueryClient, cdc codec.Marshaler, db *forbolexdb.Db) error {
	_, err := staking.UpdateValidators(block.Height, client, cdc, db)
	return err
}
