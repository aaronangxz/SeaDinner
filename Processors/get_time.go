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
	return time.Date(year, month, day, Config.OrderTime.Hour, Config.OrderTime.Minutes, Config.OrderTime.Seconds, 0, now.Location())
}
