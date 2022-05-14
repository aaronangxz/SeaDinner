package Bot

import (
	"testing"

	"github.com/aaronangxz/SeaDinner/Bot/TestHelper"
	"github.com/aaronangxz/SeaDinner/Bot/TestHelper/user_key"
)

func TestGetKey(t *testing.T) {
	u := user_key.New().Build()
	err := "I can't do this in a group chat! PM me instead 😉"
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
			want: u.GetKey(),
		},
		{
			name: "GroupChat",
			args: args{id: -1},
			want: err,
		},
		{
			name: "NotFound",
			args: args{id: 1},
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
			name:  "GroupChat",
			args:  args{id: -1},
			want1: false,
		},
		{
			name:  "NotFound",
			args:  args{id: 1},
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
			name: "UserNotExistKeyNotExist",
			args: args{u.GetUserID(), TestHelper.RandomString(40)},
			want: true,
		},
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
			name: "GroupChat",
			args: args{-1, TestHelper.RandomString(40)},
			want: false,
		},
		{
			name: "UserKeyNotExist",
			args: args{1, TestHelper.RandomString(40)},
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
