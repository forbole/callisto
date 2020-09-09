package operations

import (
	"time"

	"github.com/forbole/bdjuno/database"
	"github.com/rs/zerolog/log"
)

func UpdateBlockTime(blockTime time.Time, blockHeight int64, db database.BigDipperDb) error {
	log.Debug().
		Str("module", "staking").
		Str("operation", " tokens").
		Msg("getting total token supply")
	
	genesis,err := db.GetGenesisTime()
	if err!=nil{
		return err
	}

	if(blockTime.Sub(genesis).Minutes()>0){
		minute, err := db.GetBlockHeightTimeMinuteAgo(blockTime)
		if err != nil {
			return err
		}
		minutesub := blockTime.Sub(minute.Timestamp).Seconds()
		
	}

	/* minute, err := db.GetBlockHeightTimeMinuteAgo(blockTime)
	if err!=nil{
		return err
	}
	hour, err := db.GetBlockHeightTimeHourAgo(blockTime)
	if err!=nil{
		return err
	} */


	print(minutesub)

	return nil
}
