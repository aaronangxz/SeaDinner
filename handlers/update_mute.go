package handlers

import (
	"context"
	"github.com/aaronangxz/SeaDinner/log"
	"github.com/aaronangxz/SeaDinner/processors"
	"github.com/aaronangxz/SeaDinner/sea_dinner.pb"
)

//UpdateMute Updates the user's current status of mute state
func UpdateMute(ctx context.Context, id int64, callback string) (string, bool) {
	var (
		toUdate    = int64(sea_dinner.MuteStatus_MUTE_STATUS_YES)
		returnMsg  = "Daily reminder notifications are *OFF*.\nDo you want to turn it ON?"
		returnBool = true
	)
	txn := processors.App.StartTransaction("update_mute")
	defer txn.End()

	if callback == "UNMUTE" {
		toUdate = int64(sea_dinner.MuteStatus_MUTE_STATUS_NO)
		returnMsg = "Daily reminder notifications are *ON*.\nDo you want to turn it OFF?"
		returnBool = false
	}

	if err := processors.DB.Exec("UPDATE user_key_tab SET is_mute = ? WHERE user_id = ?", toUdate, id).Error; err != nil {
		log.Error(ctx, "Failed to update DB")
		return err.Error(), false
	}
	log.Info(ctx, "UpdateMute | Success")
	return returnMsg, returnBool
}
