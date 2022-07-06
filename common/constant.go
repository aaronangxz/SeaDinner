package common

import "time"

//goland:noinspection ALL
const (
	ONE_HOUR  = int64(60*time.Minute) / int64(1*time.Second)
	ONE_DAY   = 24 * ONE_HOUR
	ONE_WEEK  = 7 * ONE_DAY
	ONE_MONTH = 4 * ONE_WEEK
)
