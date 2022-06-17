package handlers

import (
	"context"
	"github.com/aaronangxz/SeaDinner/test_helper"
	"github.com/aaronangxz/SeaDinner/test_helper/user_key"
	"testing"
)

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
			args:  args{u.GetUserId()},
			want1: true,
		},
		{
			name:  "HappyCaseCached",
			args:  args{u.GetUserId()},
			want1: true,
		},
		{
			name:  "NotFound",
			args:  args{id: test_helper.RandomInt(999)},
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, got1 := CheckKey(context.TODO(), tt.args.id)
			if got1 != tt.want1 {
				t.Errorf("CheckKey() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	u.TearDown()
}
