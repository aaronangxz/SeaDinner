package handlers

import (
	"context"
	"fmt"
	"github.com/aaronangxz/SeaDinner/log"
	"github.com/aaronangxz/SeaDinner/processors"
	"github.com/aaronangxz/SeaDinner/sea_dinner.pb"
	"time"
)

func ListResultByUserID(ctx context.Context, id int64, timeRange int64) (string, bool) {
	var (
		start   int64
		end     int64
		res     []*sea_dinner.OrderRecord
		keyword string
	)
	txn := processors.App.StartTransaction("list_result_by_user_id")
	defer txn.End()

	if id <= 0 {
		log.Error(ctx, "Id must be > 1.")
		return "", false
	}

	switch timeRange {
	case int64(sea_dinner.ResultTimeRange_RESULT_TIME_RANGE_WEEK):
		start, end = processors.WeekStartEndDate(time.Now().Unix())
		keyword = "week"
	case int64(sea_dinner.ResultTimeRange_RESULT_TIME_RANGE_MONTH):
		start, end = processors.MonthStartEndDate(time.Now().Unix())
		keyword = "month"
	case int64(sea_dinner.ResultTimeRange_RESULT_TIME_RANGE_YEAR):
		start, end = processors.YearStartEndDate(time.Now().Unix())
		keyword = "year"
	}

	if err := processors.DbInstance().Raw("SELECT * FROM order_log_tab WHERE user_id = ? AND order_time BETWEEN ? AND ?", id, start, end).Scan(&res).Error; err != nil {
		log.Error(ctx, "id : %v | Failed to retrieve record.", id)
		return fmt.Sprintf("You have not ordered anything this %v. 😕", keyword), true
	}

	if res == nil {
		return fmt.Sprintf("You have not ordered anything this %v. 😕", keyword), true
	}
	log.Info(ctx, "ListResultByUserID | Success | user_id:%v", id)
	return GenerateResultTable(ctx, res, start, end), true
}
