package handlers

import (
	"context"
	"github.com/aaronangxz/SeaDinner/log"
	"github.com/aaronangxz/SeaDinner/processors"
)

//UpdateUserStatus Update the status record in user_key_tab
func UpdateUserStatus(ctx context.Context, id int64, status int64) error {
	txn := processors.App.StartTransaction("update_user_status")
	defer txn.End()

	if err := processors.DbInstance().Exec("UPDATE user_key_tab SET status = ? WHERE user_id = ?", status, id).Error; err != nil {
		log.Error(ctx, "UpdateUserStatus | Failed to update record | %v", err.Error())
		return err
	}
	log.Info(ctx, "UpdateUserStatus | Successfully updated status | user_id: %v", id)
	return nil
}
