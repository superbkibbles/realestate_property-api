package date_utils

import "time"

const (
	apiDateLayout = "2006-1-02T15:04:05Z"
	apiDBLayout   = "2006-1-02 15:04:05"
)

func GetNow() time.Time {
	return time.Now().UTC()
}

func GetNowString() string {
	return GetNow().Format(apiDateLayout)
}

func GetNowDBFromat() string {
	return GetNow().Format(apiDBLayout)
}
