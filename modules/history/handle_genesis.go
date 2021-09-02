package history

import (
	"time"

	"github.com/forbole/bdjuno/database"
	historyutils "github.com/forbole/bdjuno/modules/history/utils"
)

// HandleGenesis handles the genesis operations for the history module, storing the initial history data
func HandleGenesis(genesisTime time.Time, db *database.Db) error {
	accounts, err := db.GetAccounts()
	if err != nil {
		return err
	}

	for _, account := range accounts {
		err = historyutils.UpdateAccountBalanceHistoryWithTime(account, genesisTime, db)
		if err != nil {
			return err
		}
	}

	return nil
}
