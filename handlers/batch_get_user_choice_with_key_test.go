package handlers

import (
	"context"
	"fmt"
	"github.com/aaronangxz/SeaDinner/sea_dinner.pb"
	"github.com/aaronangxz/SeaDinner/test_helper"
	"github.com/aaronangxz/SeaDinner/test_helper/user_choice"
	"github.com/aaronangxz/SeaDinner/test_helper/user_key"
	"google.golang.org/protobuf/proto"
	"testing"
)

func TestBatchGetUsersChoiceWithKey(t *testing.T) {
	m := test_helper.GetLiveMenuDetails()
	u := user_choice.New().SetUserChoice(fmt.Sprint(m[0].GetId())).Build()
	uk := user_key.New().SetUserId(u.GetUserId()).Build()

	expected1 := []*sea_dinner.UserChoiceWithKey{
		{
			UserId:     proto.Int64(u.GetUserId()),
			UserKey:    proto.String(uk.GetUserKey()),
			UserChoice: proto.String(u.GetUserChoice()),
			Ctime:      proto.Int64(u.GetCtime()),
			Mtime:      proto.Int64(u.GetMtime()),
		},
	}

	u2 := user_choice.New().SetUserChoice(fmt.Sprint(m[0].GetId())).Build()
	uk2 := user_key.New().SetUserId(u2.GetUserId()).Build()

	expected2 := []*sea_dinner.UserChoiceWithKey{
		{
			UserId:     proto.Int64(u2.GetUserId()),
			UserKey:    proto.String(uk2.GetUserKey()),
			UserChoice: proto.String(u2.GetUserChoice()),
			Ctime:      proto.Int64(u2.GetCtime()),
			Mtime:      proto.Int64(u2.GetMtime()),
		},
	}

	defer func() {
		u.TearDown()
		uk.TearDown()
		u2.TearDown()
		uk2.TearDown()
	}()

	tests := []struct {
		name    string
		want    []*sea_dinner.UserChoiceWithKey
		wantErr bool
	}{
		{
			name:    "HappyCase",
			want:    expected1,
			wantErr: false,
		},
		{
			name:    "HasRAND",
			want:    expected2,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BatchGetUsersChoiceWithKey(context.TODO())
			if (err != nil) != tt.wantErr {
				t.Errorf("BatchGetUsersChoiceWithKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !test_helper.IsInSlice(tt.want, got) {
				t.Errorf("BatchGetUsersChoiceWithKey() = %v, want %v", got, tt.want)
			}
		})
	}
}
