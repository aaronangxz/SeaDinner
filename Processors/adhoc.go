package Processors

import (
	"context"
	"log"

	"github.com/aaronangxz/SeaDinner/Common"
	"github.com/aaronangxz/SeaDinner/Log"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

//SendAdHocNotification Sends a one-off message to a specific user
func SendAdHocNotification(ctx context.Context, id int64, text string) {

	bot, err := tgbotapi.NewBotAPI(Common.GetTGToken(ctx))
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	Log.Info(ctx, "Authorized on account %s", bot.Self.UserName)
	// log.Printf("Authorized on account %s", bot.Self.UserName)

	msg := tgbotapi.NewMessage(id, "")
	msg.Text = text

	if _, err := bot.Send(msg); err != nil {
		Log.Error(ctx, err.Error())
		// log.Println(err)
	}
}
