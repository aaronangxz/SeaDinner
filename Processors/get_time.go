package Processors

import (
	"fmt"
	"time"
)

const (
	timeZone = "Asia/Singapore"
)

func GetLunchTime() time.Time {
	tz, err := time.LoadLocation(timeZone)
	if err != nil {
		fmt.Println(err)
	}
	now := time.Now().In(tz)

	year, month, day := now.Date()
	return time.Date(year, month, day, 12, 30, 0, 0, now.Location())
}
