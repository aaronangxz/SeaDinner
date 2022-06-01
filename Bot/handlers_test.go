package Bot

import (
	"fmt"
	"testing"
	"time"

	"github.com/aaronangxz/SeaDinner/Bot/TestHelper"
	"github.com/aaronangxz/SeaDinner/Bot/TestHelper/user_choice"
	"github.com/aaronangxz/SeaDinner/Bot/TestHelper/user_key"
	"github.com/aaronangxz/SeaDinner/Processors"
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
	tz, _ := time.LoadLocation(Processors.TimeZone)
	var dayText string = "today"
	if time.Now().In(tz).Unix() > Processors.GetLunchTime().Unix() {
		if Processors.IsNotEOW(time.Now().In(tz)) {
			dayText = "tomorrow"
		}
	}
	defer func() {
		u.TearDown()
		stopOrder.TearDown()
		notInMenu.TearDown()
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
			want:  fmt.Sprintf("Not placing dinner order for you %v 🙅 Changed your mind? You can choose from /menu", dayText),
			want1: false,
		},
		{
			name:  "OrderNotInMenu",
			args:  args{id: notInMenu.GetUserId()},
			want:  fmt.Sprintf("Your choice %v is not available this week, so I will not order anything 🥲 Choose a new dish from /menu", notInMenu.GetUserChoice()),
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
	m := TestHelper.GetLiveMenuDetails()
	u := user_choice.New().Build()
	u1 := user_choice.New().SetUserChoice(fmt.Sprint(m[0].GetId())).Build()
	expected := "Okay got it. I will order %v for you today 😙"
	if time.Now().Unix() > Processors.GetLunchTime().Unix() {
		expected = "Okay got it. I will order %v for you tomorrow 😙"
	}

	defer func() {
		u.TearDown()
		u1.TearDown()
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
			want:  "Okay got it. I will order *NOTHING* for you and stop sending reminders in the morning.😀",
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
