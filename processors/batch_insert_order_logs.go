package processors

import (
	"context"
	"fmt"
	"github.com/aaronangxz/SeaDinner/common"
	"github.com/aaronangxz/SeaDinner/log"
	"github.com/aaronangxz/SeaDinner/sea_dinner.pb"
)

//BatchInsertOrderLogs Batch insert new order records into order_log_tab
func BatchInsertOrderLogs(ctx context.Context, records []*sea_dinner.OrderRecord) {
	txn := App.StartTransaction("batch_insert_order_logs")
	defer txn.End()

	if records == nil {
		log.Warn(ctx, "BatchInsertOrderLogs | No record to update.")
		return
	}
	if err := DB.Table(common.DB_ORDER_LOG_TAB).Create(&records).Error; err != nil {
		log.Error(ctx, fmt.Sprintf("BatchInsertOrderLogs | Failed to update records | %v", err.Error()))
		return
	}
	log.Info(ctx, fmt.Sprintf("BatchInsertOrderLogs | Successfully updated records | size: %v", len(records)))
}
