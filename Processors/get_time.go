package Processors

import (
	"fmt"
	"os"
	"time"
)

const (
	TimeZone = "Asia/Singapore"
)

func GetLunchTime() time.Time {
	tz, err := time.LoadLocation(TimeZone)
	if err != nil {
		fmt.Println(err)
	}
	now := time.Now().In(tz)

	year, month, day := now.Date()
	return time.Date(year, month, day, Config.OrderTime.Hour, Config.OrderTime.Minutes, Config.OrderTime.Seconds, 0, now.Location())
}

//time format: Mon Jan 2 15:04:05 -0700 MST 2006
func ConvertTimeStamp(timestamp int64) string {
	t := time.Unix(timestamp, 0).Local().UTC()
	tz, _ := time.LoadLocation(TimeZone)
	return fmt.Sprint(t.In(tz).Format("2006-01-02"))
}

func ConvertTimeStampMonthDay(timestamp int64) string {
	t := time.Unix(timestamp, 0).Local().UTC()
	tz, _ := time.LoadLocation(TimeZone)
	return fmt.Sprint(t.In(tz).Format("2/1"))
}

func ConvertTimeStampDayOfWeek(timestamp int64) string {
	t := time.Unix(timestamp, 0).Local().UTC()
	tz, _ := time.LoadLocation(TimeZone)
	return fmt.Sprint(t.In(tz).Format("Mon 2/1"))
}

func ConvertTimeStampTime(timestamp int64) string {
	t := time.Unix(timestamp, 0).Local().UTC()
	tz, _ := time.LoadLocation(TimeZone)
	return fmt.Sprint(t.In(tz).Format("3:04PM"))
}

func IsWeekDay(t time.Time) bool {
	tz, _ := time.LoadLocation(TimeZone)
	day := t.In(tz).Weekday()
	return day >= 1 && day <= 5
}

func IsNotEOW(t time.Time) bool {
	tz, _ := time.LoadLocation(TimeZone)
	day := t.In(tz).Weekday()
	return day >= 1 && day < 5
}

func IsPollStart() bool {
	var (
		status Current
		key    = os.Getenv("TOKEN")
	)

	_, err := Client.R().
		SetHeader("Authorization", MakeToken(key)).
		SetResult(&status).
		Get(MakeURL(URL_CURRENT, nil))

	if err != nil {
		fmt.Println(err)
	}

	return status.Menu.GetActive()
}

func WeekStartEndDate(timestamp int64) (int64, int64) {
	tz, _ := time.LoadLocation(TimeZone)
	date := time.Unix(timestamp, 0).Local().UTC().In(tz)

	startOffset := (int(time.Monday) - int(date.Weekday()) - 7) % 7
	startResult := date.Add(time.Duration(startOffset*24) * time.Hour)
	endResult := startResult.Add(time.Duration(4*24) * time.Hour)

	startYear, startMonth, startDay := startResult.Date()
	endYear, endMonth, endDay := endResult.Date()
	return time.Date(startYear, startMonth, startDay, 0, 0, 0, 0, tz).Unix(), time.Date(endYear, endMonth, endDay, 23, 59, 59, 59, tz).Unix()
}
