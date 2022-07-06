package handlers

import (
	"context"
	"github.com/aaronangxz/SeaDinner/log"
	"github.com/aaronangxz/SeaDinner/processors"
	"github.com/aaronangxz/SeaDinner/sea_dinner.pb"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

//CallbackQueryHandler Handles the call back result of menu buttons
func CallbackQueryHandler(ctx context.Context, id int64, callBack *tgbotapi.CallbackQuery) (string, bool) {
	txn := processors.App.StartTransaction("call_back_query_handler")
	defer txn.End()

	log.Info(ctx, "id: %v | CallbackQueryHandler | callback: %v", id, callBack.Data)

	switch callBack.Data {
	case "MUTE":
		fallthrough
	case "UNMUTE":
		return UpdateMute(ctx, id, callBack.Data)
	case "ATTEMPTCANCEL":
		return "", true
	case "ATTEMPTOPTOUT":
		return "", true
	case "ATTEMPTRESIGN":
		return "", true
	case "CANCEL":
		return CancelOrder(ctx, id)
	case "OPTOUT":
		return UserOptsOut(ctx, id)
	case "RESIGN":
		return UserResigns(ctx, id)
	case "SKIP":
		return "Alright 😗", true
	case "WEEKRESULT":
		return ListResultByUserID(ctx, id, int64(sea_dinner.ResultTimeRange_RESULT_TIME_RANGE_WEEK))
	case "MONTHRESULT":
		return ListResultByUserID(ctx, id, int64(sea_dinner.ResultTimeRange_RESULT_TIME_RANGE_MONTH))
		//Not active for now
		//case "YEARRESULT":
		//	return ListResultByUserID(ctx, id, int64(sea_dinner.ResultTimeRange_RESULT_TIME_RANGE_YEAR))
	}
	return UpdateChope(ctx, id, callBack.Data)
}
