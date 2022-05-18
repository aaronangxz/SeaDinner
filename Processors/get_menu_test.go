package Processors

import (
	"os"
	"reflect"
	"testing"

	"github.com/go-resty/resty/v2"
)

func TestGetMenu(t *testing.T) {
	LoadEnv()
	LoadConfig()
	Config.Adhoc = true
	r := Init()
	key := os.Getenv("TOKEN")
	type args struct {
		client resty.Client
		ID     int
		key    string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "HappyCase",
			args: args{client: r, ID: 3521, key: key},
			want: 8,
		},
		{
			name: "InvalidID",
			args: args{client: r, ID: 0, key: key},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetMenu(tt.args.client, tt.args.ID, tt.args.key); !reflect.DeepEqual(len(got.DinnerArr), tt.want) {
				t.Errorf("GetMenu() = %v, want %v", len(got.DinnerArr), tt.want)
			}
		})
	}
}

func AdhocTestOutputMenu(t *testing.T) {
	LoadEnv()
	Init()
	key := os.Getenv("TOKEN")
	type args struct {
		key string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "HappyCase",
			args: args{key},
			want: "There is no dinner order today! 😕",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := OutputMenu(tt.args.key); got != tt.want {
				t.Errorf("OutputMenu() = %v, want %v", got, tt.want)
			}
		})
	}
}
