package processors

import (
	"context"
	"github.com/aaronangxz/SeaDinner/common"
	"github.com/aaronangxz/SeaDinner/log"
	"github.com/aaronangxz/SeaDinner/sea_dinner.pb"
	"sync"
)

//BatchOrderDinnerMultiThreadedWithWait Spawns individual go routines before lunchtime
func BatchOrderDinnerMultiThreadedWithWait(ctx context.Context, userQueue []*sea_dinner.UserChoiceWithKey) {
	var (
		wg      sync.WaitGroup
		records []*sea_dinner.OrderRecord
	)
	txn := App.StartTransaction("batch_order_dinner_multi_threaded_with_wait")
	defer txn.End()

	m := make(map[int64]int64)

	for _, user := range userQueue {
		if !common.IsInGrayScale(user.GetUserId()) {
			log.Info(ctx, "BatchOrderDinnerMultiThreadedWithWait | Not in grayscale, skipping | user_id:%v", user.GetUserId())
			continue
		}
		//Increment group
		wg.Add(1)
		go func(u *sea_dinner.UserChoiceWithKey) {
			//Release group
			defer wg.Done()
			var record *sea_dinner.OrderRecord
			for {
				if IsOrderTime() && IsPollStart() {
					log.Info(ctx, "BatchOrderDinnerMultiThreadedWithWait | Begin | user_id: %v", u.GetUserId())
					m[u.GetUserId()], record = OrderDinnerWithUpdate(ctx, u)
					records = append(records, record)
					break
				}
			}
		}(user)
	}

	//Wait for all groups to release
	wg.Wait()
	log.Info(ctx, "BatchOrderDinnerMultiThreadedWithWait | Done")
	BatchInsertOrderLogs(ctx, records)
	OutputResults(ctx, m, "BatchOrderDinnerMultiThreadedWithWait")
}
