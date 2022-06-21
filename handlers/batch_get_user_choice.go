package handlers

import (
	"context"
	"fmt"
	"github.com/aaronangxz/SeaDinner/common"
	"github.com/aaronangxz/SeaDinner/log"
	"github.com/aaronangxz/SeaDinner/processors"
	"github.com/aaronangxz/SeaDinner/sea_dinner.pb"
	"time"
)

//BatchGetUsersChoice Retrieves order_choice of all users
func BatchGetUsersChoice(ctx context.Context) []*sea_dinner.UserChoice {
	var (
		res    []*sea_dinner.UserChoice
		expiry = 7200 * time.Second
	)
	txn := processors.App.StartTransaction("batch_get_user_choice")
	defer txn.End()

	if err := processors.DbInstance().Raw("SELECT uc.* FROM user_choice_tab uc, user_key_tab uk WHERE uc.user_id = uk.user_id AND uk.is_mute <> ?", sea_dinner.MuteStatus_MUTE_STATUS_YES).Scan(&res).Error; err != nil {
		log.Error(ctx, "BatchGetUsersChoice | Failed to retrieve record: %v", err.Error())
		return nil
	}

	//Save into cache
	//For Morning Reminder callback
	for _, r := range res {
		//Not necessary to cache -1 orders because we never send reminder for those
		if r.GetUserChoice() != "-1" {
			key := fmt.Sprint(common.USER_CHOICE_PREFIX, r.GetUserId())
			if err := processors.CacheInstance().Set(key, r.GetUserChoice(), expiry).Err(); err != nil {
				log.Error(ctx, "BatchGetUsersChoice | Error while writing to redis: %v", err.Error())
			} else {
				log.Info(ctx, "BatchGetUsersChoice | Successful | Written %v to redis", key)
			}
		}
	}
	log.Info(ctx, "BatchGetUsersChoice | size: %v", len(res))
	return res
}
