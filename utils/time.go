package utils

import "time"

func ISO8601ToTimestamp(timeStr string) (int, error) {
	layout := "2006-01-02T15:04:05Z"
	t, err := time.Parse(layout, timeStr)
	if err != nil {
		return 0, err
	}
	return int(t.Unix()), nil
}
