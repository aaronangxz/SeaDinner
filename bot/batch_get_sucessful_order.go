package handlers

import (
	"context"
	"github.com/aaronangxz/SeaDinner/log"
	"github.com/aaronangxz/SeaDinner/processors"
)

//BatchGetSuccessfulOrder Calls Sea API to verify the user's current order
func BatchGetSuccessfulOrder(ctx context.Context) []int64 {
	var (
		success []int64
	)
	txn := processors.App.StartTransaction("batch_get_successful_order")
	defer txn.End()

	records, err := BatchGetUsersChoiceWithKey(ctx)
	if err != nil {
		log.Error(ctx, "BatchGetSuccessfulOrder | Failed to fetch user_records: %v", err.Error())
		return nil
	}

	for _, r := range records {
		ok := processors.GetSuccessfulOrder(ctx, r.GetUserKey())
		if ok {
			success = append(success, r.GetUserId())
		} else {
			log.Error(ctx, "BatchGetSuccessfulOrder | Failed | user_id: %v", r.GetUserId())
		}
	}
	log.Info(ctx, "BatchGetSuccessfulOrder | Done | size: %v", len(success))
	return success
}
