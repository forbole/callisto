package utils

import (
	"time"
)

func AreTimesEqual(first *time.Time, second *time.Time) bool {
	if first == nil && second == nil {
		return true
	}

	if first == nil || second == nil {
		return false
	}

	return first.Equal(*second)
}
