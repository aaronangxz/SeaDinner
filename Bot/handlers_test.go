package Bot

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/aaronangxz/SeaDinner/Common"
	"github.com/aaronangxz/SeaDinner/Processors"
	"github.com/aaronangxz/SeaDinner/TestHelper"
	"github.com/aaronangxz/SeaDinner/TestHelper/user_choice"
	"github.com/aaronangxz/SeaDinner/TestHelper/user_key"
	"github.com/aaronangxz/SeaDinner/sea_dinner.pb"
)

func TestGetKey(t *testing.T) {
	u := user_key.New().Build()
	randId := TestHelper.RandomInt(9999999)
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
			args: args{id: randId},
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
	user_key.DeleteUserKey(randId)
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
	randU := TestHelper.RandomInt(999)
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
			args: args{u.GetUserId(), TestHelper.RandomString(39)},
			want: false,
		},
		{
			name: "KeyEmpty",
			args: args{u.GetUserId(), ""},
			want: false,
		},
		{
			name: "UserKeyNotExist",
			args: args{randU, TestHelper.RandomString(40)},
			want: true,
		},
		{
			name: "UserExistsButKeyNotExist",
			args: args{u.GetUserId(), TestHelper.RandomString(40)},
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
	user_key.DeleteUserKey(randU)
}

func TestCheckChope(t *testing.T) {
	m := TestHelper.GetLiveMenuDetails()
	u := user_choice.New().SetUserChoice(fmt.Sprint(m[0].GetId())).Build()
	stopOrder := user_choice.New().SetUserChoice(fmt.Sprint(-1)).Build()
	notInMenu := user_choice.New().SetUserChoice(fmt.Sprint(999999)).Build()
	randOrder := user_choice.New().SetUserChoice("RAND").Build()
	tz, _ := time.LoadLocation(Processors.TimeZone)
	var expected string = "Not placing dinner order for you today 🙅 Changed your mind? You can choose from /menu"
	if time.Now().In(tz).Unix() > Processors.GetLunchTime().Unix() {
		if Processors.IsNotEOW(time.Now().In(tz)) {
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
			got, got1 := CheckChope(tt.args.id)
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
	tz, _ := time.LoadLocation(Processors.TimeZone)
	if !Processors.IsNotEOW(time.Now().In(tz)) {
		log.Printf("TestGetChope | Skipping")
		return
	}

	expiry := 60 * time.Second
	m := TestHelper.GetLiveMenuDetails()
	u := user_choice.New().Build()
	u1 := user_choice.New().SetUserChoice(fmt.Sprint(m[0].GetId())).Build()
	u5 := user_choice.New().SetUserChoice(fmt.Sprint(m[0].GetId())).Build()

	var us = []*user_choice.UserChoice{
		user_choice.New().SetUserChoice(fmt.Sprint(m[1].GetId())).Build(),
		user_choice.New().SetUserChoice("RAND").Build(),
		user_choice.New().SetUserChoice("-1").Build()}

	for _, uu := range us {
		key := fmt.Sprint(Common.USER_CHOICE_PREFIX, uu.GetUserId())
		if err := Processors.RedisClient.Set(key, uu.GetUserChoice(), expiry).Err(); err != nil {
			log.Printf("TestGetChope | Error while writing to redis: %v", err.Error())
		} else {
			log.Printf("TestGetChope | Successful | Written %v to redis", key)
		}
	}

	expected := "Okay got it. I will order %v for you today 😙"
	if time.Now().Unix() > Processors.GetLunchTime().Unix() {
		expected = "Okay got it. I will order %v for you tomorrow 😙"
	}

	defer func() {
		u.TearDown()
		u1.TearDown()
		for _, uu := range us {
			key := fmt.Sprint(Common.USER_CHOICE_PREFIX, uu.GetUserId())
			if _, err := Processors.RedisClient.Del(key).Result(); err != nil {
				log.Printf("TestGetChope | Failed to invalidate cache: %v. %v", key, err)
			}
			log.Printf("TestGetChope | Successfully invalidated cache: %v", key)
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
			got, got1 := GetChope(tt.args.id, tt.args.s)
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
			got, _ := CheckMute(tt.args.id)
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
			got, got1 := UpdateMute(tt.args.id, tt.args.callback)
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
