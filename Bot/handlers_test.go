package Bot

import (
	"os"
	"testing"

	"github.com/aaronangxz/SeaDinner/Bot/TestHelper"
	"github.com/aaronangxz/SeaDinner/Bot/TestHelper/user_key"
	"github.com/aaronangxz/SeaDinner/Processors"
)

func TestGetKey(t *testing.T) {
	u := user_key.New().Build()
	type args struct {
		id int64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "HappyCase",
			args: args{id: u.GetUserID()},
			want: Processors.DecryptKey(u.GetUserKey(), os.Getenv("AES_KEY")),
		},
		{
			name: "NotFound",
			args: args{id: TestHelper.RandomInt(999)},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetKey(tt.args.id); got != tt.want {
				t.Errorf("GetKey() = %v, want %v", got, tt.want)
			}
		})
	}
	u.TearDown()
}

func TestCheckKey(t *testing.T) {
	u := user_key.New().Build()
	type args struct {
		id int64
	}
	tests := []struct {
		name  string
		args  args
		want1 bool
	}{
		{
			name:  "HappyCase",
			args:  args{u.GetUserID()},
			want1: true,
		},
		{
			name:  "NotFound",
			args:  args{id: TestHelper.RandomInt(999)},
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, got1 := CheckKey(tt.args.id)
			if got1 != tt.want1 {
				t.Errorf("CheckKey() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	u.TearDown()
}

func TestUpdateKey(t *testing.T) {
	u := user_key.New().Build()
	type args struct {
		id int64
		s  string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "KeyInvalidLen",
			args: args{u.GetUserID(), TestHelper.RandomString(39)},
			want: false,
		},
		{
			name: "KeyEmpty",
			args: args{u.GetUserID(), ""},
			want: false,
		},
		{
			name: "UserKeyNotExist",
			args: args{TestHelper.RandomInt(999), TestHelper.RandomString(40)},
			want: true,
		},
		{
			name: "UserExistsButKeyNotExist",
			args: args{u.GetUserID(), TestHelper.RandomString(40)},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, got := UpdateKey(tt.args.id, tt.args.s)
			if got != tt.want {
				t.Errorf("UpdateKey() got = %v, want %v", got, tt.want)
			}
		})
	}
	u.TearDown()
}
