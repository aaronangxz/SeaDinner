package processors

import (
	"github.com/aaronangxz/SeaDinner/common"
	"reflect"
	"testing"
	"time"
)

func TestGetLunchTime(t *testing.T) {
	tz, _ := time.LoadLocation(TimeZone)
	now := time.Now().In(tz)
	year, month, day := now.Date()
	expectedTime := time.Date(year, month, day, common.Config.OrderTime.Hour, common.Config.OrderTime.Minutes, common.Config.OrderTime.Seconds, 0, now.Location())
	tests := []struct {
		name string
		want time.Time
	}{
		{
			name: "HappyCase",
			want: expectedTime,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetLunchTime(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetLunchTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConvertTimeStamp(t *testing.T) {
	timeNow := time.Unix(time.Now().Unix(), 0).Local().UTC()
	tz, _ := time.LoadLocation(TimeZone)
	s := timeNow.In(tz).Format("2006-01-02")

	type args struct {
		timestamp int64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "HappyCase",
			args: args{timestamp: time.Now().Unix()},
			want: s,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConvertTimeStamp(tt.args.timestamp); got != tt.want {
				t.Errorf("ConvertTimeStamp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConvertTimeStampTime(t *testing.T) {
	timeNow := time.Unix(time.Now().Unix(), 0).Local().UTC()
	tz, _ := time.LoadLocation(TimeZone)
	s := timeNow.In(tz).Format("3:04PM")

	type args struct {
		timestamp int64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "HappyCase",
			args: args{time.Now().Unix()},
			want: s,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConvertTimeStampTime(tt.args.timestamp); got != tt.want {
				t.Errorf("ConvertTimeStampTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWeekStartEndDate(t *testing.T) {
	type args struct {
		timestamp int64
	}
	tests := []struct {
		name  string
		args  args
		want  int64
		want1 int64
	}{
		{
			name:  "HappyCase",
			args:  args{1653262392},
			want:  1653235200,
			want1: 1653667199,
		},
		{
			name:  "EndsInNewMonth",
			args:  args{1653969600},
			want:  1653840000,
			want1: 1654271999,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := WeekStartEndDate(tt.args.timestamp)
			if got != tt.want {
				t.Errorf("WeekStartEndDate() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("WeekStartEndDate() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestConvertTimeStampMonthDay(t *testing.T) {
	timeNow := time.Unix(time.Now().Unix(), 0).Local().UTC()
	tz, _ := time.LoadLocation(TimeZone)
	s := timeNow.In(tz).Format("2/1")

	type args struct {
		timestamp int64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "HappyCase",
			args: args{time.Now().Unix()},
			want: s,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConvertTimeStampMonthDay(tt.args.timestamp); got != tt.want {
				t.Errorf("ConvertTimeStampMonthDay() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConvertTimeStampDayOfWeek(t *testing.T) {
	timeNow := time.Unix(time.Now().Unix(), 0).Local().UTC()
	tz, _ := time.LoadLocation(TimeZone)
	s := timeNow.In(tz).Format("Mon 02/01")

	type args struct {
		timestamp int64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "HappyCase",
			args: args{time.Now().Unix()},
			want: s,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConvertTimeStampDayOfWeek(tt.args.timestamp); got != tt.want {
				t.Errorf("ConvertTimeStampDayOfWeek() = %v, want %v", got, tt.want)
			}
		})
	}
}

func AdhocTestIsWeekDay(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{
			name: "HappyCase",
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsWeekDay(); got != tt.want {
				t.Errorf("IsWeekDay() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetPreviousDayLunchTime(t *testing.T) {
	tz, _ := time.LoadLocation(TimeZone)
	now := time.Now().In(tz)
	year, month, day := now.Date()
	expectedTime := time.Date(year, month, day-1, common.Config.OrderTime.Hour, common.Config.OrderTime.Minutes, common.Config.OrderTime.Seconds, 0, now.Location())
	tests := []struct {
		name string
		want time.Time
	}{
		{
			name: "HappyCase",
			want: expectedTime,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetPreviousDayLunchTime(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetPreviousDayLunchTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetOffWorkTime(t *testing.T) {
	tz, _ := time.LoadLocation(TimeZone)
	now := time.Now().In(tz)
	year, month, day := now.Date()
	expectedTime := time.Date(year, month, day, 19, 0, 0, 0, now.Location())
	tests := []struct {
		name string
		want time.Time
	}{
		{
			name: "HappyCase",
			want: expectedTime,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetOffWorkTime(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetOffWorkTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMonthStartEndDate(t *testing.T) {
	now := int64(1655658904) //20/6/2022
	expected := int64(1654012800)
	expected1 := int64(1656604799)
	type args struct {
		timestamp int64
	}
	tests := []struct {
		name  string
		args  args
		want  int64
		want1 int64
	}{
		{
			"HappyCase",
			args{now},
			expected,
			expected1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := MonthStartEndDate(tt.args.timestamp)
			if got != tt.want {
				t.Errorf("MonthStartEndDate() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("MonthStartEndDate() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestYearStartEndDate(t *testing.T) {
	now := int64(1655658904) //20/6/2022
	expected := int64(1640966400)
	expected1 := int64(1672502399)
	type args struct {
		timestamp int64
	}
	tests := []struct {
		name  string
		args  args
		want  int64
		want1 int64
	}{
		{
			"HappyCase",
			args{now},
			expected,
			expected1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := YearStartEndDate(tt.args.timestamp)
			if got != tt.want {
				t.Errorf("YearStartEndDate() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("YearStartEndDate() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestConvertTimeStampWeekOfYear(t *testing.T) {
	ts := int64(1656649800)  //1/7/2022
	ts1 := int64(1656563400) //30/6/2022
	type args struct {
		timestamp int64
	}
	tests := []struct {
		name  string
		args  args
		want  int64
		want1 int64
	}{
		{
			"HappyCase",
			args{timestamp: ts},
			2022,
			26,
		},
		{
			"HappyCase1",
			args{timestamp: ts1},
			2022,
			26,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := ConvertTimeStampWeekOfYear(tt.args.timestamp)
			if got != tt.want {
				t.Errorf("ConvertTimeStampWeekOfYear() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("ConvertTimeStampWeekOfYear() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
