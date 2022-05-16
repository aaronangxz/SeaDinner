package Processors

// import (
// 	"os"
// 	"testing"
// )

// func TestGetDayId(t *testing.T) {
// 	LoadEnv()
// 	Init()
// 	Config.Prefix.TokenPrefix = "Token "
// 	Config.Prefix.UrlPrefix = "https://dinner.sea.com"
// 	type args struct {
// 		key string
// 	}
// 	tests := []struct {
// 		name   string
// 		args   args
// 		wantID int
// 	}{
// 		{
// 			name:   "HappyCase",
// 			args:   args{key: os.Getenv("TOKEN")},
// 			wantID: 3521,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if gotID := GetDayId(tt.args.key); gotID != tt.wantID {
// 				t.Errorf("GetDayId() = %v, want %v", gotID, tt.wantID)
// 			}
// 		})
// 	}
// }
