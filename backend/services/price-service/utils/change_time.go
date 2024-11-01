package utils

import (
	"fmt"
	"time"
)

func ConvertMillisecondsToHHMMSS(ms int64) string {
	duration := time.Duration(ms) * time.Millisecond

	// Format to hh:mm:ss
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	seconds := int(duration.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}

func ConvertMillisecondsToTimestamp(ms int64) string {
	t := time.Unix(ms/1000, (ms%1000)*1000000)
	return t.Format("2006-01-02 15:04:05")
}
