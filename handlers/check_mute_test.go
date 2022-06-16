package handlers

import (
	"context"
	"github.com/aaronangxz/SeaDinner/sea_dinner.pb"
	"github.com/aaronangxz/SeaDinner/test_helper/user_key"
	"testing"
)

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
