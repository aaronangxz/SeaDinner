package Processors

import (
	"os"
	"reflect"
	"testing"

	"github.com/aaronangxz/SeaDinner/Common"
	"github.com/go-resty/resty/v2"
)

func TestGetMenu(t *testing.T) {
	LoadEnv()
	Common.LoadConfig()
	ConnectTestRedis()
	ConnectTestMySQL()
	r := InitClient()
	key := os.Getenv("TOKEN")
	type args struct {
		client resty.Client
		key    string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "HappyCase",
			args: args{client: r, key: key},
			want: 12,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetMenu(tt.args.client, tt.args.key); !reflect.DeepEqual(len(got.GetFood()), tt.want) {
				t.Errorf("GetMenu() = %v, want %v", len(got.GetFood()), tt.want)
			}
		})
	}
}
