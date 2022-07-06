package handlers

import (
	"context"
	"github.com/aaronangxz/SeaDinner/log"
	"github.com/aaronangxz/SeaDinner/sea_dinner.pb"
)

func UserOptsOut(ctx context.Context, id int64) (string, bool) {
	if err := UpdateUserStatus(ctx, id, int64(sea_dinner.UserStatus_USER_STATUS_INACTIVE)); err != nil {
		log.Error(ctx, "UserOptsOut | Error | %v", err.Error())
	}
	return "Okay, hope to see you soon again. Goodbye 🥺", true
}
