package types

import (
	"fmt"
	"net/http"
	"time"

	actionlogging "github.com/forbole/bdjuno/v2/cmd/actions/logging"
)

type GraphQLError struct {
	Message string `json:"message"`
}

func promActionCounter(path string) {
	actionlogging.ActionCounter.
		WithLabelValues(path, fmt.Sprintf("%d", http.StatusOK)).Inc()
}

func promErrorCounter(path string) {
	actionlogging.ActionErrorCounter.
		WithLabelValues(path, fmt.Sprintf("%d", http.StatusInternalServerError)).Inc()
}

func promReponseTimeLog(path string, start time.Time) {
	actionlogging.ActionResponseTime.
		WithLabelValues(path, fmt.Sprintf("%v", time.Since(start).Seconds())).
		Observe(time.Since(start).Seconds())
}
