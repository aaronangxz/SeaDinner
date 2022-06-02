package Processors

import (
	"os"
	"testing"

	"github.com/aaronangxz/SeaDinner/Common"
	"github.com/aaronangxz/SeaDinner/sea_dinner.pb"
	"google.golang.org/protobuf/proto"
)

func TestMakeToken(t *testing.T) {
	LoadEnv()
	key := "ogNiXZrVyXZglYPZHmhoF7J9JvQzxaIINBRgntOA"
	shortKey := "abcde"
	eKey := EncryptKey(key, os.Getenv("AES_KEY"))
	shortEKey := EncryptKey(shortKey, os.Getenv("AES_KEY"))
	type args struct {
		key string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "HappyCase",
			args: args{key: eKey},
			want: false,
		},
		{
			name: "EmptyString",
			args: args{key: ""},
			want: true,
		},
		{
			name: "LessThanMinLength",
			args: args{key: shortEKey},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MakeToken(tt.args.key); (got == tt.args.key) != tt.want {
				t.Errorf("MakeToken() = %v, want %v", (got == tt.args.key), tt.want)
			}
		})
	}
}

func TestMakeURL(t *testing.T) {
	Common.Config.Prefix.UrlPrefix = "https://dinner.sea.com"
	type args struct {
		opt int
		id  *int64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "DayId_happy_case",
			args: args{opt: int(sea_dinner.URLType_URL_CURRENT), id: nil},
			want: "https://dinner.sea.com/api/current",
		}, {
			name: "Menu_happy_case",
			args: args{opt: int(sea_dinner.URLType_URL_MENU), id: proto.Int64(1)},
			want: "https://dinner.sea.com/api/menu/1",
		}, {
			name: "Menu_no_id",
			args: args{opt: int(sea_dinner.URLType_URL_MENU), id: nil},
			want: "",
		}, {
			name: "Order_happy_case",
			args: args{opt: int(sea_dinner.URLType_URL_ORDER), id: proto.Int64(1)},
			want: "https://dinner.sea.com/api/order/1",
		}, {
			name: "Order_no_id",
			args: args{opt: int(sea_dinner.URLType_URL_MENU), id: nil},
			want: "",
		},
		{
			name: "InvalidOpt",
			args: args{opt: 3, id: nil},
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

func TestOutputResults(t *testing.T) {
	m := make(map[int64]int64)
	m[1] = int64(sea_dinner.OrderStatus_ORDER_STATUS_OK)
	m[2] = int64(sea_dinner.OrderStatus_ORDER_STATUS_FAIL)
	type args struct {
		resultMap map[int64]int64
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Happy Case",
			args: args{m},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			OutputResults(tt.args.resultMap)
		})
	}
}

func TestIsNotNumber(t *testing.T) {
	type args struct {
		a string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "HappyCase",
			args: args{"12345"},
			want: false,
		},
		{
			name: "SpecialChar",
			args: args{"!@#$%"},
			want: true,
		},
		{
			name: "ChineseChar",
			args: args{"哈哈哈哈"},
			want: true,
		},
		{
			name: "Emojis",
			args: args{"😍😍😍😍"},
			want: true,
		},
		{
			name: "Alphabets",
			args: args{"ABCDE"},
			want: true,
		},
		{
			name: "Alphanumeric",
			args: args{"ABC123"},
			want: true,
		},
		{
			name: "BeginWithNumber",
			args: args{"123ABC"},
			want: true,
		},
		{
			name: "EmptyString",
			args: args{""},
			want: true,
		},
		{
			name: "Minus One (-1)",
			args: args{"-1"},
			want: false,
		},
		{
			name: "Minus 1.0 (-1.0)",
			args: args{"-1.0"},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsNotNumber(tt.args.a); got != tt.want {
				t.Errorf("IsNotNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEncryptKey(t *testing.T) {
	type args struct {
		stringToEncrypt string
		keyString       string
	}
	tests := []struct {
		name         string
		args         args
		isSameString bool
	}{
		{
			name:         "HappyCase",
			args:         args{stringToEncrypt: "SOMESTRING", keyString: MakeKey()},
			isSameString: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotEncryptedString := EncryptKey(tt.args.stringToEncrypt, tt.args.keyString); (gotEncryptedString == tt.args.stringToEncrypt) != tt.isSameString {
				t.Errorf("EncryptKey() = %v, want %v", gotEncryptedString, tt.isSameString)
			}
		})
	}
}

func TestDecryptKey(t *testing.T) {
	originalString := "SomeString"
	key := MakeKey()
	enc := EncryptKey(originalString, key)
	type args struct {
		encryptedString string
		keyString       string
	}
	tests := []struct {
		name                string
		args                args
		wantDecryptedString string
	}{
		{
			name:                "HappyCase",
			args:                args{encryptedString: enc, keyString: key},
			wantDecryptedString: originalString,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotDecryptedString := DecryptKey(tt.args.encryptedString, tt.args.keyString); gotDecryptedString != tt.wantDecryptedString {
				t.Errorf("DecryptKey() = %v, want %v", gotDecryptedString, tt.wantDecryptedString)
			}
		})
	}
}

func TestRandomFood(t *testing.T) {
	m := make(map[string]string)
	m["A"] = "1"
	m["B"] = "2"
	m["C"] = "3"
	m["D"] = "4"
	m["E"] = "5"
	m["F"] = "6"
	m["RAND"] = "RAND"
	m["-1"] = "-1"

	type args struct {
		m map[string]string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "HappyCase",
			args: args{m},
			want: true,
		},
		{
			name: "HappyCaseAgain",
			args: args{m},
			want: true,
		},
		{
			name: "HappyCaseAndAgain",
			args: args{m},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RandomFood(tt.args.m); got != "" != tt.want {
				t.Errorf("RandomFood() = %v, want %v", got, tt.want)
			}
		})
	}
}
