package processors

import (
	"context"
	"fmt"
	"github.com/aaronangxz/SeaDinner/common"
	"os"
	"strconv"
	"time"

	"github.com/aaronangxz/SeaDinner/log"
	"github.com/aaronangxz/SeaDinner/sea_dinner.pb"
)

const (
	//TimeZone constant timezone
	TimeZone = "Asia/Singapore"
)

var (
	tz, _ = time.LoadLocation(TimeZone)
)

//UnixToUTC Converts current unix time to UTC time object
func UnixToUTC(timestamp int64) time.Time {
	return time.Unix(timestamp, 0).Local().UTC()
}

//GetLunchTime Returns the lunchtime of today, defined in Config, as time object
func GetLunchTime() time.Time {
	now := time.Now().In(tz)
	year, month, day := now.Date()
	return time.Date(year, month, day, common.Config.OrderTime.Hour, common.Config.OrderTime.Minutes, common.Config.OrderTime.Seconds, 0, now.Location())
}

//GetPreviousDayLunchTime Returns the lunchtime of yesterday, defined in Config, as time object
func GetPreviousDayLunchTime() time.Time {
	now := time.Now().In(tz)
	year, month, day := now.Add(time.Duration(-1*24) * time.Hour).Date()
	return time.Date(year, month, day, common.Config.OrderTime.Hour, common.Config.OrderTime.Minutes, common.Config.OrderTime.Seconds, 0, now.Location())
}

//ConvertTimeStamp Converts current unix timestamp to yyyy-mm-dd format
//time format: Mon Jan 2 15:04:05 -0700 MST 2006
func ConvertTimeStamp(timestamp int64) string {
	return fmt.Sprint(UnixToUTC(timestamp).In(tz).Format("2006-01-02"))
}

//ConvertTimeStampMonthDay Converts current unix timestamp to d/y format
func ConvertTimeStampMonthDay(timestamp int64) string {
	return fmt.Sprint(UnixToUTC(timestamp).In(tz).Format("2/1"))
}

//ConvertTimeStampDayOfWeek Converts current unix timestamp to DDD dd/mm format
func ConvertTimeStampDayOfWeek(timestamp int64) string {
	return fmt.Sprint(UnixToUTC(timestamp).In(tz).Format("Mon 02/01"))
}

//ConvertTimeStampWeekOfYear Converts current unix timestamp to the week number of year
func ConvertTimeStampWeekOfYear(timestamp int64) (int64, int64) {
	year, week := UnixToUTC(timestamp).In(tz).ISOWeek()
	return int64(year), int64(week)
}

//ConvertTimeStampTime Converts current unix timestamp to m:ss format
func ConvertTimeStampTime(timestamp int64) string {
	return fmt.Sprint(UnixToUTC(timestamp).In(tz).Format("3:04PM"))
}

//ShouldOrder Checks if today is weekday + has dinner
func ShouldOrder() bool {
	return IsWeekDay() && IsActiveDay()
}

//IsWeekDay Checks if today is a weekday
func IsWeekDay() bool {
	day := time.Now().In(tz).Weekday()
	return day >= 1 && day <= 5
}

//IsActiveDay Checks if today has dinner (If today is a holiday)
func IsActiveDay() bool {
	return GetDayID(context.TODO()) != 0
}

//IsNotEOW Checks if today is not friday, saturday, sunday
func IsNotEOW(t time.Time) bool {
	tz, _ := time.LoadLocation(TimeZone)
	day := t.In(tz).Weekday()
	return day >= 1 && day < 5
}

//IsSOW Checks if today is a monday
func IsSOW(t time.Time) bool {
	tz, _ := time.LoadLocation(TimeZone)
	day := t.In(tz).Weekday()
	return day == 1
}

//IsPollStart Checks if order polling has started
func IsPollStart() bool {
	var (
		status *sea_dinner.Current
		key    = os.Getenv("TOKEN")
	)
	txn := App.StartTransaction("is_poll_start")
	defer txn.End()

	if common.Config.Adhoc {
		return true
	}
	_, err := Client.R().
		SetHeader("Authorization", MakeToken(context.TODO(), key)).
		SetResult(&status).
		Get(MakeURL(int(sea_dinner.URLType_URL_CURRENT), nil))

	if err != nil {
		log.Error(context.TODO(), err.Error())
	}

	return status.GetMenu().GetActive()
}

//WeekStartEndDate Returns the start and end day of the current week in SGT unix time
func WeekStartEndDate(timestamp int64) (int64, int64) {
	date := UnixToUTC(timestamp).In(tz)

	startOffset := (int(time.Monday) - int(date.Weekday()) - 7) % 7
	startResult := date.Add(time.Duration(startOffset*24) * time.Hour)
	endResult := startResult.Add(time.Duration(4*24) * time.Hour)

	startYear, startMonth, startDay := startResult.Date()
	endYear, endMonth, endDay := endResult.Date()
	return time.Date(startYear, startMonth, startDay, 0, 0, 0, 0, tz).Unix(), time.Date(endYear, endMonth, endDay, 23, 59, 59, 59, tz).Unix()
}

//MonthStartEndDate Returns the start and end day of the current month in SGT unix time
func MonthStartEndDate(timestamp int64) (int64, int64) {
	date := UnixToUTC(timestamp).In(tz)
	currentYear, currentMonth, _ := date.Date()
	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, tz)
	lastOfMonth := time.Date(currentYear, currentMonth+1, 0, 23, 59, 59, 59, tz)
	return firstOfMonth.Unix(), lastOfMonth.Unix()
}

//YearStartEndDate Returns the start and end day of the current year in SGT unix time
func YearStartEndDate(timestamp int64) (int64, int64) {
	date := UnixToUTC(timestamp).In(tz)
	currentYear, _, _ := date.Date()
	return time.Date(currentYear, time.January, 1, 0, 0, 0, 0, tz).Unix(), time.Date(currentYear, time.December, 31, 23, 59, 59, 59, tz).Unix()
}

//IsSendReminderTime Checks if it is 2 hours prior to the pre-defined lunch time
func IsSendReminderTime() bool {
	offset, _ := strconv.ParseInt(os.Getenv("SEND_REMINDER_OFFSET_SECONDS"), 10, 64)
	return ShouldOrder() && time.Now().Unix() == GetLunchTime().Add(time.Duration(-offset)*time.Second).Unix()
}

//IsPrepOrderTime Checks if it is within 1 minute before and 15 seconds before the pre-defined lunch time
func IsPrepOrderTime() bool {
	return ShouldOrder() &&
		time.Now().Unix() >= GetLunchTime().Add(time.Duration(-60)*time.Second).Unix() &&
		time.Now().Unix() <= GetLunchTime().Add(time.Duration(-15)*time.Second).Unix()
}

//IsOrderTime Checks if it is lunch time
func IsOrderTime() bool {
	return ShouldOrder() && time.Now().Unix() == GetLunchTime().Unix()
}

//GetOffWorkTime Returns the off-work time
func GetOffWorkTime() time.Time {
	now := time.Now().In(tz)
	year, month, day := now.Date()
	return time.Date(year, month, day, 19, 0, 0, 0, now.Location())
}

//GetEOD Returns 23:59:59 of today
func GetEOD() time.Time {
	now := time.Now().In(tz)
	year, month, day := now.Date()
	return time.Date(year, month, day, 23, 59, 59, 0, now.Location())
}

//IsSendCheckInTime Checks if it is eod
func IsSendCheckInTime() bool {
	return ShouldOrder() && time.Now().Unix() == GetOffWorkTime().Unix()
}

//IsDeleteCheckInTime Checks if it is check-in expiry time
func IsDeleteCheckInTime() bool {
	return ShouldOrder() && time.Now().Unix() == GetOffWorkTime().Add(time.Duration(5400)*time.Second).Unix()
}
