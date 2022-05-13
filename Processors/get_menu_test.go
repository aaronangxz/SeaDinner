package Processors

import (
	"os"
	"reflect"
	"testing"

	"github.com/go-resty/resty/v2"
)

func TestGetMenu(t *testing.T) {
	LoadEnv()
	r := Init()

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
			args: args{client: r, ID: 3521, key: os.Getenv("TOKEN")},
			want: 8,
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
