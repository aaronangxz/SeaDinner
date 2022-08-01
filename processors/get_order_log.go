package processors

import (
	"context"
	"github.com/aaronangxz/SeaDinner/log"
	"github.com/aaronangxz/SeaDinner/sea_dinner.pb"
)

//GetOrderLog Retrieves a single record in order_log_tab
func GetOrderLog(ctx context.Context, id int64) (*sea_dinner.OrderRecord, error) {
	var (
		order *sea_dinner.OrderRecord
	)
	txn := App.StartTransaction("update_order_log")
	defer txn.End()

	if err := DbInstance().Raw("SELECT * FROM order_log_tab WHERE user_id = ? ORDER BY id DESC LIMIT 1", id).Scan(&order).Error; err != nil {
		log.Error(ctx, "GetOrderLog | Failed to retrieve record | %v", err.Error())
		return nil, err
	}
	log.Info(ctx, "GetOrderLog | Successfully retrieved record | user_id: %v", id)
	return order, nil
}
