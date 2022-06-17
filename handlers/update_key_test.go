package handlers

import (
	"context"
	"github.com/aaronangxz/SeaDinner/test_helper"
	"github.com/aaronangxz/SeaDinner/test_helper/user_key"
	"testing"
)

func TestUpdateKey(t *testing.T) {
	u := user_key.New().Build()
	randU := test_helper.RandomInt(999)
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
			args: args{u.GetUserId(), test_helper.RandomString(39)},
			want: false,
		},
		{
			name: "KeyEmpty",
			args: args{u.GetUserId(), ""},
			want: false,
		},
		{
			name: "UserKeyNotExist",
			args: args{randU, test_helper.RandomString(40)},
			want: true,
		},
		{
			name: "UserExistsButKeyNotExist",
			args: args{u.GetUserId(), test_helper.RandomString(40)},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, got := UpdateKey(context.TODO(), tt.args.id, tt.args.s)
			if got != tt.want {
				t.Errorf("UpdateKey() got = %v, want %v", got, tt.want)
			}
		})
	}
	u.TearDown()
	user_key.DeleteUserKey(randU)
}
