package slashing

import (
	"time"
)

type StakingModule interface {
	RefreshValidatorInfos(height int64, timestamp time.Time, valOper string) error
}
