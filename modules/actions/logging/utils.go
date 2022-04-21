package logging

import (
	"fmt"
	"net/http"
	"time"
)

func SuccessCounter(path string) {
	ActionCounter.WithLabelValues(path, fmt.Sprintf("%d", http.StatusOK)).Inc()
}

func ErrorCounter(path string) {
	ActionErrorCounter.WithLabelValues(path, fmt.Sprintf("%d", http.StatusInternalServerError)).Inc()
}

func ReponseTimeBuckets(path string, start time.Time) {
	ActionResponseTime.WithLabelValues(path).
		Observe(time.Since(start).Seconds())
}
