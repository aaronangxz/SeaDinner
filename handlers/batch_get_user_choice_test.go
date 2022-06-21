package handlers

import (
	"context"
	"github.com/aaronangxz/SeaDinner/sea_dinner.pb"
	"github.com/aaronangxz/SeaDinner/test_helper"
	"github.com/aaronangxz/SeaDinner/test_helper/user_choice"
	"github.com/aaronangxz/SeaDinner/test_helper/user_key"
	"google.golang.org/protobuf/proto"
	"testing"
)

func TestBatchGetUsersChoice(t *testing.T) {
	ctx := context.TODO()
	uk := user_key.New().Build()
	uc := user_choice.New().SetUserId(uk.GetUserId()).Build()
	defer func() {
		uc.TearDown()
		uk.TearDown()
	}()
	expected := []*sea_dinner.UserChoice{
		{
			UserId:     proto.Int64(uc.GetUserId()),
			UserChoice: proto.String(uc.GetUserChoice()),
			Ctime:      proto.Int64(uc.GetCtime()),
			Mtime:      proto.Int64(uc.GetMtime()),
		},
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want []*sea_dinner.UserChoice
	}{
		{
			"HappyCase",
			args{ctx: ctx},
			expected,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BatchGetUsersChoice(tt.args.ctx)
			if !test_helper.IsInSlice(tt.want, got) {
				t.Errorf("BatchGetUsersChoice() = %v, want %v", got, tt.want)
			}
		})
	}
}
