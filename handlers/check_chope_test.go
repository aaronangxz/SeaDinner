package handlers

import (
	"context"
	"fmt"
	"github.com/aaronangxz/SeaDinner/processors"
	"github.com/aaronangxz/SeaDinner/test_helper"
	"github.com/aaronangxz/SeaDinner/test_helper/user_choice"
	"testing"
	"time"
)

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
