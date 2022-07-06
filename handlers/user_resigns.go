package handlers

import (
	"context"
	"github.com/aaronangxz/SeaDinner/log"
	"github.com/aaronangxz/SeaDinner/sea_dinner.pb"
)

func UserResigns(ctx context.Context, id int64) (string, bool) {
	if err := UpdateUserStatus(ctx, id, int64(sea_dinner.UserStatus_USER_STATUS_RESIGNED)); err != nil {
		log.Error(ctx, "UserResigns | Error | %v", err.Error())
	}
	return "Okay, it was great serving you. Goodbye 🥺", true
}
