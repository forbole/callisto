package operations

import (
	"time"

	"github.com/forbole/bdjuno/database"
	"github.com/rs/zerolog/log"
)

func UpdateBlockTimeInMinute(blockTime time.Time, blockHeight int64, db database.BigDipperDb) error {
	log.Debug().
		Str("module", "staking").
		Str("operation", " tokens").
		Msg("getting total token supply")
	
	genesis,err := db.GetGenesisTime()
	if err!=nil{
		return err
	}

	//check if chain is not created minutes ago
	if(blockTime.Sub(genesis).Minutes()<0){
		return nil
	}

	minute, err := db.GetBlockHeightTimeMinuteAgo(blockTime)
	if err != nil {
		return err
	}
	newBlockTime := blockTime.Sub(minute.Timestamp).Seconds()/float64((blockHeight-minute.Height))

	return db.SaveAverageBlockTimePerMin(newBlockTime,blockTime,blockHeight)
}
