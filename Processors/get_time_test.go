package Processors

import (
	"reflect"
	"testing"
	"time"
)

func TestGetLunchTime(t *testing.T) {
	tz, _ := time.LoadLocation(TimeZone)
	now := time.Now().In(tz)
	year, month, day := now.Date()
	expectedTime := time.Date(year, month, day, Config.OrderTime.Hour, Config.OrderTime.Minutes, Config.OrderTime.Seconds, 0, now.Location())
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
	s := (timeNow.In(tz).Format("2006-01-02"))

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
	s := (timeNow.In(tz).Format("3:04PM"))

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

func TestIsWeekDay(t *testing.T) {
	type args struct {
		t time.Time
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "HappyCase_weekday",
			args: args{time.Now()},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsWeekDay(tt.args.t); got != tt.want {
				t.Errorf("IsWeekDay() = %v, want %v", got, tt.want)
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
	s := (timeNow.In(tz).Format("2/1"))

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
	s := (timeNow.In(tz).Format("Mon 2/1"))

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
