package handlers

import (
	"context"
	"github.com/aaronangxz/SeaDinner/log"
	"github.com/aaronangxz/SeaDinner/processors"
	"github.com/aaronangxz/SeaDinner/sea_dinner.pb"
)

//BatchGetLatestResult Retrieves the most recent failed orders
func BatchGetLatestResult(ctx context.Context) []*sea_dinner.OrderRecord {
	var (
		res []*sea_dinner.OrderRecord
	)
	txn := processors.App.StartTransaction("batch_get_latest_result")
	defer txn.End()

	if err := processors.DB.Raw("SELECT ol.* FROM order_log_tab ol INNER JOIN "+
		"(SELECT MAX(order_time) AS max_order_time FROM order_log_tab WHERE status <> ? AND order_time BETWEEN ? AND ? GROUP BY user_id) nestedQ "+
		"ON ol.order_time = nestedQ.max_order_time GROUP BY user_id",
		sea_dinner.OrderStatus_ORDER_STATUS_OK, processors.GetLunchTime().Unix()-300, processors.GetLunchTime().Unix()+300).
		Scan(&res).Error; err != nil {
		log.Error(ctx, "Failed to retrieve record.")
		return nil
	}
	log.Info(ctx, "BatchGetLatestResult: %v", len(res))
	return res
}
