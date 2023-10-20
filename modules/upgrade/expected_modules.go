package upgrade

import (
	"time"
)

type StakingModule interface {
	RefreshAllValidatorInfos(height int64, timestamp time.Time) error
}
