package handlers

import (
	"context"
	"github.com/aaronangxz/SeaDinner/sea_dinner.pb"
	"github.com/aaronangxz/SeaDinner/test_helper/user_key"
	"testing"
)

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
