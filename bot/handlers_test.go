package handlers

import (
	"context"
	"fmt"
	"github.com/aaronangxz/SeaDinner/common"
	"testing"
	"time"

	"github.com/aaronangxz/SeaDinner/log"
	"github.com/aaronangxz/SeaDinner/processors"
	"github.com/aaronangxz/SeaDinner/sea_dinner.pb"
	"github.com/aaronangxz/SeaDinner/test_helper"
	"github.com/aaronangxz/SeaDinner/test_helper/user_choice"
	"github.com/aaronangxz/SeaDinner/test_helper/user_key"
	"google.golang.org/protobuf/proto"
)

func TestMain(m *testing.M) {
	log.InitializeLogger()
	m.Run()
}
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

func TestCheckChope(t *testing.T) {
	m := test_helper.GetLiveMenuDetails()
	u := user_choice.New().SetUserChoice(fmt.Sprint(m[0].GetId())).Build()
	stopOrder := user_choice.New().SetUserChoice(fmt.Sprint(-1)).Build()
	notInMenu := user_choice.New().SetUserChoice(fmt.Sprint(999999)).Build()
	randOrder := user_choice.New().SetUserChoice("RAND").Build()
	tz, _ := time.LoadLocation(processors.TimeZone)
	var expected string = "Not placing dinner order for you today 🙅 Changed your mind? You can choose from /menu"
	if time.Now().In(tz).Unix() > processors.GetLunchTime().Unix() {
		if processors.IsNotEOW(time.Now().In(tz)) {
			expected = "Not placing dinner order for you tomorrow 🙅 Changed your mind? You can choose from /menu"
		} else {
			expected = "We are done for this week! You can tell me your order again next week 😀"
		}
	}

	defer func() {
		u.TearDown()
		stopOrder.TearDown()
		notInMenu.TearDown()
		randOrder.TearDown()
	}()

	type args struct {
		id int64
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 bool
	}{
		{
			name:  "HappyCase",
			args:  args{u.GetUserId()},
			want:  fmt.Sprintf("I'm tasked to snatch %v for you 😀 Changed your mind? You can choose from /menu", m[0].GetName()),
			want1: true,
		},
		{
			name:  "InvalidId",
			args:  args{id: -1},
			want:  "",
			want1: false,
		},
		{
			name:  "NoOrder",
			args:  args{id: 1},
			want:  "I have yet to receive your order 🥲 You can choose from /menu",
			want1: false,
		},
		{
			name:  "StopOrder",
			args:  args{id: stopOrder.GetUserId()},
			want:  expected,
			want1: false,
		},
		{
			name:  "OrderNotInMenu",
			args:  args{id: notInMenu.GetUserId()},
			want:  fmt.Sprintf("Your choice %v is not available this week, so I will not order anything 🥲 Choose a new dish from /menu", notInMenu.GetUserChoice()),
			want1: true,
		},
		{
			name:  "RandomOrder",
			args:  args{id: randOrder.GetUserId()},
			want:  "I'm tasked to snatch a random dish for you 😀 Changed your mind? You can choose from /menu",
			want1: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := CheckChope(context.TODO(), tt.args.id)
			if got != tt.want {
				t.Errorf("CheckChope() got1 = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("CheckChope() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestGetChope(t *testing.T) {
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
			want:  "Okay got it. I will order *NOTHING* for you and stop sending reminders in the morning for the rest of the week.😀",
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
			want:  "Okay got it. I will order *NOTHING* for you and stop sending reminders in the morning for the rest of the week.😀",
			want1: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := GetChope(context.TODO(), tt.args.id, tt.args.s)
			if got != tt.want {
				t.Errorf("GetChope() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetChope() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestCheckMute(t *testing.T) {
	u1 := user_key.New().Build()
	u2 := user_key.New().SetIsMute(int64(sea_dinner.MuteStatus_MUTE_STATUS_YES)).Build()
	defer func() {
		u1.TearDown()
		u2.TearDown()
	}()
	type args struct {
		id int64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "HappyCase-UnMute",
			args: args{u1.GetUserId()},
			want: "Daily reminder notifications are *ON*.\nDo you want to turn it OFF?",
		},
		{
			name: "HappyCase-Mute",
			args: args{u2.GetUserId()},
			want: "Daily reminder notifications are *OFF*.\nDo you want to turn it ON?",
		},
		{
			name: "NotExist",
			args: args{12345},
			want: "Record not found.",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := CheckMute(context.TODO(), tt.args.id)
			if got != tt.want {
				t.Errorf("CheckMute() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpdateMute(t *testing.T) {
	u1 := user_key.New().Build()
	u2 := user_key.New().Build()
	u3 := user_key.New().SetIsMute(int64(sea_dinner.MuteStatus_MUTE_STATUS_YES)).Build()
	u4 := user_key.New().SetIsMute(int64(sea_dinner.MuteStatus_MUTE_STATUS_NO)).Build()
	defer func() {
		u1.TearDown()
		u2.TearDown()
		u3.TearDown()
		u4.TearDown()
	}()
	type args struct {
		id       int64
		callback string
	}
	tests := []struct {
		name     string
		args     args
		want     string
		want1    bool
		expected int64
	}{
		{
			name:     "HappyCaseNoRecord-Mutes",
			args:     args{u1.GetUserId(), "MUTE"},
			want:     "Daily reminder notifications are *OFF*.\nDo you want to turn it ON?",
			want1:    true,
			expected: int64(sea_dinner.MuteStatus_MUTE_STATUS_YES),
		},
		{
			name:     "HappyCaseNoRecord-Unmutes",
			args:     args{u2.GetUserId(), "UNMUTE"},
			want:     "Daily reminder notifications are *ON*.\nDo you want to turn it OFF?",
			want1:    false,
			expected: int64(sea_dinner.MuteStatus_MUTE_STATUS_NO),
		},
		{
			name:     "HappyCaseMuted-Unmutes",
			args:     args{u3.GetUserId(), "UNMUTE"},
			want:     "Daily reminder notifications are *ON*.\nDo you want to turn it OFF?",
			want1:    false,
			expected: int64(sea_dinner.MuteStatus_MUTE_STATUS_NO),
		},
		{
			name:     "HappyCaseUnmuted-Mutes",
			args:     args{u4.GetUserId(), "MUTE"},
			want:     "Daily reminder notifications are *OFF*.\nDo you want to turn it ON?",
			want1:    true,
			expected: int64(sea_dinner.MuteStatus_MUTE_STATUS_YES),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := UpdateMute(context.TODO(), tt.args.id, tt.args.callback)
			if got != tt.want {
				t.Errorf("UpdateMute() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("UpdateMute() got1 = %v, want %v", got1, tt.want1)
			}

			c := user_key.CheckUserKey(tt.args.id)

			if c.IsMute == nil {
				t.Errorf("UpdateMute() record = %v, expected %v", nil, tt.expected)
			} else {
				if c.GetIsMute() != tt.expected {
					t.Errorf("UpdateMute() record = %v, expected %v", c.GetIsMute(), tt.expected)
				}
			}
		})
	}
}

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
