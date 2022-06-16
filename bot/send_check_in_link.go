package handlers

import (
	"context"
	"fmt"
	"github.com/aaronangxz/SeaDinner/common"
	"github.com/aaronangxz/SeaDinner/log"
	"github.com/aaronangxz/SeaDinner/processors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

//SendCheckInLink Verify if the user indeed has a valid order and sends the updated check-in link of the day
func SendCheckInLink(ctx context.Context) {
	var (
		txt        = "Check in now to collect your food!\nLink will expire at 8.30pm."
		buttonText = "Check in"
		out        []tgbotapi.InlineKeyboardMarkup
	)
	txn := processors.App.StartTransaction("send_check_in_link")
	defer txn.End()

	//Decode dynamic URL from static QR
	url, err := common.DecodeQR()
	if err != nil {
		log.Error(ctx, "SendCheckInLink | error:%v", err.Error())
		return
	}

	orders := BatchGetSuccessfulOrder(ctx)
	bot, err := tgbotapi.NewBotAPI(common.GetTGToken(ctx))
	if err != nil {
		log.Error(ctx, err.Error())
	}
	bot.Debug = true
	log.Info(ctx, "Authorized on account %s", bot.Self.UserName)

	for _, user := range orders {
		var buttons []tgbotapi.InlineKeyboardButton
		buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonURL(buttonText, url))
		out = append(out, tgbotapi.NewInlineKeyboardMarkup(buttons))

		msg := tgbotapi.NewMessage(user, "")
		msg.Text = txt
		msg.ReplyMarkup = out[0]

		if msgTrace, err := bot.Send(msg); err != nil {
			log.Error(ctx, err.Error())
		} else {
			//Save into set as <user_id>:<message_id>
			toWrite := fmt.Sprint(user, ":", msgTrace.MessageID)
			if err := processors.RedisClient.SAdd("checkin_link", toWrite).Err(); err != nil {
				log.Error(ctx, "SendCheckInLink | Error while writing to redis: %v", err.Error())
			} else {
				log.Info(ctx, "SendCheckInLink | Successful | Written %v to checkin_link set", toWrite)
			}
		}
	}
}
