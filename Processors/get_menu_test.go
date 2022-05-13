package Processors

import (
	"os"
	"reflect"
	"testing"

	"github.com/go-resty/resty/v2"
)

func TestGetMenu(t *testing.T) {
	r := Init()
	LoadEnv()
	exp := DinnerMenuArr{
		Status: "success",
	}

	type args struct {
		client resty.Client
		ID     int
		key    string
	}
	tests := []struct {
		name string
		args args
		want DinnerMenuArr
	}{
		{
			name: "HappyCase",
			args: args{client: r, ID: 1234, key: os.Getenv("Token")},
			want: exp,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetMenu(tt.args.client, tt.args.ID, tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetMenu() = %v, want %v", got, tt.want)
			}
		})
	}
}
