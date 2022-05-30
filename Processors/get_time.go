package Processors

import (
	"fmt"
	"os"
	"time"
)

const (
	TimeZone = "Asia/Singapore"
)

var (
	tz, _ = time.LoadLocation(TimeZone)
)

func UnixToUTC(timestamp int64) time.Time {
	return time.Unix(timestamp, 0).Local().UTC()
}

func GetLunchTime() time.Time {
	now := time.Now().In(tz)
	year, month, day := now.Date()
	return time.Date(year, month, day, Config.OrderTime.Hour, Config.OrderTime.Minutes, Config.OrderTime.Seconds, 0, now.Location())
}

func GetOffWorkTime() time.Time {
	now := time.Now().In(tz)
	year, month, day := now.Date()
	return time.Date(year, month, day, 19, 15, 0, 0, now.Location())
}

//time format: Mon Jan 2 15:04:05 -0700 MST 2006
func ConvertTimeStamp(timestamp int64) string {
	return fmt.Sprint(UnixToUTC(timestamp).In(tz).Format("2006-01-02"))
}

func ConvertTimeStampMonthDay(timestamp int64) string {
	return fmt.Sprint(UnixToUTC(timestamp).In(tz).Format("2/1"))
}

func ConvertTimeStampDayOfWeek(timestamp int64) string {
	return fmt.Sprint(UnixToUTC(timestamp).In(tz).Format("Mon 2/1"))
}

func ConvertTimeStampTime(timestamp int64) string {
	return fmt.Sprint(UnixToUTC(timestamp).In(tz).Format("3:04PM"))
}

func ShouldOrder() bool {
	return IsWeekDay() && IsActiveDay()
}

func IsWeekDay() bool {
	day := time.Now().In(tz).Weekday()
	return day >= 1 && day <= 5
}

func IsActiveDay() bool {
	return GetDayId() != 0
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
	date := UnixToUTC(timestamp).In(tz)

	startOffset := (int(time.Monday) - int(date.Weekday()) - 7) % 7
	startResult := date.Add(time.Duration(startOffset*24) * time.Hour)
	endResult := startResult.Add(time.Duration(4*24) * time.Hour)

	startYear, startMonth, startDay := startResult.Date()
	endYear, endMonth, endDay := endResult.Date()
	return time.Date(startYear, startMonth, startDay, 0, 0, 0, 0, tz).Unix(), time.Date(endYear, endMonth, endDay, 23, 59, 59, 59, tz).Unix()
}

func IsSendReminderTime() bool {
	return ShouldOrder() && time.Now().Unix() == GetLunchTime().Add(time.Duration(-2)*time.Hour).Unix()
}

func IsPrepOrderTime() bool {
	return ShouldOrder() &&
		time.Now().Unix() >= GetLunchTime().Add(time.Duration(-60)*time.Second).Unix() &&
		time.Now().Unix() <= GetLunchTime().Add(time.Duration(-15)*time.Second).Unix()
}

func IsOrderTime() bool {
	return ShouldOrder() && time.Now().Unix() == GetLunchTime().Unix()
}

func IsSendCheckInTime() bool {
	return ShouldOrder() && time.Now().Unix() == GetOffWorkTime().Unix()
}
