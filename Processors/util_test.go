package Processors

import (
	"testing"
)

func TestMakeToken(t *testing.T) {
	Config.Prefix.TokenPrefix = "Token "
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
			args: args{key: "ABCDEFG"},
			want: "Token ABCDEFG",
		}, {
			name: "EmptyString",
			args: args{key: ""},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MakeToken(tt.args.key); got != tt.want {
				t.Errorf("MakeToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMakeURL(t *testing.T) {
	Config.Prefix.UrlPrefix = "https://dinner.sea.com"
	type args struct {
		opt int
		id  *int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "DayId_happy_case",
			args: args{opt: URL_CURRENT, id: nil},
			want: "https://dinner.sea.com/api/current",
		}, {
			name: "Menu_happy_case",
			args: args{opt: URL_MENU, id: Int(1)},
			want: "https://dinner.sea.com/api/menu/1",
		}, {
			name: "Menu_no_id",
			args: args{opt: URL_MENU, id: nil},
			want: "",
		}, {
			name: "Order_happy_case",
			args: args{opt: URL_ORDER, id: Int(1)},
			want: "https://dinner.sea.com/api/order/1",
		}, {
			name: "Order_no_id",
			args: args{opt: URL_MENU, id: nil},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MakeURL(tt.args.opt, tt.args.id); got != tt.want {
				t.Errorf("MakeURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
