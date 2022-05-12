package Processors

import (
	"reflect"
	"testing"

	"github.com/go-resty/resty/v2"
)

func TestGetMenu(t *testing.T) {
	r := Init()

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
			args: args{client: r, ID: 1234, key: "8f983bf2f8dfb706713896c8aa9174646e3e37c2"},
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
