package handlers

import (
	"context"
	"github.com/aaronangxz/SeaDinner/log"
	"github.com/aaronangxz/SeaDinner/processors"
	"github.com/aaronangxz/SeaDinner/sea_dinner.pb"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

//CheckMute Checks the user's current status of mute state
func CheckMute(ctx context.Context, id int64) (string, []tgbotapi.InlineKeyboardMarkup) {
	var (
		res *sea_dinner.UserKey
		out []tgbotapi.InlineKeyboardMarkup
	)
	txn := processors.App.StartTransaction("check_mute")
	defer txn.End()

	if err := processors.DbInstance().Raw("SELECT * FROM user_key_tab WHERE user_id = ?", id).Scan(&res).Error; err != nil {
		log.Error(ctx, "CheckMute | Failed to retrieve record: %v", err.Error())
		return "", nil
	}

	if res == nil {
		log.Error(ctx, "CheckMute | Record not found | user_id:%v", id)
		return "Record not found.", nil
	}

	if res.GetIsMute() == int64(sea_dinner.MuteStatus_MUTE_STATUS_NO) {
		var rows []tgbotapi.InlineKeyboardButton
		muteBotton := tgbotapi.NewInlineKeyboardButtonData("Turn OFF 🔕", "MUTE")
		rows = append(rows, muteBotton)
		out = append(out, tgbotapi.NewInlineKeyboardMarkup(rows))
		return "Daily reminder notifications are *ON*.\nDo you want to turn it OFF?", out
	}
	var rows []tgbotapi.InlineKeyboardButton
	unmuteBotton := tgbotapi.NewInlineKeyboardButtonData("Turn ON 🔔", "UNMUTE")
	rows = append(rows, unmuteBotton)
	out = append(out, tgbotapi.NewInlineKeyboardMarkup(rows))
	return "Daily reminder notifications are *OFF*.\nDo you want to turn it ON?", out
}
