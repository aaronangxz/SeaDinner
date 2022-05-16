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
