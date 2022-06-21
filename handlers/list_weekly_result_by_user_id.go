package handlers

import (
	"context"
	"github.com/aaronangxz/SeaDinner/log"
	"github.com/aaronangxz/SeaDinner/processors"
	"github.com/aaronangxz/SeaDinner/sea_dinner.pb"
	"time"
)

//ListWeeklyResultByUserID Returns the order records of a user in the current week
func ListWeeklyResultByUserID(ctx context.Context, id int64) string {
	var (
		res []*sea_dinner.OrderRecord
	)
	txn := processors.App.StartTransaction("list_weekly_result_by_user_id")
	defer txn.End()

	if !processors.IsWeekDay() {
		log.Warn(ctx, "ListWeeklyResultByUserId | Not a weekday.")
		return "We are done for this week! Check again next week 😀"
	}

	start, end := processors.WeekStartEndDate(time.Now().Unix())

	if id <= 0 {
		log.Error(ctx, "Id must be > 1.")
		return ""
	}

	if err := processors.DB.Raw("SELECT * FROM order_log_tab WHERE user_id = ? AND order_time BETWEEN ? AND ?", id, start, end).Scan(&res).Error; err != nil {
		log.Error(ctx, "id : %v | Failed to retrieve record.", id)
		return "You have not ordered anything this week. 😕"
	}

	if res == nil {
		return "You have not ordered anything this week. 😕"
	}
	log.Info(ctx, "ListWeeklyResultByUserId | Success | user_id:%v", id)
	return GenerateWeeklyResultTableWithFoodMapping(ctx, res)
}
