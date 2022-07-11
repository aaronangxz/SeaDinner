package processors

import (
	"context"
	"github.com/aaronangxz/SeaDinner/log"
)

//UpdateOrderLog Update a single record in order_log_tab
func UpdateOrderLog(ctx context.Context, userId int64, status int64) error {
	txn := App.StartTransaction("update_order_log")
	defer txn.End()

	order, getOrderLogErr := GetOrderLog(ctx, userId)
	if getOrderLogErr != nil {
		log.Error(ctx, "UpdateOrderLog | Failed to retrieve record | %v", getOrderLogErr.Error())
		return getOrderLogErr
	}

	if err := DbInstance().Exec("UPDATE order_log_tab SET status = ? WHERE id = ? AND user_id = ?", status, order.GetId(), userId).Error; err != nil {
		log.Error(ctx, "UpdateOrderLog | Failed to update record | %v", err.Error())
		return err
	}
	log.Info(ctx, "UpdateOrderLog | Successfully updated record | user_id: %v", userId)
	return nil
}
