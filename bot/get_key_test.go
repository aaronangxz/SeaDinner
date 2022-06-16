package handlers

import (
	"context"
	"github.com/aaronangxz/SeaDinner/test_helper"
	"github.com/aaronangxz/SeaDinner/test_helper/user_key"
	"testing"
)

func TestGetKey(t *testing.T) {
	u := user_key.New().Build()
	randID := test_helper.RandomInt(9999999)
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
			args: args{id: u.GetUserId()},
			want: u.GetUserKey(),
		},
		{
			name: "HappyCaseCached",
			args: args{id: u.GetUserId()},
			want: u.GetUserKey(),
		},
		{
			name: "NotFound",
			args: args{id: randID},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetKey(context.TODO(), tt.args.id); got != tt.want {
				t.Errorf("GetKey() = %v, want %v", got, tt.want)
			}
		})
	}
	u.TearDown()
	user_key.DeleteUserKey(randID)
}
