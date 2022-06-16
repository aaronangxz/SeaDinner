package handlers

import (
	"context"
	"fmt"
	"github.com/aaronangxz/SeaDinner/common"
	"github.com/aaronangxz/SeaDinner/log"
	"github.com/aaronangxz/SeaDinner/processors"
	"github.com/aaronangxz/SeaDinner/sea_dinner.pb"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

//SendNotifications Sends out notifications based on order status from BatchGetLatestResult
//Used to send failed orders only
func SendNotifications(ctx context.Context) {
	var (
		msg string
	)
	txn := processors.App.StartTransaction("send_notifications")
	defer txn.End()

	bot, err := tgbotapi.NewBotAPI(common.GetTGToken(ctx))
	if err != nil {
		log.Error(ctx, err.Error())
	}
	bot.Debug = true
	log.Info(ctx, "Authorized on account %s", bot.Self.UserName)

	res := BatchGetLatestResult(processors.Ctx)
	menu := MakeMenuNameMap(processors.Ctx)
	log.Info(ctx, "SendNotifications | size: %v", len(res))

	for _, r := range res {
		if r.GetStatus() == int64(sea_dinner.OrderStatus_ORDER_STATUS_OK) {
			msg = fmt.Sprintf("Successfully ordered %v! 🥳", menu[r.GetFoodId()])
		} else {
			msg = fmt.Sprintf("Failed to order %v today. %v 😔", menu[r.GetFoodId()], r.GetErrorMsg())
		}

		if _, err := bot.Send(tgbotapi.NewMessage(r.GetUserId(), msg)); err != nil {
			log.Error(ctx, err.Error())
		}
	}
}
