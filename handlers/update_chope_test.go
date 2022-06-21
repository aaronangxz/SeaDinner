package handlers

import (
	"context"
	"fmt"
	"github.com/aaronangxz/SeaDinner/common"
	"github.com/aaronangxz/SeaDinner/log"
	"github.com/aaronangxz/SeaDinner/processors"
	"github.com/aaronangxz/SeaDinner/test_helper"
	"github.com/aaronangxz/SeaDinner/test_helper/user_choice"
	"testing"
	"time"
)

func TestUpdateChope(t *testing.T) {
	tz, _ := time.LoadLocation(processors.TimeZone)
	if !processors.IsNotEOW(time.Now().In(tz)) {
		log.Info(context.TODO(), "TestGetChope | Skipping")
		return
	}

	expiry := 60 * time.Second
	m := test_helper.GetLiveMenuDetails()
	u := user_choice.New().Build()
	u1 := user_choice.New().SetUserChoice(fmt.Sprint(m[0].GetId())).Build()
	u5 := user_choice.New().SetUserChoice(fmt.Sprint(m[0].GetId())).Build()

	var us = []*user_choice.UserChoice{
		user_choice.New().SetUserChoice(fmt.Sprint(m[1].GetId())).Build(),
		user_choice.New().SetUserChoice("RAND").Build(),
		user_choice.New().SetUserChoice("-1").Build()}

	for _, uu := range us {
		key := fmt.Sprint(common.USER_CHOICE_PREFIX, uu.GetUserId())
		if err := processors.RedisClient.Set(key, uu.GetUserChoice(), expiry).Err(); err != nil {
			log.Error(context.TODO(), "TestGetChope | Error while writing to redis: %v", err.Error())
		} else {
			log.Info(context.TODO(), "TestGetChope | Successful | Written %v to redis", key)
		}
	}

	expected := "Okay got it. I will order %v for you today 😙"
	if time.Now().Unix() > processors.GetLunchTime().Unix() {
		expected = "Okay got it. I will order %v for you tomorrow 😙"
	}

	defer func() {
		u.TearDown()
		u1.TearDown()
		for _, uu := range us {
			key := fmt.Sprint(common.USER_CHOICE_PREFIX, uu.GetUserId())
			if _, err := processors.RedisClient.Del(key).Result(); err != nil {
				log.Error(context.TODO(), "TestGetChope | Failed to invalidate cache: %v. %v", key, err)
			}
			log.Info(context.TODO(), "TestGetChope | Successfully invalidated cache: %v", key)
			uu.TearDown()
		}
		u5.TearDown()
	}()

	type args struct {
		id int64
		s  string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 bool
	}{
		{
			name:  "HappyCase",
			args:  args{id: u.GetUserId(), s: fmt.Sprint(m[0].GetId())},
			want:  fmt.Sprintf(expected, m[0].GetName()),
			want1: true,
		},
		{
			name:  "InvalidId",
			args:  args{id: -1},
			want:  "",
			want1: false,
		},
		{
			name:  "Alphabets",
			args:  args{id: u.GetUserId(), s: "ABCDEF"},
			want:  "Are you sure that is a valid FoodID? Tell me another one. 😟",
			want1: false,
		},
		{
			name:  "SpecialChar",
			args:  args{id: u.GetUserId(), s: "!@#$%^"},
			want:  "Are you sure that is a valid FoodID? Tell me another one. 😟",
			want1: false,
		},
		{
			name:  "NotInMenu",
			args:  args{id: u.GetUserId(), s: fmt.Sprint(6969)},
			want:  "This dish is not available today. Tell me another one.😟",
			want1: false,
		},
		{
			name:  "UpdateEntry",
			args:  args{id: u1.GetUserId(), s: fmt.Sprint(m[1].GetId())},
			want:  fmt.Sprintf(expected, m[1].GetName()),
			want1: true,
		},
		{
			name:  "StopOrder",
			args:  args{id: u.GetUserId(), s: fmt.Sprint(-1)},
			want:  "Okay got it. I will order <b>NOTHING</b> for you and stop sending reminders in the morning for the rest of the week.😀",
			want1: true,
		},
		{
			name:  "RandomOrder",
			args:  args{id: u.GetUserId(), s: "RAND"},
			want:  "Okay got it. I will give you a surprise 😙",
			want1: true,
		},
		{
			name:  "SameOrderWithFoodId",
			args:  args{id: us[0].GetUserId(), s: "SAME"},
			want:  fmt.Sprintf("Okay got it! I will order %v 😙", m[1].GetName()),
			want1: true,
		},
		{
			name:  "SameOrderWithRAND",
			args:  args{id: us[1].GetUserId(), s: "SAME"},
			want:  "Okay got it. I will give you a surprise 😙",
			want1: true,
		},
		{
			name:  "SameOrderWith-1",
			args:  args{id: us[2].GetUserId(), s: "SAME"},
			want:  "Okay got it. I will not order anything for you instead.😀",
			want1: true,
		},
		{
			name:  "OrderWithFoodId",
			args:  args{id: u5.GetUserId(), s: fmt.Sprint(m[0].GetId())},
			want:  fmt.Sprintf(expected, m[0].GetName()),
			want1: true,
		},
		{
			name:  "OrderWithRAND",
			args:  args{id: u5.GetUserId(), s: "RAND"},
			want:  "Okay got it. I will give you a surprise 😙",
			want1: true,
		},
		{
			name:  "OrderWith-1",
			args:  args{id: u5.GetUserId(), s: "-1"},
			want:  "Okay got it. I will order <b>NOTHING</b> for you and stop sending reminders in the morning for the rest of the week.😀",
			want1: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := UpdateChope(context.TODO(), tt.args.id, tt.args.s)
			if got != tt.want {
				t.Errorf("UpdateChope() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("UpdateChope() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
