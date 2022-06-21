package processors

import (
	"context"
	"fmt"
	"github.com/aaronangxz/SeaDinner/common"
	"github.com/aaronangxz/SeaDinner/log"
	"github.com/aaronangxz/SeaDinner/sea_dinner.pb"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

//SendInstantNotification Spawns a one-time telegram handlers instance and send notification to user
func SendInstantNotification(ctx context.Context, u *sea_dinner.UserChoiceWithKey, took int64) {
	var (
		mk   tgbotapi.InlineKeyboardMarkup
		out  [][]tgbotapi.InlineKeyboardButton
		rows []tgbotapi.InlineKeyboardButton
	)
	txn := App.StartTransaction("send_instant_notifications")
	defer txn.End()

	bot, err := tgbotapi.NewBotAPI(common.GetTGToken(ctx))
	if err != nil {
		log.Error(ctx, err.Error())
	}
	bot.Debug = true
	log.Info(ctx, "Authorized on account %s", bot.Self.UserName)

	menu := MakeMenuMap(ctx)
	msg := tgbotapi.NewMessage(u.GetUserId(), "")
	msg.Text = fmt.Sprintf("Successfully ordered %v in %vms! 🥳", menu[u.GetUserChoice()], took)

	skipBotton := tgbotapi.NewInlineKeyboardButtonData("I DON'T NEED IT 🙅", "ATTEMPTCANCEL")
	rows = append(rows, skipBotton)
	out = append(out, rows)
	mk.InlineKeyboard = out
	msg.ReplyMarkup = mk
	if _, err := bot.Send(msg); err != nil {
		log.Error(ctx, err.Error())
	}
	log.Info(ctx, "SendInstantNotification | user_id:%v | msg: %v", u.GetUserId(), msg)
}
