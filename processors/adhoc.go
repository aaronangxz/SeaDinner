package processors

import (
	"context"
	"github.com/aaronangxz/SeaDinner/common"

	"github.com/aaronangxz/SeaDinner/log"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

//SendAdHocNotification Sends a one-off message to a specific user
func SendAdHocNotification(ctx context.Context, id int64, text string) {

	bot, err := tgbotapi.NewBotAPI(common.GetTGToken(ctx))
	if err != nil {
		log.Error(ctx, err.Error())
	}
	bot.Debug = true
	log.Info(ctx, "Authorized on account %s", bot.Self.UserName)

	msg := tgbotapi.NewMessage(id, "")
	msg.Text = text

	if _, err := bot.Send(msg); err != nil {
		log.Error(ctx, err.Error())
	}
}
