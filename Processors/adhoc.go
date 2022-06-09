package Processors

import (
	"log"

	"github.com/aaronangxz/SeaDinner/Common"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

//SendAdHocNotification Sends a one-off message to a specific user
func SendAdHocNotification(id int64, text string) {

	bot, err := tgbotapi.NewBotAPI(Common.GetTGToken())
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	msg := tgbotapi.NewMessage(id, "")
	msg.Text = text

	if _, err := bot.Send(msg); err != nil {
		log.Println(err)
	}
}
