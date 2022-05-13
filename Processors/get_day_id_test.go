package Processors

import (
	"testing"
)

func TestGetDayId(t *testing.T) {
	Init()
	Config.Prefix.TokenPrefix = "Token "
	Config.Prefix.UrlPrefix = "https://dinner.sea.com"
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		args   args
		wantID int
	}{
		{
			name:   "HappyCase",
			args:   args{key: "8f983bf2f8dfb706713896c8aa9174646e3e37c2"},
			wantID: 3521,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotID := GetDayId(tt.args.key); gotID != tt.wantID {
				t.Errorf("GetDayId() = %v, want %v", gotID, tt.wantID)
			}
		})
	}
}
