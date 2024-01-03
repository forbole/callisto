package testutils

import (
	"time"
)

func NewDurationPointer(duration time.Duration) *time.Duration {
	return &duration
}

func NewTimePointer(time time.Time) *time.Time {
	return &time
}
