package processors

import (
	"context"
	"github.com/aaronangxz/SeaDinner/common"
	"github.com/aaronangxz/SeaDinner/log"
	"github.com/aaronangxz/SeaDinner/sea_dinner.pb"
	"sync"
)

//BatchOrderDinnerMultiThreaded Spawns multiple Order goroutines, and update order_log_tab with the respective results.
//Guaranteed to execute goroutines for all users in the queue.
func BatchOrderDinnerMultiThreaded(ctx context.Context, userQueue []*sea_dinner.UserChoiceWithKey) {
	var (
		wg      sync.WaitGroup
		records []*sea_dinner.OrderRecord
	)
	txn := App.StartTransaction("batch_order_dinner_multi_threaded")
	defer txn.End()

	m := make(map[int64]int64)
	log.Info(ctx, "BatchOrderDinnerMultiThreaded | Begin | size: %v", len(userQueue))

	for _, user := range userQueue {
		if common.IsInGrayScale(user.GetUserId()) {
			log.Info(ctx, "BatchOrderDinnerMultiThreaded | In grayscale, skipping | user_id:%v", user.GetUserId())
			continue
		}
		//Increment group
		wg.Add(1)
		go func(u *sea_dinner.UserChoiceWithKey) {
			//Release group
			defer wg.Done()
			var record *sea_dinner.OrderRecord
			m[u.GetUserId()], record = OrderDinnerWithUpdate(ctx, u)
			records = append(records, record)
		}(user)
	}

	//Wait for all groups to release
	wg.Wait()

	log.Info(ctx, "BatchOrderDinnerMultiThreaded | Done")
	BatchInsertOrderLogs(ctx, records)
	OutputResults(ctx, m, "BatchOrderDinnerMultiThreaded")
}
